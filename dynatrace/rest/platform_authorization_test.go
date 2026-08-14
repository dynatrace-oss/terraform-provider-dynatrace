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
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/wif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// idToken builds a JWT that expires an hour from now. The signature segment is not a real signature:
// the provider never verifies these tokens, it only reads the expiry claim out of them.
func idToken(subject string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString(fmt.Appendf(nil, `{"sub":%q,"exp":%d}`, subject, time.Now().Add(time.Hour).Unix()))
	return header + "." + payload + ".not-a-real-signature"
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

// callPlatformAPI issues one platform request with the given credentials.
func callPlatformAPI(t *testing.T, serverURL string, platform PlatformCredentials) {
	t.Helper()

	platform.EnvironmentURL = serverURL
	credentials := &Credentials{URL: serverURL, Platform: platform}

	require.NoError(t, HybridClient(credentials).Get(t.Context(), "/platform/management/v1/environment").Finish())
}

func TestPlatformRequestSendsPlatformTokenAsBearerToken(t *testing.T) {
	serverURL, authorization := platformAPI(t)

	callPlatformAPI(t, serverURL, PlatformCredentials{PlatformToken: "dt0s16.ABCDEF"})

	assert.Equal(t, "Bearer dt0s16.ABCDEF", *authorization)
}

// The counterpart of the test above, differing only in which credential carries the token. Together
// they pin down that a federated token reaches the API in exactly the same wire format as a platform
// token, which is the premise the whole feature rests on.
func TestPlatformRequestSendsSuppliedOIDCTokenAsBearerToken(t *testing.T) {
	serverURL, authorization := platformAPI(t)
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")

	callPlatformAPI(t, serverURL, PlatformCredentials{WIF: wif.Config{StaticToken: token}})

	assert.Equal(t, "Bearer "+token, *authorization)
}

func TestPlatformRequestPrefersWorkloadIdentityFederationOverPlatformToken(t *testing.T) {
	serverURL, authorization := platformAPI(t)
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")

	callPlatformAPI(t, serverURL, PlatformCredentials{
		PlatformToken: "dt0s16.ABCDEF",
		WIF:           wif.Config{StaticToken: token},
	})

	assert.Equal(t, "Bearer "+token, *authorization)
}

func TestPlatformRequestPrefersWorkloadIdentityFederationOverOAuthCredentials(t *testing.T) {
	serverURL, authorization := platformAPI(t)
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")

	callPlatformAPI(t, serverURL, PlatformCredentials{
		ClientID:     "dt0s02.CLIENT",
		ClientSecret: "dt0s02.SECRET",
		TokenURL:     "https://sso.invalid/sso/oauth2/token",
		WIF:          wif.Config{StaticToken: token},
	})

	assert.Equal(t, "Bearer "+token, *authorization)
}

// The end to end path: the provider asks a token service for a token and puts that exact token on
// the platform request. Nothing is stubbed out - the token service is a second HTTP server, reached
// through the environment variables GitHub sets inside a job.
func TestPlatformRequestSendsMintedOIDCTokenAsBearerToken(t *testing.T) {
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	tokenService := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, `{"count":1,"value":%q}`, token)
	}))
	t.Cleanup(tokenService.Close)

	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", tokenService.URL)
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "request-token")

	serverURL, authorization := platformAPI(t)

	// The audience is unique to this test, because token sources are cached per audience for the
	// lifetime of the process and a shared one would outlive this test's token service.
	callPlatformAPI(t, serverURL, PlatformCredentials{
		WIF: wif.Config{Vendor: wif.VendorGitHub, Audience: t.Name()},
	})

	assert.Equal(t, "Bearer "+token, *authorization)
}

// A platform client is built for every single request. Without a token source that outlives those
// clients, each API call would go and fetch its own token, which is both slow and a good way to be
// rate limited by the token service.
func TestPlatformRequestReusesMintedOIDCTokenAcrossRequests(t *testing.T) {
	var mintRequests atomic.Int64
	token := idToken("repo:dynatrace-oss/terraform-provider-dynatrace:ref:refs/heads/main")
	tokenService := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mintRequests.Add(1)
		_, _ = fmt.Fprintf(writer, `{"count":1,"value":%q}`, token)
	}))
	t.Cleanup(tokenService.Close)

	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", tokenService.URL)
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "request-token")

	serverURL, _ := platformAPI(t)
	platform := PlatformCredentials{WIF: wif.Config{Vendor: wif.VendorGitHub, Audience: t.Name()}}

	callPlatformAPI(t, serverURL, platform)
	callPlatformAPI(t, serverURL, platform)

	assert.Equal(t, int64(1), mintRequests.Load())
}

// An apply can run for far longer than a federated token lives. The token is replaced ahead of its
// expiry, but a runner whose clock differs from the platform's is enough for one to be turned down
// anyway, and a run must not end over that. Both the token service and the platform are real HTTP
// servers here, so this covers the wiring between them rather than any single piece of it.
func TestPlatformRequestRetriesWithNewlyMintedOIDCTokenAfterRejection(t *testing.T) {
	tokens := []string{idToken("first-run"), idToken("second-run")}
	var mintRequests atomic.Int64
	tokenService := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		index := int(mintRequests.Add(1)) - 1
		if index >= len(tokens) {
			http.Error(writer, "asked for more tokens than this test prepared", http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprintf(writer, `{"count":1,"value":%q}`, tokens[index])
	}))
	t.Cleanup(tokenService.Close)

	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", tokenService.URL)
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "request-token")

	var accepted string
	platform := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authorization := request.Header.Get("Authorization")
		if authorization == "Bearer "+tokens[0] {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		accepted = authorization
		_, _ = writer.Write([]byte("{}"))
	}))
	t.Cleanup(platform.Close)

	callPlatformAPI(t, platform.URL, PlatformCredentials{WIF: wif.Config{Vendor: wif.VendorGitHub, Audience: t.Name()}})

	assert.Equal(t, "Bearer "+tokens[1], accepted)
}

func TestPlatformRequestReportsFailureToObtainOIDCToken(t *testing.T) {
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", "")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "")

	serverURL, _ := platformAPI(t)
	credentials := &Credentials{
		URL:      serverURL,
		Platform: PlatformCredentials{EnvironmentURL: serverURL, WIF: wif.Config{Vendor: wif.VendorGitHub, Audience: t.Name()}},
	}

	err := HybridClient(credentials).Get(t.Context(), "/platform/management/v1/environment").Finish()

	assert.EqualError(t, err, "Unable to get ACTIONS_ID_TOKEN_REQUEST_URL environment variable. `wif_vendor = \"github\"` only works inside a GitHub Actions job that is allowed to request an OIDC token. Add `permissions: { id-token: write }` to the workflow or to the job.")
}
