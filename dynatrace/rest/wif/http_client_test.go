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
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// answering starts a server that replies with the given status codes in order, repeating the last one
// once they run out, and returns a retrying client aimed at it together with the request counter.
//
// The backoff is negligible on purpose, so that the retry behaviour can be observed without the test
// waiting for the production one.
func answering(t *testing.T, statusCodes ...int) (*http.Client, string, *atomic.Int64) {
	t.Helper()

	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		attempt := int(requests.Add(1)) - 1
		writer.WriteHeader(statusCodes[min(attempt, len(statusCodes)-1)])
	}))
	t.Cleanup(server.Close)

	client := &http.Client{Transport: &retryTransport{base: http.DefaultTransport, backoff: time.Millisecond}}

	return client, server.URL, &requests
}

func TestRetryTransportRetriesServerError(t *testing.T) {
	client, url, requests := answering(t, http.StatusInternalServerError, http.StatusOK)

	response, err := client.Get(url)
	require.NoError(t, err)
	defer drainAndClose(response.Body)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(2), requests.Load())
}

// The sibling of the test above, differing only in which of the two retriable answers the token
// service gives. Being rate limited is the one a job shares with every other job on the runner fleet.
func TestRetryTransportRetriesTooManyRequests(t *testing.T) {
	client, url, requests := answering(t, http.StatusTooManyRequests, http.StatusOK)

	response, err := client.Get(url)
	require.NoError(t, err)
	defer drainAndClose(response.Body)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(2), requests.Load())
}

func TestRetryTransportDoesNotRetryClientError(t *testing.T) {
	client, url, requests := answering(t, http.StatusForbidden)

	response, err := client.Get(url)
	require.NoError(t, err)
	defer drainAndClose(response.Body)

	assert.Equal(t, http.StatusForbidden, response.StatusCode)
	assert.Equal(t, int64(1), requests.Load())
}

func TestRetryTransportGivesUpAfterMaxAttempts(t *testing.T) {
	client, url, requests := answering(t, http.StatusInternalServerError)

	response, err := client.Get(url)
	require.NoError(t, err)
	defer drainAndClose(response.Body)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, int64(mintAttempts), requests.Load())
}
