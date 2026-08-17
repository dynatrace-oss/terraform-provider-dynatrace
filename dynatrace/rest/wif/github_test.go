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
	"errors"
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

// A stub rather than a closed port, so the asserted error is the one given to the transport rather
// than however the operating system words a refused connection.
type unreachableTransport struct {
	err error
}

func (transport *unreachableTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, transport.err
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

func TestGitHubMinterReportsUnreachableTokenService(t *testing.T) {
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "https://token.service.invalid/")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "request-token")
	client := &http.Client{Transport: &unreachableTransport{err: errors.New("no route to the token service")}}
	minter, err := newGitHubMinter("dynatrace", client)
	require.NoError(t, err)

	_, err = minter.mint(t.Context())

	assert.EqualError(t, err, `failed to get ID token: Get "https://token.service.invalid/?audience=dynatrace": no route to the token service`)
}

func TestGitHubMinterRejectsUnsuccessfulResponse(t *testing.T) {
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
		_, _ = writer.Write([]byte(`{"message":"forbidden"}`))
	})

	_, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	assert.EqualError(t, err, "failed to get ID token: ACTIONS_ID_TOKEN_REQUEST_URL responded with HTTP 403")
}

func TestGitHubMinterRejectsResponseThatIsNotJSON(t *testing.T) {
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("not JSON"))
	})

	_, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	assert.EqualError(t, err, "failed to get ID token: the response of the token service is not valid JSON")
}

func TestGitHubMinterRejectsResponseWithoutToken(t *testing.T) {
	server := tokenServiceHandler(t, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"count":0}`))
	})

	_, err := gitHubMinterFor(t, server.URL, "dynatrace").mint(t.Context())

	assert.EqualError(t, err, "failed to get ID token: the response of the token service does not contain a token")
}

func TestGitHubMinterRequiresRequestURLVariable(t *testing.T) {
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "request-token")

	_, err := newGitHubMinter("dynatrace", http.DefaultClient)

	assert.EqualError(t, err, "Unable to get ACTIONS_ID_TOKEN_REQUEST_URL environment variable. `wif_vendor = \"github\"` only works inside a GitHub Actions job that is allowed to request an OIDC token. Add `permissions: { id-token: write }` to the workflow or to the job.")
}

func TestGitHubMinterRequiresRequestTokenVariable(t *testing.T) {
	t.Setenv(envutils.ActionsIDTokenRequestURL.Key, "https://token.service.invalid/")
	t.Setenv(envutils.ActionsIDTokenRequestToken.Key, "")

	_, err := newGitHubMinter("dynatrace", http.DefaultClient)

	assert.EqualError(t, err, "Unable to get ACTIONS_ID_TOKEN_REQUEST_TOKEN environment variable. `wif_vendor = \"github\"` only works inside a GitHub Actions job that is allowed to request an OIDC token. Add `permissions: { id-token: write }` to the workflow or to the job.")
}
