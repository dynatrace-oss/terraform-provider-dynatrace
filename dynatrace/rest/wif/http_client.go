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
	"io"
	"net/http"
	"time"
)

const (
	// Bounds one attempt to obtain a token including its retries, because the client timeout cancels
	// the request context, which also ends the wait between attempts.
	mintTimeout = 30 * time.Second

	mintAttempts     = 3
	mintRetryBackoff = 500 * time.Millisecond
)

// mintingHTTPClient builds the client used to talk to a vendor's token service. It must not use
// http.DefaultTransport: the provider mutates that one when DYNATRACE_HTTP_INSECURE is set, and an
// escape hatch for self-signed Managed clusters must not disable certificate verification on the
// call carrying the minting credentials.
func mintingHTTPClient() *http.Client {
	return &http.Client{
		Timeout: mintTimeout,
		Transport: &retryTransport{
			base:    &http.Transport{Proxy: http.ProxyFromEnvironment},
			backoff: mintRetryBackoff,
		},
	}
}

// retryTransport is local rather than the provider's RetryTransport because this package must not
// import the rest package. Requests to a token service are bodyless, so there is nothing to replay.
type retryTransport struct {
	base    http.RoundTripper
	backoff time.Duration
}

func (transport *retryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	for attempt := 0; ; attempt++ {
		response, err := transport.base.RoundTrip(request.Clone(request.Context()))
		if err == nil && !isRetriableStatus(response.StatusCode) {
			return response, nil
		}
		if attempt >= mintAttempts-1 {
			return response, err
		}

		if response != nil {
			drainAndClose(response.Body)
		}

		select {
		case <-time.After(transport.backoff << attempt):
		case <-request.Context().Done():
			return nil, request.Context().Err()
		}
	}
}

func isRetriableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

// drainAndClose drains a response body and closes it, so the connection can be reused. Errors are
// discarded: the caller is either already returning one or has what it needs.
func drainAndClose(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}
