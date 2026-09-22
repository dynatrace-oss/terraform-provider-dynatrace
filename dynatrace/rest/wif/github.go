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
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
)

const maxTokenResponseSize = 1 << 20

var (
	errResponseNotJSON      = errors.New("failed to get ID token: the response of the token service is not valid JSON")
	errResponseWithoutToken = errors.New("failed to get ID token: the response of the token service does not contain a token")
)

type gitHubConfig struct {
	audience          string
	tokenRequestURL   string
	tokenRequestToken string
}

// GitHub injects both credentials of its token service into every job holding the id-token: write
// permission, and they stay valid for the whole job - which is what lets the provider mint a fresh
// token whenever it needs one. Their presence is what identifies GitHub as the vendor.
func inferGitHubConfig(audience string) (gitHubConfig, bool) {
	tokenRequestURL := strings.TrimSpace(envutils.ActionsIDTokenRequestURL.Get())
	tokenRequestToken := strings.TrimSpace(envutils.ActionsIDTokenRequestToken.Get())
	if len(tokenRequestURL) == 0 || len(tokenRequestToken) == 0 {
		return gitHubConfig{}, false
	}

	return gitHubConfig{
		audience:          audience,
		tokenRequestURL:   tokenRequestURL,
		tokenRequestToken: tokenRequestToken,
	}, true
}

// The URL is parsed once, here, so that a malformed one is reported while the provider is being
// configured rather than on every mint.
func (config gitHubConfig) createMinter(httpClient *http.Client) (minter, error) {
	requestURL, err := url.Parse(config.tokenRequestURL)
	if err != nil {
		return nil, fmt.Errorf("the GitHub Actions token request URL is not valid: %w", err)
	}

	return &gitHubMinter{
		requestURL:   requestURL,
		requestToken: config.tokenRequestToken,
		audience:     config.audience,
		httpClient:   httpClient,
	}, nil
}

type gitHubMinter struct {
	requestURL   *url.URL
	requestToken string
	audience     string
	httpClient   *http.Client
}

func (minter *gitHubMinter) mint(ctx context.Context) (string, error) {
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
