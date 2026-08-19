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

package rest

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type internalRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f internalRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func internalEmptyResponse(status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}
}

func TestRetryTransport_NonIdempotentTransportErrorOptedIn(t *testing.T) {
	calls := 0
	tr := &RetryTransport{
		Transport: internalRoundTripperFunc(func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return nil, io.ErrUnexpectedEOF
			}
			return internalEmptyResponse(http.StatusOK), nil
		}),
		BaseBackoff: time.Millisecond,
	}

	ctx := withNonIdempotentRetry(t.Context())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://dynatrace.com", nil)
	require.NoError(t, err)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
}

func TestRetryTransport_BodyReplayedAfterNonIdempotentTransportError(t *testing.T) {
	const payload = "hello-transport-error-retry"
	var receivedBodies []string

	tr := &RetryTransport{
		Transport: internalRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.Body != nil {
				data, _ := io.ReadAll(req.Body)
				receivedBodies = append(receivedBodies, string(data))
			}
			if len(receivedBodies) < 2 {
				return nil, io.ErrUnexpectedEOF
			}
			return internalEmptyResponse(http.StatusOK), nil
		}),
		BaseBackoff: time.Millisecond,
	}

	ctx := withNonIdempotentRetry(t.Context())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://dynatrace.com", bytes.NewBufferString(payload))
	require.NoError(t, err)
	resp, err := tr.RoundTrip(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, receivedBodies, 2, "body must be delivered on every attempt")
	assert.Equal(t, payload, receivedBodies[0])
	assert.Equal(t, payload, receivedBodies[1])
}
