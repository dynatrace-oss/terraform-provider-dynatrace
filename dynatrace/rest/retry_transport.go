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
	"bytes"
	"context"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"golang.org/x/oauth2"
)

const (
	defaultMaxRetries  = 5
	defaultBaseBackoff = 5 * time.Second
	defaultMaxBackoff  = 60 * time.Second
)

// RetryTransport is an http.RoundTripper that automatically retries requests
// when the server returns HTTP 429 (Too Many Requests) or HTTP 503 (Service Unavailable).
// MaxRetries, BaseBackoff, and MaxBackoff can be set to non-zero values to override the
// defaults, which is useful in tests to keep execution time short.
type RetryTransport struct {
	Transport   http.RoundTripper
	MaxRetries  int
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
}

// transport returns the configured RoundTripper, falling back to http.DefaultTransport when nil.
func (t *RetryTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *RetryTransport) maxRetries() int {
	if t.MaxRetries > 0 {
		return t.MaxRetries
	}
	return defaultMaxRetries
}

func (t *RetryTransport) baseBackoff() time.Duration {
	if t.BaseBackoff > 0 {
		return t.BaseBackoff
	}
	return defaultBaseBackoff
}

func (t *RetryTransport) maxBackoff() time.Duration {
	if t.MaxBackoff > 0 {
		return t.MaxBackoff
	}
	return defaultMaxBackoff
}

func isRetriableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == http.StatusServiceUnavailable
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

// drainAndClose drains body into /dev/null and closes it, allowing TCP
// connection reuse. The drain error takes precedence over the close error.
func drainAndClose(body io.ReadCloser) error {
	_, copyErr := io.Copy(io.Discard, body)
	closeErr := body.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the request body so it can be replayed on retries.
	bodyBytes, err := bufferRequestBody(req)
	if err != nil {
		return nil, err
	}

	maxRetries := t.maxRetries()

	for attempt := 0; ; attempt++ {
		// Clone the request and restore the body for each attempt.
		clone := req.Clone(req.Context())
		if bodyBytes != nil {
			clone.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err := t.transport().RoundTrip(clone)
		if err != nil {
			return nil, err
		}

		if !isRetriableStatus(resp.StatusCode) || attempt >= maxRetries {
			return resp, nil
		}

		ctx := clone.Context()
		wait := t.sleepDuration(resp, attempt)
		logging.Logger.Printf(ctx, "[RetryTransport] Received HTTP %d, retrying in %s (attempt %d/%d)", resp.StatusCode, wait, attempt+1, maxRetries)
		if err := drainAndClose(resp.Body); err != nil {
			return nil, err
		}
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// getRetryAfterHeaderAsSleepTime parses the Retry-After response header and returns the
// indicated wait duration. Returns 0 if the header is absent,
// cannot be parsed, or indicates a non-positive delay.
func getRetryAfterHeaderAsSleepTime(resp *http.Response) time.Duration {
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
func (t *RetryTransport) sleepDuration(resp *http.Response, attempt int) time.Duration {
	d := getRetryAfterHeaderAsSleepTime(resp)
	if d <= 0 {
		d = computeBackoffSleepTime(t.baseBackoff(), attempt)
	}
	if maxWait := t.maxBackoff(); d > maxWait {
		return maxWait
	}
	return d
}

func NewContextWithOAuthRetryClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Transport: &RetryTransport{},
	})
}
