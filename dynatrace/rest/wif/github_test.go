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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tokenServiceHandler is the shape of the GitHub Actions token service: it answers with the ID token
// wrapped in a "value" field.
func tokenServiceHandler(t *testing.T, respond func(writer http.ResponseWriter, request *http.Request)) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(respond))
	t.Cleanup(server.Close)

	return server
}

func gitHubMinterFor(t *testing.T, requestURL string, audience string) minter {
	t.Helper()

	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, requestURL)
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "request-token")

	minter, err := newGitHubMinter(audience, http.DefaultClient)
	require.NoError(t, err)

	return minter
}

func TestGitHubMinterReturnsTokenFromTokenService(t *testing.T) {
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"count":1,"value":"the.id.token"}`))
	})

	token, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	require.NoError(t, err)
	assert.Equal(t, "the.id.token", token)
}

func TestGitHubMinterRequestsConfiguredAudience(t *testing.T) {
	var audience string
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		audience = request.URL.Query().Get("audience")
		_, _ = writer.Write([]byte(`{"value":"the.id.token"}`))
	})

	_, err := gitHubMinterFor(t, server.URL, "https://dynatrace.com").mint(t.Context())

	require.NoError(t, err)
	assert.Equal(t, "https://dynatrace.com", audience)
}

func TestGitHubMinterAuthenticatesWithRequestToken(t *testing.T) {
	var authorization string
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		authorization = request.Header.Get("Authorization")
		_, _ = writer.Write([]byte(`{"value":"the.id.token"}`))
	})

	_, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	require.NoError(t, err)
	assert.Equal(t, "Bearer request-token", authorization)
}

// The URL GitHub injects already carries an api-version parameter.
func TestGitHubMinterKeepsExistingQueryParameters(t *testing.T) {
	var query string
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		query = request.URL.RawQuery
		_, _ = writer.Write([]byte(`{"value":"the.id.token"}`))
	})

	_, err := gitHubMinterFor(t, server.URL+"/?api-version=2.0", "dynatrace").mint(t.Context())

	require.NoError(t, err)
	assert.Equal(t, "api-version=2.0&audience=dynatrace", query)
}

// Asserting on the cause rather than the message is what proves the round trip failure is wrapped:
// swapping the %w for a %v produces an identical string but breaks this.
func TestGitHubMinterPropagatesCancellation(t *testing.T) {
	cancelled, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := gitHubMinterFor(t, "https://token.service.invalid/", "dynatrace").mint(cancelled)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestGitHubMinterRejectsUnsuccessfulResponse(t *testing.T) {
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
		_, _ = writer.Write([]byte(`{"message":"forbidden"}`))
	})

	_, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	assert.EqualError(t, err, "failed to get ID token: the token service responded with HTTP 403")
}

func TestGitHubMinterRejectsResponseThatIsNotJSON(t *testing.T) {
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("not JSON"))
	})

	_, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	assert.ErrorIs(t, err, errResponseNotJSON)
}

func TestGitHubMinterRejectsResponseWithoutToken(t *testing.T) {
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"count":0}`))
	})

	_, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	assert.ErrorIs(t, err, errResponseWithoutToken)
}

// The parse happens once, at construction, so a malformed URL is reported alongside the missing one
// rather than on every mint. The wrapped cause carries the offending value.
func TestGitHubMinterRejectsRequestURLThatIsNotAURL(t *testing.T) {
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "://")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "request-token")

	_, err := newGitHubMinter("dynatrace", http.DefaultClient)

	assert.EqualError(t, err, `ACTIONS_ID_TOKEN_REQUEST_URL does not hold a valid URL: parse "://": missing protocol scheme`)
}

func TestGitHubMinterRequiresRequestURLVariable(t *testing.T) {
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "request-token")

	_, err := newGitHubMinter("dynatrace", http.DefaultClient)

	assert.EqualError(t, err, "unable to get ACTIONS_ID_TOKEN_REQUEST_URL environment variable: an OIDC token can only be requested from a GitHub Actions job with `permissions: { id-token: write }`")
}

func TestGitHubMinterRequiresRequestTokenVariable(t *testing.T) {
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "https://token.service.invalid/")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "")

	_, err := newGitHubMinter("dynatrace", http.DefaultClient)

	assert.EqualError(t, err, "unable to get ACTIONS_ID_TOKEN_REQUEST_TOKEN environment variable: an OIDC token can only be requested from a GitHub Actions job with `permissions: { id-token: write }`")
}
