//go:build unit

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

package rest_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreRest "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

// idToken builds a JWT that expires an hour from now.
func idToken(subject string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString(fmt.Appendf(nil, `{"sub":%q,"exp":%d}`, subject, time.Now().Add(time.Hour).Unix()))
	return header + "." + payload + ".not-a-real-signature"
}

// gitHubActionsTokenService starts a stand-in for the GitHub Actions token service answering with the
// given ID token, and injects its credentials the way GitHub does for a job holding the
// id-token: write permission. The returned counter holds the number of tokens it has handed out.
func gitHubActionsTokenService(t *testing.T, token string) *atomic.Int64 {
	t.Helper()

	var mintRequests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mintRequests.Add(1)
		_, _ = fmt.Fprintf(writer, `{"count":1,"value":%q}`, token)
	}))
	t.Cleanup(server.Close)

	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", server.URL)
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "request-token")

	return &mintRequests
}

// platformAPI starts a stand-in for the Dynatrace platform API that records the Authorization header
// of the request it receives.
func platformAPI(t *testing.T) (serverURL string, authorization *string) {
	t.Helper()

	var received string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		received = request.Header.Get("Authorization")
		_, _ = writer.Write([]byte("{}"))
	}))
	t.Cleanup(server.Close)

	return server.URL, &received
}

func callPlatformAPI(t *testing.T, serverURL string, platform rest.PlatformCredentials) error {
	t.Helper()

	platform.EnvironmentURL = serverURL
	credentials := &rest.Credentials{Platform: platform}

	return rest.HybridClient(createMockClientSet(t, credentials)).Get(t.Context(), "/platform/management/v1/environment").Finish()
}

func TestCreatePlatformClientUsesPlatformTokenAsBearerToken(t *testing.T) {
	serverURL, authorization := platformAPI(t)

	require.NoError(t, callPlatformAPI(t, serverURL, rest.PlatformCredentials{PlatformToken: "dt0s16.ABCDEF"}))

	assert.Equal(t, "Bearer dt0s16.ABCDEF", *authorization)
}

func TestCreatePlatformClientPrefersWIFOverPlatformToken(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	gitHubActionsTokenService(t, token)
	serverURL, authorization := platformAPI(t)

	require.NoError(t, callPlatformAPI(t, serverURL, rest.PlatformCredentials{
		PlatformToken:                      "dt0s16.ABCDEF",
		WorkloadIdentityFederationAudience: t.Name(),
	}))

	assert.Equal(t, "Bearer "+token, *authorization)
}

func TestCreatePlatformClientPrefersWIFOverOAuth(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	gitHubActionsTokenService(t, token)
	serverURL, authorization := platformAPI(t)

	require.NoError(t, callPlatformAPI(t, serverURL, rest.PlatformCredentials{
		ClientID:                           "dt0s02.CLIENT",
		ClientSecret:                       "dt0s02.SECRET",
		TokenURL:                           "https://sso.invalid/sso/oauth2/token",
		WorkloadIdentityFederationAudience: t.Name(),
	}))

	assert.Equal(t, "Bearer "+token, *authorization)
}

func TestCreatePlatformClientMintsOIDCTokenFromGitHubActions(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	gitHubActionsTokenService(t, token)
	serverURL, authorization := platformAPI(t)

	require.NoError(t, callPlatformAPI(t, serverURL, rest.PlatformCredentials{WorkloadIdentityFederationAudience: t.Name()}))

	assert.Equal(t, "Bearer "+token, *authorization)
}

func TestCreatePlatformClientReusesMintedTokenAcrossRequests(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	mintRequests := gitHubActionsTokenService(t, token)
	serverURL, _ := platformAPI(t)
	platform := rest.PlatformCredentials{EnvironmentURL: serverURL, WorkloadIdentityFederationAudience: t.Name()}
	clientSet := createMockClientSet(t, &rest.Credentials{Platform: platform})

	require.NoError(t, rest.HybridClient(clientSet).Get(t.Context(), "/platform/management/v1/environment").Finish())
	require.NoError(t, rest.HybridClient(clientSet).Get(t.Context(), "/platform/management/v1/environment").Finish())

	assert.Equal(t, int64(1), mintRequests.Load())
}

