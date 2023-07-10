/*
 * @license
 * Copyright 2023 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const RESPONSE_HEADER_RATE_LIMIT = "X-RateLimit-Limit"
const RESPONSE_HEADER_RATE_LIMIT_RESET = "X-RateLimit-Reset"

type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
	Pagination Pagination
}

func (resp Response) Success() bool {
	return resp.StatusCode >= 200 && resp.StatusCode <= 299
}

func (resp Response) Is5xxError() bool {
	return resp.StatusCode >= 500 && resp.StatusCode <= 599
}

func (resp Response) Is4xxError() bool {
	return resp.StatusCode >= 400 && resp.StatusCode <= 499
}

func (resp Response) SleepDuration() (sleepDuration time.Duration, err error) {
	_, timeInMicroseconds, err := resp.evalRateLimit()
	if err != nil {
		return 0, fmt.Errorf("encountered response code 'STATUS_TOO_MANY_REQUESTS (429)' but failed to extract rate limit header: %w", err)
	}
	// we're correlating server time and client time here
	return microsecondsToUnixTime(timeInMicroseconds).Sub(time.Now()), nil
}

func (resp Response) evalRateLimit() (limit string, resetTimeInMicroseconds int64, err error) {
	limitAsArray := resp.Headers[http.CanonicalHeaderKey(RESPONSE_HEADER_RATE_LIMIT)]
	resetAsArray := resp.Headers[http.CanonicalHeaderKey(RESPONSE_HEADER_RATE_LIMIT_RESET)]

	if limitAsArray == nil || limitAsArray[0] == "" {
		return "", 0, fmt.Errorf("rate limit header '%s' not found", RESPONSE_HEADER_RATE_LIMIT)
	}
	if resetAsArray == nil || resetAsArray[0] == "" {
		return "", 0, fmt.Errorf("rate limit header '%s' not found", RESPONSE_HEADER_RATE_LIMIT_RESET)
	}

	limit = limitAsArray[0]
	resetTimeInMicroseconds, err = ParseTimeStamp(resetAsArray[0])

	return
}

// ParseTimeStamp parses and sanity-checks a unix timestamp as string and returns it as int64
func ParseTimeStamp(unixTimestampAsString string) (parsedTimestamp int64, err error) {
	if parsedTimestamp, err = strconv.ParseInt(unixTimestampAsString, 10, 64); err != nil {
		return 0, fmt.Errorf("%s is not a valid unix timestamp", unixTimestampAsString)
	}
	return parsedTimestamp, nil
}
