/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package rest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"golang.org/x/sync/semaphore"
)

const MinWaitTime = 5 * time.Second
const MaxWaitTime = 1 * time.Minute

var jar = createJar()

func createJar() *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	return jar
}

type statuscodes []int

func (me statuscodes) contains(code int) bool {
	for _, c := range me {
		if c == code {
			return true
		}
	}
	return false
}

type legacy_request request

func (me *legacy_request) authenticate(req *http.Request) {
	req.Header.Add("Authorization", "Api-Token "+me.client.Credentials().Token)
	req.Header.Set("User-Agent", "Dynatrace Terraform Provider")
}

func (me *legacy_request) Finish(vs ...any) error {
	if err := me.finish(vs...); err != nil {
		// If the Settings 2.0 API returned with an error
		// that contains constraint violations complaining
		// about unknown properties we may be able to
		// retry with payload that doesn't contain these
		// properties
		// if CorrectPayload(err, me) {
		// 	return me.finishRetry(1, vs...)
		// }
		return err
	}
	return nil
}

func (me *legacy_request) finish(vs ...any) error {
	if shutdown.System.Stopped() {
		return nil
	}
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}
	var err error
	var data []byte
	if data, err = me.Raw(); err != nil {
		return err
	}
	if shutdown.System.Stopped() {
		return nil
	}
	if v != nil {
		if err = json.Unmarshal(data, &v); err != nil {
			return fmt.Errorf("%s %s: unmarshal error: %s\n%s", me.method, me.url, err.Error(), string(data))
		}
	}
	return nil
}

