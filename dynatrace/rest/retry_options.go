/*
 * @license
 * Copyright 2026 Dynatrace LLC
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
	"net/http"
	"time"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

var defaultRetryOptions = &rest.RetryOptions{
	MaxRetries:      30,
	DelayAfterRetry: 10 * time.Second,
	ShouldRetryFunc: RetryIfTemporaryUnavailable,
}

var eofRetryOptions = &rest.RetryOptions{
	MaxRetries:             30,
	DelayAfterRetry:        10 * time.Second,
	ShouldRetryFunc:        RetryIfTemporaryUnavailable,
	ShouldRetryOnErrorFunc: rest.RetryOnUnexpectedEOF,
}

// RetryIfTemporaryUnavailable returns true for responses with status code 429, 502, 503, or 504.
func RetryIfTemporaryUnavailable(resp *http.Response) bool {
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}
