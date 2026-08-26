/**
* @license
* Copyright 2026 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
)

const maxTokenResponseSize = 1 << 20

var (
	errMissingCredential    = errors.New("an OIDC token can only be requested from a GitHub Actions job with `permissions: { id-token: write }`")
	errResponseNotJSON      = errors.New("failed to get ID token: the response of the token service is not valid JSON")
	errResponseWithoutToken = errors.New("failed to get ID token: the response of the token service does not contain a token")
)

type githubMinter struct {
	requestURL   *url.URL
	requestToken string
	audience     string
	httpClient   *http.Client
}

func missingCredentialError(variableKey string) error {
	return fmt.Errorf("unable to get %s environment variable: %w", variableKey, errMissingCredential)
}

// The credentials GitHub injects stay valid for the whole job, which is what lets the provider mint
// a fresh token whenever it needs one.
func newGitHubMinter(audience string, httpClient *http.Client) (minter, error) {
	rawRequestURL := envutils.ActionsIDTokenRequestURL.Get()
	if len(rawRequestURL) == 0 {
		return nil, missingCredentialError(envutils.ActionsIDTokenRequestURL.Key)
	}

	requestURL, err := url.Parse(rawRequestURL)
	if err != nil {
		return nil, fmt.Errorf("%s does not hold a valid URL: %w", envutils.ActionsIDTokenRequestURL.Key, err)
	}

	requestToken := envutils.ActionsIDTokenRequestToken.Get()
	if len(requestToken) == 0 {
		return nil, missingCredentialError(envutils.ActionsIDTokenRequestToken.Key)
	}

	return &githubMinter{
		requestURL:   requestURL,
		requestToken: requestToken,
		audience:     audience,
		httpClient:   httpClient,
	}, nil
}

func (minter *githubMinter) mint(ctx context.Context) (string, error) {
	// Copied, so that setting the audience does not mutate the URL the minter keeps.
	endpoint := *minter.requestURL
	query := endpoint.Query()
	query.Set("audience", minter.audience)
	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to get ID token: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+minter.requestToken)
	request.Header.Set("Accept", "application/json; api-version=2.0")

	response, err := minter.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to get ID token: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		// The body is left out: on success it is the ID token itself.
		return "", fmt.Errorf("failed to get ID token: the token service responded with HTTP %d", response.StatusCode)
	}

	var body struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, maxTokenResponseSize)).Decode(&body); err != nil {
		return "", errResponseNotJSON
	}
	if len(body.Value) == 0 {
		return "", errResponseWithoutToken
	}

	return body.Value, nil
}