// The platform client and the classic-platform client are built by two separate calls to
// CreatePlatformClient, and deliberately do not share a token source between them - so each mints its
// own token rather than reusing the other's.
func TestCreatePlatformClientMintsIndependentlyPerCall(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	mintRequests := gitHubActionsTokenService(t, token)
	serverURL, _ := platformAPI(t)
	credentials := &rest.Credentials{Platform: rest.PlatformCredentials{EnvironmentURL: serverURL, WorkloadIdentityFederationAudience: t.Name()}}

	first, err := rest.CreatePlatformClient(t.Context(), serverURL, credentials)
	require.NoError(t, err)
	second, err := rest.CreatePlatformClient(t.Context(), serverURL, credentials)
	require.NoError(t, err)

	_, err = first.GET(t.Context(), "/platform/management/v1/environment", coreRest.RequestOptions{})
	require.NoError(t, err)
	_, err = second.GET(t.Context(), "/platform/management/v1/environment", coreRest.RequestOptions{})
	require.NoError(t, err)

	assert.Equal(t, int64(2), mintRequests.Load())
}

// federatedCredentials points the two hosts a platform request can be routed to at their own
// stand-in, so that a test can tell which of them served the request.
func federatedCredentials(t *testing.T, token string) (credentials *rest.Credentials, classicAuthorization, platformAuthorization *string) {
	t.Helper()

	gitHubActionsTokenService(t, token)
	classicURL, classicAuthorization := platformAPI(t)
	platformURL, platformAuthorization := platformAPI(t)

	return &rest.Credentials{
		ClassicEnvironmentURL: classicURL,
		Platform: rest.PlatformCredentials{
			EnvironmentURL:                     platformURL,
			WorkloadIdentityFederationAudience: t.Name(),
		},
	}, classicAuthorization, platformAuthorization
}

func TestClassicConfigRequestsGoToTheClassicHostWithTheFederatedToken(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	credentials, classicAuthorization, platformAuthorization := federatedCredentials(t, token)

	require.NoError(t, rest.HybridClient(createMockClientSet(t, credentials)).Get(t.Context(), "/api/config/v1/service/requestAttributes").Finish())

	assert.Equal(t, "Bearer "+token, *classicAuthorization)
	assert.Equal(t, "", *platformAuthorization)
}

// The counterpart of the test above, differing only in the requested path.
func TestNonClassicRequestsGoToThePlatformHostWithTheFederatedToken(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	credentials, classicAuthorization, platformAuthorization := federatedCredentials(t, token)

	require.NoError(t, rest.HybridClient(createMockClientSet(t, credentials)).Get(t.Context(), "/platform/management/v1/environment").Finish())

	assert.Equal(t, "", *classicAuthorization)
	assert.Equal(t, "Bearer "+token, *platformAuthorization)
}

// The audience asks for a federated token, but nothing in the environment can issue one.
func TestCreatePlatformClientReportsThatNoWorkloadIdentityProviderWasDetected(t *testing.T) {
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", "")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "")
	serverURL, _ := platformAPI(t)

	err := callPlatformAPI(t, serverURL, rest.PlatformCredentials{WorkloadIdentityFederationAudience: t.Name()})

	assert.EqualError(t, err, "no supported workload identity provider was detected: the only supported one is GitHub Actions, which injects the required `ACTIONS_ID_TOKEN_REQUEST_*` variables only into a job with `permissions: { id-token: write }`")
}

func TestCreatePlatformClientReportsNoCredentials(t *testing.T) {
	_, err := rest.CreatePlatformClient(t.Context(), "https://irrelevant.example", &rest.Credentials{})

	assert.ErrorIs(t, err, rest.ErrNoPlatformCredentials)
}