func (me *legacy_request) Raw() ([]byte, error) {
	url := me.client.Credentials().URL + me.url
	var err error
	var body io.Reader
	var data []byte
	if me.payload != nil {
		if data, err = json.Marshal(me.payload); err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}

	if len(data) > 0 {
		Logger.Printf(me.ctx, "[legacy  ] [%s] [REQUEST ] %s %s", me.id, me.method, url)
		Logger.Printf(me.ctx, "           [%s] [PAYLOAD ] %s", me.id, string(data))
	} else {
		Logger.Printf(me.ctx, "[legacy  ] [%s] [REQUEST ] %s %s", me.id, me.method, url)
	}

	var req *http.Request
	if req, err = http.NewRequest(me.method, url, body); err != nil {
		return nil, err
	}
	me.authenticate(req)
	if len(me.headers) > 0 {
		for headerName, headerValue := range me.headers {
			req.Header.Add(headerName, headerValue)
		}
	}
	var res *http.Response

	httpClient := &http.Client{
		Jar:       jar,
		Transport: http.DefaultTransport,
	}
	if strings.TrimSpace(os.Getenv("DYNATRACE_HTTP_INSECURE")) == "true" {
		httpClient.Transport = &http.Transport{
			ForceAttemptHTTP2:     http.DefaultTransport.(*http.Transport).ForceAttemptHTTP2,
			Proxy:                 http.DefaultTransport.(*http.Transport).Proxy,
			DialContext:           http.DefaultTransport.(*http.Transport).DialContext,
			MaxIdleConns:          http.DefaultTransport.(*http.Transport).MaxIdleConns,
			IdleConnTimeout:       http.DefaultTransport.(*http.Transport).IdleConnTimeout,
			TLSHandshakeTimeout:   http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout,
			ExpectContinueTimeout: http.DefaultTransport.(*http.Transport).ExpectContinueTimeout,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		httpClient.Transport = http.DefaultTransport
	}
	response, err := me.execute(me.ctx, func() (*http.Response, error) {
		if res, err = httpClient.Do(req); err != nil {
			return nil, err
		}
		return res, nil
	})
	if shutdown.System.Stopped() {
		return nil, nil
	}
	if me.onResponse != nil {
		me.onResponse(response)
	}
	if err != nil {
		return nil, err
	}
	if data, err = io.ReadAll(res.Body); err != nil {
		return nil, err
	}
	if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
		if data != nil {
			Logger.Printf(me.ctx, "           [%s] [RESPONSE] %s %s", me.id, res.Status, string(data))
		} else {
			Logger.Printf(me.ctx, "           [%s] [RESPONSE] %s", me.id, res.Status)
		}
	}
	if len(me.expect) > 0 && !me.expect.contains(res.StatusCode) {
		var env errorEnvelope
		if err = json.Unmarshal(data, &env); err == nil && env.Error != nil {
			return nil, Error{Code: env.Error.Code, Method: me.method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
		} else {
			var envs []errorEnvelope
			if err = json.Unmarshal(data, &envs); err == nil && len(envs) > 0 {
				env = envs[0]
				return nil, Error{Code: env.Error.Code, Method: me.method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
			}
		}
		if len(data) > 0 {
			return nil, fmt.Errorf("status code %d (expected: %d): %s", res.StatusCode, me.expect, string(data))
		}
		return nil, fmt.Errorf("status code %d (expected: %d)", res.StatusCode, me.expect)
	}

	return data, nil
}

func Envelope(data []byte, url string, method string) error {
	if len(data) == 0 {
		return nil
	}
	var err error
	var env errorEnvelope
	if err = json.Unmarshal(data, &env); err == nil && env.Error != nil {
		return Error{Code: env.Error.Code, Method: method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
	} else {
		var envs []errorEnvelope
		if err = json.Unmarshal(data, &envs); err == nil && len(envs) > 0 {
			env = envs[0]
			return Error{Code: env.Error.Code, Method: method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
		}
	}
	return nil
}

const defaultMaxWorkers = 20
const highLimitMaxWorkers = 50

var maxWorkers = resolveMaxWorkers()

func resolveMaxWorkers() int64 {
	sMaxWorkers := os.Getenv("DYNATRACE_MAX_HTTP_WORKERS")
	if len(sMaxWorkers) == 0 {
		return defaultMaxWorkers
	}
	mw, err := strconv.Atoi(sMaxWorkers)
	if err != nil {
		return defaultMaxWorkers
	}
	if mw > highLimitMaxWorkers {
		return highLimitMaxWorkers
	}
	if mw < 1 {
		return 1
	}
	return int64(mw)
}

var sem = semaphore.NewWeighted(maxWorkers)

func (s *legacy_request) execute(ctx context.Context, callback func() (*http.Response, error)) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	err := sem.Acquire(ctx, 1)
	if err != nil {
		return nil, err
	}
	defer sem.Release(1)
	if shutdown.System.Stopped() {
		return nil, nil
	}

	response, err := callback()
	if err != nil {
		return nil, err
	}

	maxIterationCount := 500
	currentIteration := 0

	for response.StatusCode == http.StatusTooManyRequests && currentIteration < maxIterationCount {

		if limit, humanReadableTimestamp, timeInMicroseconds, err := s.extractRateLimitHeaders(response); err == nil {
			Logger.Printf(ctx, "Rate limit of %s requests/min reached (iteration: %d)", limit, currentIteration+1)
			Logger.Printf(ctx, "Attempting to sleep until %s", humanReadableTimestamp)

			now := Now()                                            // client time
			resetTime := MicrosecondsToUnixTime(timeInMicroseconds) // server time
			// mixing server and client time here - sanity check necessary
			sleepDuration := min(max(resetTime.Sub(now), MinWaitTime), MaxWaitTime)

			time.Sleep(sleepDuration)

			currentIteration++
			if response, err = callback(); err != nil {
				return nil, err
			}
		} else {
			// fallback if there are no response headers available
			time.Sleep(30 * time.Second)
			currentIteration++
			if response, err = callback(); err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}

func (s *legacy_request) extractRateLimitHeaders(response *http.Response) (limit string, humanReadableResetTimestamp string, resetTimeInMicroseconds int64, err error) {
	limit = response.Header.Get("X-RateLimit-Limit")
	reset := response.Header.Get("X-RateLimit-Reset")

	if len(limit) == 0 {
		return "", "", 0, errors.New("rate limit header 'X-RateLimit-Limit' not found")
	}
	if len(reset) == 0 {
		return "", "", 0, errors.New("rate limit header 'X-RateLimit-Reset' not found")
	}

	humanReadableResetTimestamp, resetTimeInMicroseconds, err = StringTimestampToHumanReadableFormat(reset)
	if err != nil {
		return "", "", 0, err
	}

	return limit, humanReadableResetTimestamp, resetTimeInMicroseconds, nil
}

func min(a, b time.Duration) time.Duration {
	if a.Nanoseconds() < b.Nanoseconds() {
		return a
	}

	return b
}

func max(a, b time.Duration) time.Duration {
	if a.Nanoseconds() < b.Nanoseconds() {
		return b
	}
	return a
}

func Now() time.Time {
	nowInLocalTimeZone := time.Now()
	location, _ := time.LoadLocation("UTC")
	return nowInLocalTimeZone.In(location)
}

// StringTimestampToHumanReadableFormat parses and sanity-checks a unix timestamp as string and returns it
// as int64 and a human-readable representation of it
func StringTimestampToHumanReadableFormat(unixTimestampAsString string) (humanReadable string, parsedTimestamp int64, err error) {
	if parsedTimestamp, err = strconv.ParseInt(unixTimestampAsString, 10, 64); err != nil {
		return "", 0, fmt.Errorf("%s is not a valid unix timestamp", unixTimestampAsString)
	}
	return time.Unix(parsedTimestamp, 0).Format(time.RFC3339), parsedTimestamp, nil
}

// MicrosecondsToUnixTime converts the UTC time in microseconds to a time.Time struct (unix time)
func MicrosecondsToUnixTime(timeInMicroseconds int64) time.Time {
	return time.Unix(timeInMicroseconds/1000000, (timeInMicroseconds%1000000)*1000)
}
