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
	"context"
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
// than 429 and 503 is returned on the very first attempt.
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
			assert.Equal(t, 1, calls, "expected exactly 1 call for non-retriable status %d", status)
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
				return emptyResponse(http.StatusTooManyRequests, nil), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
		BaseBackoff: 10 * time.Millisecond,
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
				return emptyResponse(http.StatusTooManyRequests, nil), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
		BaseBackoff: 10 * time.Millisecond,
	}

	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, calls)
}

// TestRetryTransport_RetriesOn503_EventualSuccess verifies that a 503 response triggers
// a retry and that the eventual non-503 response is returned to the caller.
// Uses Retry-After: 1 – test duration ~2 seconds.
func TestRetryTransport_RetriesOn503_EventualSuccess(t *testing.T) {
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls < 3 {
				return emptyResponse(http.StatusServiceUnavailable, nil), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
		BaseBackoff: 10 * time.Millisecond,
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
// present the transport falls back to exponential backoff (base × 1.5^attempt).
func TestRetryTransport_ExponentialBackoff(t *testing.T) {
	const base = 100 * time.Millisecond
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
		BaseBackoff: base,
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	// base × 1.5^0 = base for the first retry.
	assert.GreaterOrEqual(t, elapsed, base, "expected exponential backoff of at least %s", base)
}

// TestRetryTransport_MaxRetriesExhausted verifies that the transport stops after
// MaxRetries attempts and returns the final retriable response to the caller.
// Uses Retry-After: 1 – test duration ~10 seconds (10 retries × 1 s each).
func TestRetryTransport_MaxRetriesExhausted(t *testing.T) {
	const maxRetries = 3
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			return emptyResponse(http.StatusTooManyRequests, nil), nil
		}),
		MaxRetries:  maxRetries,
		BaseBackoff: 10 * time.Millisecond,
	}

	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	// MaxRetries = 3, so the loop runs for attempts 0..3 → 4 total calls.
	assert.Equal(t, maxRetries+1, calls, "expected MaxRetries+1 = %d total attempts", maxRetries+1)
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

// TestRetryTransport_RetryAfterZero verifies that a Retry-After value of 0 is treated as
// invalid and the transport falls back to exponential back-off instead.
func TestRetryTransport_RetryAfterZero_FallsBackToBackoff(t *testing.T) {
	const base = 30 * time.Millisecond
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": "0"}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
		BaseBackoff: base,
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	assert.GreaterOrEqual(t, elapsed, base, "Retry-After: 0 should fall back to exponential backoff")
}

// TestRetryTransport_RetryAfterNegative verifies that a negative Retry-After value is
// treated as invalid and the transport falls back to exponential back-off instead.
func TestRetryTransport_RetryAfterNegative_FallsBackToBackoff(t *testing.T) {
	const base = 30 * time.Millisecond
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": "-5"}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
		BaseBackoff: base,
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	assert.GreaterOrEqual(t, elapsed, base, "negative Retry-After should fall back to exponential backoff")
}

// TestRetryTransport_RetryAfterPastDate verifies that a Retry-After HTTP-date in the past
// is treated as invalid and the transport falls back to exponential back-off instead.
func TestRetryTransport_RetryAfterPastDate_FallsBackToBackoff(t *testing.T) {
	const base = 30 * time.Millisecond
	pastDate := time.Now().Add(-10 * time.Second).UTC().Format(http.TimeFormat)
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": pastDate}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
		BaseBackoff: base,
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	assert.GreaterOrEqual(t, elapsed, base, "past Retry-After date should fall back to exponential backoff")
}

// TestRetryTransport_RetryAfterCapped verifies that an excessively large Retry-After value
// is capped at MaxBackoff rather than waited in full.
func TestRetryTransport_RetryAfterCapped(t *testing.T) {
	const cap = 50 * time.Millisecond
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return emptyResponse(http.StatusTooManyRequests, map[string]string{"Retry-After": "99999"}), nil
			}
			return emptyResponse(http.StatusOK, nil), nil
		}),
		MaxBackoff: cap,
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
	assert.Less(t, elapsed, time.Second, "Retry-After of 99999 s should have been capped")
	assert.GreaterOrEqual(t, elapsed, cap, "should have waited at least MaxBackoff")
}

// TestRetryTransport_BackoffCapped verifies that the exponential back-off is capped at
// MaxBackoff when the computed value would otherwise exceed it.
func TestRetryTransport_BackoffCapped(t *testing.T) {
	const cap = 50 * time.Millisecond
	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			return emptyResponse(http.StatusTooManyRequests, nil), nil
		}),
		MaxRetries:  3,
		BaseBackoff: cap * 10, // base far exceeds cap → every wait should be capped
		MaxBackoff:  cap,
	}

	start := time.Now()
	req, _ := http.NewRequest(http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, 4, calls) // MaxRetries+1
	// 3 waits each capped at ~55 ms (cap + ≤10 % jitter) → well under 1 s.
	assert.Less(t, elapsed, time.Second, "backoff should have been capped at MaxBackoff")
}

// TestRetryTransport_ContextCancelled verifies that a context cancellation during the
// retry wait is honoured immediately and the error is propagated to the caller.
func TestRetryTransport_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	calls := 0
	tr := &rest.RetryTransport{
		Transport: rtFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			// Cancel the context so the upcoming retry wait fires the Done branch.
			cancel()
			return emptyResponse(http.StatusTooManyRequests, nil), nil
		}),
		BaseBackoff: time.Hour, // large enough that only context cancellation can unblock
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://dynatrace.com", nil)
	resp, err := tr.RoundTrip(req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 1, calls, "should not retry after context is cancelled")
}

// failReader is an io.Reader that always returns an error.
type failReader struct{ err error }

func (f *failReader) Read(_ []byte) (int, error) { return 0, f.err }
