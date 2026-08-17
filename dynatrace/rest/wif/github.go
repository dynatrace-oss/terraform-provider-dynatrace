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

const (
	// A missing permission is by far the likeliest reason for the credentials to be absent.
	githubPermissionHint = "`wif_vendor = \"github\"` only works inside a GitHub Actions job that is allowed to request an OIDC token. Add `permissions: { id-token: write }` to the workflow or to the job."

	// A JWT is a few kilobytes.
	maxTokenResponseSize = 1 << 20
)

type githubMinter struct {
	requestURL   string
	requestToken string
	audience     string
	httpClient   *http.Client
}

// newGitHubMinter reads the credentials GitHub injects into the job. They stay valid for the whole
// job, which is what lets the provider mint a fresh token whenever it needs one.
func newGitHubMinter(audience string, httpClient *http.Client) (minter, error) {
	requestURL := envutils.ActionsIDTokenRequestURL.Get()
	if len(requestURL) == 0 {
		return nil, ConfigError{fmt.Sprintf("Unable to get %s environment variable. %s", envutils.ActionsIDTokenRequestURL.Key, githubPermissionHint)}
	}

	requestToken := envutils.ActionsIDTokenRequestToken.Get()
	if len(requestToken) == 0 {
		return nil, ConfigError{fmt.Sprintf("Unable to get %s environment variable. %s", envutils.ActionsIDTokenRequestToken.Key, githubPermissionHint)}
	}

	return &githubMinter{
		requestURL:   requestURL,
		requestToken: requestToken,
		audience:     audience,
		httpClient:   httpClient,
	}, nil
}

func (minter *githubMinter) mint(ctx context.Context) (string, error) {
	endpoint, err := url.Parse(minter.requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to get ID token: the value of %s is not a valid URL", envutils.ActionsIDTokenRequestURL.Key)
	}
	// Merged into the query rather than appended, because the URL already carries an api-version.
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
	defer drainAndClose(response.Body)

	if response.StatusCode != http.StatusOK {
		// The body is left out: on success it is the ID token itself, and echoing it on failure is
		// one refactor away from leaking it.
		return "", fmt.Errorf("failed to get ID token: %s responded with HTTP %d", envutils.ActionsIDTokenRequestURL.Key, response.StatusCode)
	}

	var body struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, maxTokenResponseSize)).Decode(&body); err != nil {
		return "", errors.New("failed to get ID token: the response of the token service is not valid JSON")
	}
	if len(body.Value) == 0 {
		return "", errors.New("failed to get ID token: the response of the token service does not contain a token")
	}

	return body.Value, nil
}
