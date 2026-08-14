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
	"context"
	"net/http"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/wif"
	"golang.org/x/oauth2"
)

// bearerPrefix is the authentication scheme oauth2 puts in front of a bearer token.
const bearerPrefix = "Bearer "

// ReauthenticatingTransport sends a platform request once more, with a freshly obtained token, when
// the platform answers that the token it was sent with is not accepted.
//
// A federated token lives minutes while an apply can run for an hour, so [wif.TokenSource] replaces
// its token well ahead of the expiry it claims. That handles the ordinary case, but not every one:
// a runner whose clock differs from the platform's, or a token that is invalidated where it is
// trusted rather than where it is issued, produces a 401 for a token this provider still holds to be
// current. Obtaining a new token and sending the request again is then the difference between an
// apply that carries on and one that stops half way through.
//
// It sits below the oauth2 transport that puts the token on the request, which is what lets it read
// the token that was actually rejected off the request it is handed.
type ReauthenticatingTransport struct {
	Transport   http.RoundTripper
	TokenSource wif.TokenSource
}

// transport returns the configured RoundTripper, falling back to http.DefaultTransport when nil.
func (t *ReauthenticatingTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *ReauthenticatingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Buffer the request body up front, because sending the request consumes it and the retry needs
	// to send the same one.
	bodyBytes, err := bufferRequestBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := t.transport().RoundTrip(cloneRequestWithBody(req, bodyBytes))
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	ctx := req.Context()
	rejected := strings.TrimPrefix(req.Header.Get("Authorization"), bearerPrefix)

	token, replaceErr := t.TokenSource.Replace(rejected)
	if replaceErr != nil {
		// The original response is returned rather than this error. What the caller is entitled to
		// hear about is the platform turning its credential down; that a further attempt at getting
		// a new one also failed belongs in the log, where it explains why no retry was made.
		logging.Logger.Printf(ctx, "[ReauthenticatingTransport] Received HTTP 401 but obtaining a replacement token failed: %s", replaceErr)
		return resp, nil
	}
	if token.AccessToken == rejected {
		// Nothing was replaced - either the token cannot be replaced at all, or the same one was
		// handed out again. Sending the request again would produce the same 401.
		return resp, nil
	}

	logging.Logger.Printf(ctx, "[ReauthenticatingTransport] Received HTTP 401, retrying with a newly obtained token")
	if err := drainAndClose(resp.Body); err != nil {
		return nil, err
	}

	retry := cloneRequestWithBody(req, bodyBytes)
	token.SetAuthHeader(retry)
	return t.transport().RoundTrip(retry)
}

// NewContextWithReauthenticatingClient returns a context carrying the HTTP client that platform
// requests authenticated by the given token source are to be sent with.
//
// The oauth2 transport the platform client is built around picks this client up as its base and
// adds the token on top of it, so the retry ends up wrapped in the very transport that authenticates
// the request.
func NewContextWithReauthenticatingClient(ctx context.Context, tokenSource wif.TokenSource) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Transport: &ReauthenticatingTransport{TokenSource: tokenSource},
	})
}
