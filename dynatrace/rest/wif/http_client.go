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
	// mintTimeout bounds a single attempt to obtain a token, including its retries. The client
	// timeout cancels the request context, which also ends the wait between attempts.
	mintTimeout = 30 * time.Second

	mintAttempts     = 3
	mintRetryBackoff = 500 * time.Millisecond
)

// mintingHTTPClient builds the client used to talk to a vendor's token service.
//
// It deliberately does not use http.DefaultTransport. The provider mutates that transport when
// DYNATRACE_HTTP_INSECURE is set, and an escape hatch meant for self-signed Dynatrace Managed
// clusters must not disable certificate verification on the call that carries the credentials used
// to obtain an OIDC token.
func mintingHTTPClient() *http.Client {
	return &http.Client{
		Timeout: mintTimeout,
		Transport: &retryTransport{
			// ProxyFromEnvironment keeps the client working on runners that reach the internet
			// through a proxy.
			base:    &http.Transport{Proxy: http.ProxyFromEnvironment},
			backoff: mintRetryBackoff,
		},
	}
}

// retryTransport retries a token request that failed with a transient server error.
//
// It is a local implementation rather than the provider's RetryTransport because this package must
// not import the rest package, which imports this one. It is also much narrower: requests to a token
// service are bodyless reads, so there is no request body to buffer and replay.
type retryTransport struct {
	base http.RoundTripper
	// backoff is the wait before the second attempt, doubled for each further one.
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

// drainAndClose drains a response body and closes it, so the connection can be reused. Any error is
// discarded on purpose: the caller is either already returning an error or has what it needs.
func drainAndClose(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}
