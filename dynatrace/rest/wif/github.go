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
	"os"
)

const (
	// githubRequestURLVariable and githubRequestTokenVariable are injected into a GitHub Actions job
	// that holds the id-token: write permission. Together they are the minting credentials: they are
	// valid for the whole job and let the provider request a fresh OIDC token whenever it needs one.
	githubRequestURLVariable   = "ACTIONS_ID_TOKEN_REQUEST_URL"
	githubRequestTokenVariable = "ACTIONS_ID_TOKEN_REQUEST_TOKEN"

	// githubPermissionHint is appended to the errors above, because a missing permission is by far
	// the most likely reason for the variables to be absent inside a workflow.
	githubPermissionHint = "`wif_vendor = \"github\"` only works inside a GitHub Actions job that is allowed to request an OIDC token. Add `permissions: { id-token: write }` to the workflow or to the job."

	// maxTokenResponseSize bounds what is read from a token service. A JWT is a few kilobytes.
	maxTokenResponseSize = 1 << 20
)

// githubMinter requests OIDC tokens from the GitHub Actions token service.
type githubMinter struct {
	requestURL   string
	requestToken string
	audience     string
	httpClient   *http.Client
}

func newGitHubMinter(audience string, httpClient *http.Client) (minter, error) {
	requestURL := os.Getenv(githubRequestURLVariable)
	if len(requestURL) == 0 {
		return nil, ConfigError{fmt.Sprintf("Unable to get %s environment variable. %s", githubRequestURLVariable, githubPermissionHint)}
	}

	requestToken := os.Getenv(githubRequestTokenVariable)
	if len(requestToken) == 0 {
		return nil, ConfigError{fmt.Sprintf("Unable to get %s environment variable. %s", githubRequestTokenVariable, githubPermissionHint)}
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
		return "", fmt.Errorf("failed to get ID token: the value of %s is not a valid URL", githubRequestURLVariable)
	}
	// The URL GitHub provides already carries an api-version query parameter, so the audience is
	// merged into the existing query rather than appended to the raw string.
	query := endpoint.Query()
	query.Set("audience", minter.audience)
	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to get ID token: %w", err)
	}
	// The request token authenticates this call to the token service. It is a credential in its own
	// right and must never reach a log, an error message or a Terraform diagnostic.
	request.Header.Set("Authorization", "Bearer "+minter.requestToken)
	request.Header.Set("Accept", "application/json; api-version=2.0")

	response, err := minter.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to get ID token: %w", err)
	}
	defer drainAndClose(response.Body)

	if response.StatusCode != http.StatusOK {
		// The response body is deliberately left out of this message. On success that body is the
		// ID token itself, and echoing it on failure is one refactor away from leaking it.
		return "", fmt.Errorf("failed to get ID token: %s responded with HTTP %d", githubRequestURLVariable, response.StatusCode)
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
