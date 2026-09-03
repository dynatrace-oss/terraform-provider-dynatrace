/**
* @license
* Copyright 2026 Dynatrace LLC
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

package wif

import (
	"net/http"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/retry"
)

// Covers the retries too: the client timeout cancels the request context, which also ends the wait
// between attempts.
const mintTimeout = 30 * time.Second

// Must not use http.DefaultTransport: the provider mutates that one when DYNATRACE_HTTP_INSECURE is
// set, and that escape hatch must not disable certificate verification on the call carrying the
// minting credentials.
func mintingHTTPClient() *http.Client {
	return &http.Client{
		Timeout: mintTimeout,
		Transport: &retry.Transport{
			Base:        &http.Transport{Proxy: http.ProxyFromEnvironment},
			MaxRetries:  2,
			BaseBackoff: 500 * time.Millisecond,
			MaxBackoff:  5 * time.Second,
			ShouldRetry: retriableMintFailure,
		},
	}
}

// Wider than the shared default, which stops at 429 and 503.
func retriableMintFailure(response *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= http.StatusInternalServerError
}
