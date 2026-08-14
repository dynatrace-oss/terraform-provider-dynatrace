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
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// roundTripperFunc lets an ordinary function act as an http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (roundTrip roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTrip(request)
}

// platformResponse builds the minimal response the transport looks at.
func platformResponse(statusCode int) *http.Response {
	return &http.Response{StatusCode: statusCode, Header: make(http.Header), Body: io.NopCloser(bytes.NewReader(nil))}
}

// replacingTokenSource hands out the given tokens in order, one per rejection.
type replacingTokenSource struct {
	tokens        []string
	replacements  int
	lastRejection string
}

func (source *replacingTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: source.tokens[source.replacements]}, nil
}

func (source *replacingTokenSource) Replace(rejected string) (*oauth2.Token, error) {
	source.lastRejection = rejected
	source.replacements++
	return source.Token()
}

// irreplaceableTokenSource stands in for a token the provider did not mint: it answers a rejection
// with the very token that was rejected.
type irreplaceableTokenSource struct {
	token string
}

func (source *irreplaceableTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: source.token}, nil
}

func (source *irreplaceableTokenSource) Replace(string) (*oauth2.Token, error) {
	return source.Token()
}

// failingTokenSource stands in for a token service that stopped answering.
type failingTokenSource struct {
	err error
}

func (source *failingTokenSource) Token() (*oauth2.Token, error) {
	return nil, source.err
}

func (source *failingTokenSource) Replace(string) (*oauth2.Token, error) {
	return nil, source.err
}

// authenticatedRequest builds the request the transport is handed, which the oauth2 transport above
// it has already put a token on.
func authenticatedRequest(t *testing.T, token string, body string) *http.Request {
	t.Helper()

	var payload io.Reader
	if len(body) > 0 {
		payload = strings.NewReader(body)
	}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "https://environment.apps.dynatrace.com/platform/management/v1/environment", payload)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+token)

	return request
}

func TestReauthenticatingTransportPassesOnAcceptedRequest(t *testing.T) {
	attempts := 0
	transport := &ReauthenticatingTransport{
		TokenSource: &replacingTokenSource{tokens: []string{"first-token", "second-token"}},
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			attempts++
			return platformResponse(http.StatusOK), nil
		}),
	}

	response, err := transport.RoundTrip(authenticatedRequest(t, "first-token", ""))

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, attempts)
}

// The point of the whole transport: a token that the platform no longer accepts costs one more round
// trip rather than the resource it was sent for.
func TestReauthenticatingTransportRetriesRejectedRequestWithReplacementToken(t *testing.T) {
	var sentTokens []string
	transport := &ReauthenticatingTransport{
		TokenSource: &replacingTokenSource{tokens: []string{"first-token", "second-token"}},
		Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
			sentTokens = append(sentTokens, request.Header.Get("Authorization"))
			if len(sentTokens) == 1 {
				return platformResponse(http.StatusUnauthorized), nil
			}
			return platformResponse(http.StatusOK), nil
		}),
	}

	response, err := transport.RoundTrip(authenticatedRequest(t, "first-token", ""))

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, []string{"Bearer first-token", "Bearer second-token"}, sentTokens)
}

// Which token was turned down decides whether a replacement is minted at all, so the token that was
// on the request has to be the one reported - not whatever the source happens to hold by then.
func TestReauthenticatingTransportReportsTokenThatWasRejected(t *testing.T) {
	tokenSource := &replacingTokenSource{tokens: []string{"first-token", "second-token"}}
	transport := &ReauthenticatingTransport{
		TokenSource: tokenSource,
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return platformResponse(http.StatusUnauthorized), nil
		}),
	}

	_, err := transport.RoundTrip(authenticatedRequest(t, "first-token", ""))

	require.NoError(t, err)
	assert.Equal(t, "first-token", tokenSource.lastRejection)
}

// Sending the request consumes its body, so a retry that does not put it back would reach the
// platform as an empty write.
func TestReauthenticatingTransportRepeatsRequestBody(t *testing.T) {
	var sentBodies []string
	transport := &ReauthenticatingTransport{
		TokenSource: &replacingTokenSource{tokens: []string{"first-token", "second-token"}},
		Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			sentBodies = append(sentBodies, string(body))
			if len(sentBodies) == 1 {
				return platformResponse(http.StatusUnauthorized), nil
			}
			return platformResponse(http.StatusOK), nil
		}),
	}

	_, err := transport.RoundTrip(authenticatedRequest(t, "first-token", `{"name":"a setting"}`))

	require.NoError(t, err)
	assert.Equal(t, []string{`{"name":"a setting"}`, `{"name":"a setting"}`}, sentBodies)
}

// A token that cannot be replaced would be sent again unchanged, and the platform would turn it down
// again for the same reason.
func TestReauthenticatingTransportDoesNotRetryWithUnchangedToken(t *testing.T) {
	attempts := 0
	transport := &ReauthenticatingTransport{
		TokenSource: &irreplaceableTokenSource{token: "supplied-token"},
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			attempts++
			return platformResponse(http.StatusUnauthorized), nil
		}),
	}

	response, err := transport.RoundTrip(authenticatedRequest(t, "supplied-token", ""))

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, 1, attempts)
}

// What the caller is entitled to hear about is the platform turning its credential down. Reporting
// the failed attempt at a new token instead would replace an answer that names the actual problem.
func TestReauthenticatingTransportReportsRejectionWhenReplacementFails(t *testing.T) {
	transport := &ReauthenticatingTransport{
		TokenSource: &failingTokenSource{err: errors.New("the token service is unreachable")},
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return platformResponse(http.StatusUnauthorized), nil
		}),
	}

	response, err := transport.RoundTrip(authenticatedRequest(t, "first-token", ""))

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

// Only an unaccepted token is worth another attempt. A permission the identity does not have is not
// going to be granted by a token minted a second later.
func TestReauthenticatingTransportDoesNotRetryForbiddenRequest(t *testing.T) {
	attempts := 0
	transport := &ReauthenticatingTransport{
		TokenSource: &replacingTokenSource{tokens: []string{"first-token", "second-token"}},
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			attempts++
			return platformResponse(http.StatusForbidden), nil
		}),
	}

	response, err := transport.RoundTrip(authenticatedRequest(t, "first-token", ""))

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
	assert.Equal(t, 1, attempts)
}

func TestReauthenticatingTransportPropagatesTransportError(t *testing.T) {
	transportErr := errors.New("connection refused")
	transport := &ReauthenticatingTransport{
		TokenSource: &replacingTokenSource{tokens: []string{"first-token", "second-token"}},
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, transportErr
		}),
	}

	response, err := transport.RoundTrip(authenticatedRequest(t, "first-token", ""))

	assert.Nil(t, response)
	assert.ErrorIs(t, err, transportErr)
}
