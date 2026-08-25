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

package retry

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
)

const (
	defaultMaxRetries  = 5
	defaultBaseBackoff = 5 * time.Second
	defaultMaxBackoff  = 60 * time.Second
)

// Transport is an http.RoundTripper that automatically repeats a request whose outcome ShouldRetry
// accepts, by default when the server returns HTTP 429 (Too Many Requests) or HTTP 503 (Service
// Unavailable). MaxRetries, BaseBackoff, and MaxBackoff can be set to non-zero values to override
// the defaults, which is useful in tests to keep execution time short.
type Transport struct {
	Base        http.RoundTripper
	MaxRetries  int
	BaseBackoff time.Duration
	MaxBackoff  time.Duration

	// ShouldRetry decides whether an attempt is worth repeating. It receives the response, or nil
	// and the error when the round trip itself did not complete. Defaults to
	// IfTooManyRequestsOrServiceUnavailable.
	ShouldRetry func(*http.Response, error) bool
}

// IfTooManyRequestsOrServiceUnavailable repeats a request the server answered with HTTP 429 or 503,
// and gives up on anything else, including a round trip that failed outright.
func IfTooManyRequestsOrServiceUnavailable(resp *http.Response, err error) bool {
	return err == nil && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable)
}

// base returns the configured RoundTripper, falling back to http.DefaultTransport when nil.
func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func (t *Transport) shouldRetry(resp *http.Response, err error) bool {
	if t.ShouldRetry != nil {
		return t.ShouldRetry(resp, err)
	}
	return IfTooManyRequestsOrServiceUnavailable(resp, err)
}

func (t *Transport) maxRetries() int {
	if t.MaxRetries > 0 {
		return t.MaxRetries
	}
	return defaultMaxRetries
}

func (t *Transport) baseBackoff() time.Duration {
	if t.BaseBackoff > 0 {
		return t.BaseBackoff
	}
	return defaultBaseBackoff
}

func (t *Transport) maxBackoff() time.Duration {
	if t.MaxBackoff > 0 {
		return t.MaxBackoff
	}
	return defaultMaxBackoff
}

// bufferRequestBody reads and closes req.Body, returning its contents.
// Returns nil, nil when the body is absent.
func bufferRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(req.Body)
	if closeErr := req.Body.Close(); err == nil {
		err = closeErr
	}
	return body, err
}

// cloneRequestWithBody returns a copy of req that can be sent on its own, with the previously
// buffered body put back in place of the one the last attempt consumed.
func cloneRequestWithBody(req *http.Request, body []byte) *http.Request {
	clone := req.Clone(req.Context())
	if body != nil {
		clone.Body = io.NopCloser(bytes.NewReader(body))
	}
	return clone
}

// DrainAndClose drains body into /dev/null and closes it, allowing TCP
// connection reuse. The drain error takes precedence over the close error.
func DrainAndClose(body io.ReadCloser) error {
	_, copyErr := io.Copy(io.Discard, body)
	closeErr := body.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the request body so it can be replayed on retries.
	bodyBytes, err := bufferRequestBody(req)
	if err != nil {
		return nil, err
	}

	maxRetries := t.maxRetries()

	for attempt := 0; ; attempt++ {
		clone := cloneRequestWithBody(req, bodyBytes)

		resp, err := t.base().RoundTrip(clone)
		if !t.shouldRetry(resp, err) || attempt >= maxRetries {
			return resp, err
		}

		ctx := clone.Context()
		wait := t.sleepDuration(resp, attempt)
		logging.Logger.Printf(ctx, "[RetryTransport] %s, retrying in %s (attempt %d/%d)", retryReason(resp, err), wait, attempt+1, maxRetries)
		if resp != nil {
			if err := DrainAndClose(resp.Body); err != nil {
				return nil, err
			}
		}
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// retryReason describes the outcome that ShouldRetry accepted, for the log line.
func retryReason(resp *http.Response, err error) string {
	if err != nil {
		return fmt.Sprintf("Request failed (%s)", err)
	}
	return fmt.Sprintf("Received HTTP %d", resp.StatusCode)
}

// getRetryAfterHeaderAsSleepTime parses the Retry-After response header and returns the
// indicated wait duration. Returns 0 if there is no response, or if the header is absent,
// cannot be parsed, or indicates a non-positive delay.
func getRetryAfterHeaderAsSleepTime(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	ra := resp.Header.Get("Retry-After")
	if ra == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(ra); err == nil && seconds > 0 {
		d := time.Duration(seconds) * time.Second
		return d
	}
	if t, err := http.ParseTime(ra); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// computeBackoffSleepTime returns the exponential back-off duration for the given attempt
// using a 1.5× multiplier (base × 1.5^attempt), with up to 10%
// random jitter to avoid thundering-herd problems.
func computeBackoffSleepTime(base time.Duration, attempt int) time.Duration {
	// random number between 1.0 and 1.1 to add up to 10% jitter
	scale := 1 + 0.1*rand.Float64()
	return time.Duration(float64(base) * scale * math.Pow(1.5, float64(attempt)))
}

// sleepDuration determines how long to wait before the next retry attempt.
// It honours the Retry-After response header when present (capped at maxBackoff);
// otherwise it falls back to exponential back-off with jitter.
func (t *Transport) sleepDuration(resp *http.Response, attempt int) time.Duration {
	d := getRetryAfterHeaderAsSleepTime(resp)
	if d <= 0 {
		d = computeBackoffSleepTime(t.baseBackoff(), attempt)
	}
	if maxWait := t.maxBackoff(); d > maxWait {
		return maxWait
	}
	return d
}
