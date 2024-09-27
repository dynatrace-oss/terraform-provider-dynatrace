package iam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

type IAMError struct {
	BError  bool   `json:"error"`
	Payload string `json:"payload"`
	Message string `json:"message"`
}

func (me IAMError) Error() string {
	return me.Message
}

type IAMClient interface {
	POST(url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	PUT(url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	GET(url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	DELETE(url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	DELETE_MULTI_RESPONSE(url string, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error)
}

type iamClient struct {
	auth Authenticator
}

func NewIAMClient(auth Authenticator) IAMClient {
	return &iamClient{auth}
}

func (me *iamClient) authenticate(httpRequest *http.Request, forceNew bool) error {
	var bearerToken string
	var err error

	if bearerToken, err = getBearer(me.auth, forceNew); err != nil {
		return err
	}

	if os.Getenv("DT_DEBUG_IAM_BEARER") == "true" {
		log.Println("--- BEARER TOKEN ---")
		log.Println(bearerToken)
		log.Println("--------------------")
	}
	if len(strings.TrimSpace(bearerToken)) == 0 {
		return errors.New(msgInvalidOAuthCredentials)
	}

	httpRequest.Header.Set("Authorization", "Bearer "+bearerToken)

	return nil
}

func (me *iamClient) POST(url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodPost, []int{expectedResponseCode}, forceNewBearer, 0, payload, map[string]string{"Content-Type": "application/json"})
}

func (me *iamClient) PUT(url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodPut, []int{expectedResponseCode}, forceNewBearer, 0, payload, map[string]string{"Content-Type": "application/json"})
}

func (me *iamClient) GET(url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodGet, []int{expectedResponseCode}, forceNewBearer, 0, nil, nil)
}

func (me *iamClient) DELETE(url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodDelete, []int{expectedResponseCode}, forceNewBearer, 0, nil, nil)
}

func (me *iamClient) DELETE_MULTI_RESPONSE(url string, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodDelete, expectedResponseCodes, forceNewBearer, 0, nil, nil)
}

type RateLimiter struct {
	lastCall time.Time
	mutex    sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

var DISABLE_RATE_LIMITER = (os.Getenv("DYNATRACE_DISABLE_IAM_RATE_LIMITER") == "true")

const MAX_RATE_LIMITER_RATE = int64(5000)
const DEFAULT_RATE_LIMITER_RATE = int64(1000)

var rateLimiterRate = evalRateLimiterRate()

func evalRateLimiterRate() int64 {
	sRateLimiterRate := strings.TrimSpace(os.Getenv("DYNATRACE_IAM_RATE_LIMITER_RATE"))
	if len(sRateLimiterRate) == 0 {
		return DEFAULT_RATE_LIMITER_RATE
	}
	var err error
	iRateLimiterRate := DEFAULT_RATE_LIMITER_RATE

	if iRateLimiterRate, err = strconv.ParseInt(sRateLimiterRate, 10, 0); err != nil {
		return DEFAULT_RATE_LIMITER_RATE
	}
	if iRateLimiterRate > MAX_RATE_LIMITER_RATE {
		return MAX_RATE_LIMITER_RATE
	}
	return iRateLimiterRate
}

func (rl *RateLimiter) CanCall() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	if rateLimiterRate <= 0 || DISABLE_RATE_LIMITER {
		return true
	}
	now := time.Now()
	if now.Sub(rl.lastCall) >= time.Duration(rateLimiterRate)*time.Millisecond {
		rl.lastCall = now
		return true
	}
	return false
}

var limiter = NewRateLimiter()

func httplog(v ...any) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	tt := fmt.Sprintf("[HTTP] [%s]", formattedTime)
	rest.Logger.Println(append([]any{tt}, v...)...)
}

func (me *iamClient) request(url string, method string, expectedResponseCodes []int, forceNewBearer bool, forceNewBearerRetryCount int, payload any, headers map[string]string) ([]byte, error) {
	for {
		if limiter.CanCall() {
			return me._request(url, method, expectedResponseCodes, forceNewBearer, forceNewBearerRetryCount, payload, headers)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (me *iamClient) _request(url string, method string, expectedResponseCodes []int, forceNewBearer bool, forceNewBearerRetryCount int, payload any, headers map[string]string) ([]byte, error) {
	// httplog(fmt.Sprintf("[%s] %s", method, url))

	num504Retries := 0
	sleepTime429 := int64(500)

	for {
		var err error
		var httpRequest *http.Request
		var httpResponse *http.Response
		var responseBytes []byte
		var requestBody []byte

		rest.Logger.Println(method, url)

		if requestBody, err = json.Marshal(payload); err != nil {
			return nil, err
		}
		if payload != nil {
			rest.Logger.Println("  ", string(requestBody))
		}

		var body io.Reader

		if payload != nil {
			body = bytes.NewReader(requestBody)
		}

		if httpRequest, err = http.NewRequest(method, url, body); err != nil {
			return nil, err
		}

		if err = me.authenticate(httpRequest, forceNewBearer); err != nil {
			return nil, err
		}

		for k, v := range headers {
			httpRequest.Header.Add(k, v)
		}

		if httpResponse, err = http.DefaultClient.Do(httpRequest); err != nil {
			return nil, err
		}

		if responseBytes, err = io.ReadAll(httpResponse.Body); err != nil {
			return nil, err
		}
		rest.Logger.Println("  ", httpResponse.StatusCode, string(responseBytes))

		if httpResponse.StatusCode == 504 {
			// httplog("-------------------- FIVE-O-FOUR --------------------")
			num504Retries++
			if num504Retries > 5 {
				return nil, fmt.Errorf("response code %d (expected: %d)", 504, expectedResponseCodes)
			}
			time.Sleep(time.Second)
			continue
		}

		isNotExpectedResponseCode := true
		for _, erc := range expectedResponseCodes {
			if httpResponse.StatusCode == erc {
				isNotExpectedResponseCode = false
				break
			}
		}

		if isNotExpectedResponseCode {
			if httpResponse.StatusCode == 429 {
				time.Sleep(time.Duration(sleepTime429) * time.Millisecond)
				// logging.File.Println(".... 429 ... waiting for another", sleepTime429, "milliseconds")
				sleepTime429 = int64(math.Round(float64(sleepTime429) * float64(1.6)))
				continue
			}
			var iamErr IAMError
			if err = json.Unmarshal(responseBytes, &iamErr); err == nil {
				if (forceNewBearerRetryCount < 20) && iamErr.Error() == "Failed to validate access token." {
					// httplog("-------------------- TOKEN-SWITCH --------------------")
					return me.request(url, method, expectedResponseCodes, true, forceNewBearerRetryCount+1, payload, headers)
				}
				return nil, iamErr
			} else {
				errEnv := struct {
					Error *rest.Error `json:"error"`
				}{}
				if len(responseBytes) > 0 {
					if err = json.Unmarshal(responseBytes, &errEnv); err == nil && errEnv.Error != nil {
						// {"error":{"code":400,"message":"Policy c1f6bc44-c66d-4964-938c-ab00823a5e22 can't be bound to levelType: environment, levelId: siz65484","errorsMap":null}}
						return nil, *errEnv.Error
					}
				}
				return nil, rest.Error{Code: httpResponse.StatusCode, Message: fmt.Sprintf("response code %d (expected: %d)", httpResponse.StatusCode, expectedResponseCodes)}
			}
		}

		return responseBytes, nil
	}
}
