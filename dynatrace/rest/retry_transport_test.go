//go:build unit

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

package rest_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// rtFunc lets an ordinary function act as an http.RoundTripper.
type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// emptyResponse returns a minimal *http.Response with the given status and optional headers.
func emptyResponse(status int, headers map[string]string) *http.Response {
	h := make(http.Header)
	for k, v := range headers {
		h.Set(k, v)
	}
	return &http.Response{
		StatusCode: status,
		Header:     h,
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}
}

// TestRetryTransport_SuccessOnFirstAttempt verifies that a successful response is
// returned immediately without any retry.
func TestRetryTransport_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			return emptyResponse(http.StatusOK, nil), nil
		}),
	}

	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, calls)
}

// TestRetryTransport_NonRetriableStatusCodes verifies that every status code other
// than 429 is returned on the very first attempt.
func TestRetryTransport_NonRetriableStatusCodes(t *testing.T) {
	nonRetryStatusCodes := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
	}
	for _, status := range nonRetryStatusCodes {
		t.Run(fmt.Sprintf("status_%d", status), func(t *testing.T) {
			calls := 0
			tr := &rest.RetryTransport{
				Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
					calls++
					return emptyResponse(status, nil), nil
				}),
			}

			req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
			resp, err := tr.RoundTrip(req)

			require.NoError(t, err)
			assert.Equal(t, status, resp.StatusCode)
			assert.Equal(t, 1, calls, "expected exactly 1 call for non-429 status %d", status)
		})
	}
}

// TestRetryTransport_TransportError verifies that an error from the underlying
// transport is propagated immediately without retrying.
func TestRetryTransport_TransportError(t *testing.T) {
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			return nil, assert.AnError
		}),
	}

	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, assert.AnError)
}

// TestRetryTransport_NilBody verifies that requests without a body are handled correctly.
func TestRetryTransport_NilBody(t *testing.T) {
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			return emptyResponse(http.StatusOK, nil), nil
		}),
	}

	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestRetryTransport_BodyReadError verifies that an error while buffering the request
// body is returned immediately without calling the underlying transport.
func TestRetryTransport_BodyReadError(t *testing.T) {
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			t.Fatal("transport should not be called when the body cannot be read")
			return nil, nil
		}),
	}

	req, _ := http.NewRequest(http.MethodPost, "https://dynatrace.com", &failReader{err: assert.AnError})
	resp, err := tr.RoundTrip(req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, assert.AnError)
}

// TestRetryTransport_BodyReplayed verifies that the request body is replayed in full
// on every retry attempt so that the upstream server always receives the same payload.
// Uses Retry-After: 1 to keep the test duration to ~1 second.
func TestRetryTransport_BodyReplayed(t *testing.T) {
	const payload = "hello-retry"
	var receivedBodies []string

	tr := &rest.RetryTransport{
		Transport: rtFunc(func(req *http.Request) (*http.Response, error) {
			if req.Body != nil {
				data, _ := io.ReadAll(req.Body)
				receivedBodies = append(receivedBodies, string(data))
			}
			if len(receivedBodies) < 2 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": "1"}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
	}

	req, _ := http.NewRequest(http.MethodPost, "https://dynatrace.com", bytes.NewBufferString(payload))
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, receivedBodies, 2, "body must be delivered on every attempt")
	assert.Equal(t, payload, receivedBodies[0])
	assert.Equal(t, payload, receivedBodies[1])
}

// TestRetryTransport_RetriesOn429_EventualSuccess verifies that a 429 response triggers
// a retry and that the eventual non-429 response is returned to the caller.
// Uses Retry-After: 1 – test duration ~2 seconds.
func TestRetryTransport_RetriesOn429_EventualSuccess(t *testing.T) {
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls < 3 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": "1"}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
	}

	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, calls)
}

// TestRetryTransport_RetryAfterSeconds verifies that an integer Retry-After header value
// is honored as a wait in seconds.
// Test duration ~1 second.
func TestRetryTransport_RetryAfterSeconds(t *testing.T) {
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": "1"}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	assert.GreaterOrEqual(t, elapsed, time.Second, "should have waited at least 1 s per Retry-After header")
}

// TestRetryTransport_RetryAfterHTTPDate verifies that an HTTP-date Retry-After header
// value is parsed and honoured as a wait duration.
// Test duration ~2 seconds.
func TestRetryTransport_RetryAfterHTTPDate(t *testing.T) {
	retryAfterDate := time.Now().Add(2 * time.Second).UTC().Format(http.TimeFormat)

	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": retryAfterDate}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	assert.GreaterOrEqual(t, elapsed, time.Second, "should have waited for the Retry-After date")
}

// TestRetryTransport_ExponentialBackoff verifies that when no Retry-After header is
// present the transport falls back to exponential backoff (5 s × 2^attempt).
// Test duration ~5 seconds (defaultBackoff × 2^0 for the first attempt).
func TestRetryTransport_ExponentialBackoff(t *testing.T) {
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				// No Retry-After header → exponential backoff applies.
				return emptyResponse(http.StatusTooManyRequests, nil), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	// defaultBackoff (5 s) × 2^0 = 5 s for the first retry.
	assert.GreaterOrEqual(t, elapsed, 5*time.Second, "expected exponential backoff of at least 5 s")
}

// TestRetryTransport_MaxRetriesExhausted verifies that the transport stops after
// maxRetries attempts and returns the final 429 response to the caller.
// Uses Retry-After: 1 – test duration ~10 seconds (10 retries × 1 s each).
func TestRetryTransport_MaxRetriesExhausted(t *testing.T) {
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": "1"}), nil
		}),
	}

	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	// maxRetries = 10, so the loop runs for attempts 0..10 → 11 total calls.
	assert.Equal(t, 11, calls, "expected maxRetries+1 = 11 total attempts")
}

// NewContextWithOAuthRetryClient verifies that the context is enriched with an
// HTTP client whose transport is a *RetryTransport.
func TestNewContextWithOAuthRetryClient(t *testing.T) {
	ctx := rest.NewContextWithOAuthRetryClient(t.Context())

	client, ok := ctx.Value(oauth2.HTTPClient).(*http.Client)
	require.True(t, ok, "context should contain an *http.Client")
	_, isRetry := client.Transport.(*rest.RetryTransport)
	assert.True(t, isRetry, "the client's transport should be a *RetryTransport")
}

// failReader is an io.Reader that always returns an error.
type failReader struct{ err error }

func (f *failReader) Read(_ []byte) (int, error) { return 0, f.err }
