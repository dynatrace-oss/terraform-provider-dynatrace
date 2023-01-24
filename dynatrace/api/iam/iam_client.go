package iam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	httpRequest.Header.Set("Authorization", "Bearer "+bearerToken)

	return nil
}

func (me *iamClient) POST(url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodPost, expectedResponseCode, forceNewBearer, payload, map[string]string{"Content-Type": "application/json"})
}

func (me *iamClient) PUT(url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodPut, expectedResponseCode, forceNewBearer, payload, map[string]string{"Content-Type": "application/json"})
}

func (me *iamClient) GET(url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodGet, expectedResponseCode, forceNewBearer, nil, nil)
}

func (me *iamClient) DELETE(url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.request(url, http.MethodDelete, expectedResponseCode, forceNewBearer, nil, nil)
}

func (me *iamClient) request(url string, method string, expectedResponseCode int, forceNewBearer bool, payload any, headers map[string]string) ([]byte, error) {
	var err error
	var httpRequest *http.Request
	var httpResponse *http.Response
	var responseBytes []byte
	var requestBody []byte

	if requestBody, err = json.Marshal(payload); err != nil {
		return nil, err
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
	if httpResponse.StatusCode != expectedResponseCode {
		var iamErr IAMError
		if err = json.Unmarshal(responseBytes, &iamErr); err == nil {
			if !forceNewBearer && iamErr.Error() == "Failed to validate access token." {
				return me.request(url, method, expectedResponseCode, true, payload, headers)
			}
			return nil, iamErr
		} else {
			return nil, fmt.Errorf("response code %d (expected: %d)", httpResponse.StatusCode, expectedResponseCode)
		}
	}

	return responseBytes, nil
}
