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

const (
	// Bounds one attempt to obtain a token including its retries, because the client timeout cancels
	// the request context, which also ends the wait between attempts.
	mintTimeout = 30 * time.Second

	mintRetries     = 2
	mintBaseBackoff = 500 * time.Millisecond

	// Waiting longer than this only burns mintTimeout without leaving room for another attempt.
	mintMaxBackoff = 5 * time.Second
)

// mintingHTTPClient builds the client used to talk to a vendor's token service. It must not use
// http.DefaultTransport: the provider mutates that one when DYNATRACE_HTTP_INSECURE is set, and an
// escape hatch for self-signed Managed clusters must not disable certificate verification on the
// call carrying the minting credentials.
func mintingHTTPClient() *http.Client {
	return &http.Client{
		Timeout: mintTimeout,
		Transport: &retry.Transport{
			Base:        &http.Transport{Proxy: http.ProxyFromEnvironment},
			MaxRetries:  mintRetries,
			BaseBackoff: mintBaseBackoff,
			MaxBackoff:  mintMaxBackoff,
			ShouldRetry: retriableMintFailure,
		},
	}
}

// retriableMintFailure widens the shared default, which only covers 429 and 503: a token service is
// worth asking again after any server error, and after a round trip that never completed at all.
// Both are ordinary flakiness on a CI runner, and the request carries nothing that could be applied
// twice.
func retriableMintFailure(response *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= http.StatusInternalServerError
}
