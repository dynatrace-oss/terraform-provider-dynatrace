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

package wif

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetriableMintFailureByStatus(t *testing.T) {
	retriableByStatus := map[int]bool{
		http.StatusOK:                  false,
		http.StatusBadRequest:          false,
		http.StatusUnauthorized:        false,
		http.StatusForbidden:           false,
		http.StatusNotFound:            false,
		http.StatusTooManyRequests:     true,
		http.StatusInternalServerError: true,
		http.StatusBadGateway:          true,
		http.StatusServiceUnavailable:  true,
		http.StatusGatewayTimeout:      true,
	}

	for status, retriable := range retriableByStatus {
		t.Run(fmt.Sprintf("status_%d", status), func(t *testing.T) {
			assert.Equal(t, retriable, retriableMintFailure(&http.Response{StatusCode: status}, nil))
		})
	}
}

func TestRetriableMintFailureOnFailedRoundTrip(t *testing.T) {
	assert.True(t, retriableMintFailure(nil, assert.AnError))
}

// TestMintingHTTPClientRetriesServerError verifies that retriableMintFailure reaches the transport:
// HTTP 500 is a status the shared default would have handed back on the first attempt.
func TestMintingHTTPClientRetriesServerError(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		if requests.Add(1) == 1 {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	response, err := mintingHTTPClient().Get(server.URL)
	require.NoError(t, err)
	defer func() { _ = response.Body.Close() }()

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(2), requests.Load())
}
