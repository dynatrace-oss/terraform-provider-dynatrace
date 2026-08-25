//go:build unit

/**
* @license
* Copyright 2025 Dynatrace LLC
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
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockToken = "########"

func TestHybridClient(t *testing.T) {
	baseURL, tokenURL := newHybridTestServer(t)

	oauth := rest.PlatformCredentials{EnvironmentURL: baseURL, ClientID: mockToken, ClientSecret: mockToken, TokenURL: tokenURL}
	oauthAndPlatformToken := rest.PlatformCredentials{EnvironmentURL: baseURL, ClientID: mockToken, ClientSecret: mockToken, TokenURL: tokenURL, PlatformToken: mockToken}
	platformToken := rest.PlatformCredentials{EnvironmentURL: baseURL, PlatformToken: mockToken}

	tests := []struct {
		name           string
		creds          *rest.Credentials
		oauthPreferred bool
		wantErr        error
		wantClassic    bool
		wantPlatform   bool
	}{
		{name: "unconfigured", creds: &rest.Credentials{}, wantErr: rest.NoAPITokenError},
		{name: "unconfigured, oauth preferred", creds: &rest.Credentials{}, oauthPreferred: true, wantErr: rest.NoOAuthCredentialsError},

		{name: "api token", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken}, wantClassic: true},
		{name: "api token, oauth preferred", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken}, oauthPreferred: true, wantClassic: true},
		{name: "oauth", creds: &rest.Credentials{Platform: oauth}, wantPlatform: true},
		{name: "oauth, oauth preferred", creds: &rest.Credentials{Platform: oauth}, oauthPreferred: true, wantPlatform: true},

		{name: "oauth and platform token", creds: &rest.Credentials{Platform: oauthAndPlatformToken}, wantPlatform: true},
		{name: "oauth and platform token, oauth preferred", creds: &rest.Credentials{Platform: oauthAndPlatformToken}, oauthPreferred: true, wantPlatform: true},

		{name: "api token and oauth", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken, Platform: oauth}, wantClassic: true},
		{name: "api token and oauth, oauth preferred", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken, Platform: oauth}, oauthPreferred: true, wantPlatform: true},

		{name: "api token and platform token", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken, Platform: platformToken}, wantClassic: true},
		{name: "api token and platform token, oauth preferred", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken, Platform: platformToken}, oauthPreferred: true, wantPlatform: true},

		{name: "api token, oauth and platform token", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken, Platform: oauthAndPlatformToken}, wantClassic: true},
		{name: "api token, oauth and platform token, oauth preferred", creds: &rest.Credentials{ClassicEnvironmentURL: baseURL, Token: mockToken, Platform: oauthAndPlatformToken}, oauthPreferred: true, wantPlatform: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			if tt.oauthPreferred {
				ctx = context.WithValue(ctx, envutils.DynatraceHTTPOAuthPreference.Key, true)
			}

			var v struct {
				Classic  bool `json:"classic"`
				Platform bool `json:"platform"`
			}
			err := rest.HybridClient(tt.creds).Get(ctx, "").Finish(&v)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantClassic, v.Classic)
			assert.Equal(t, tt.wantPlatform, v.Platform)
		})
	}
}

func newHybridTestServer(t *testing.T) (baseURL, tokenURL string) {
	tokenPath := "/sso/oauth2/token"

	writeJSON := func(w http.ResponseWriter, body string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, `{"access_token": "tok","token_type":"Bearer","expires_in":3600}`)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Authorization"), "Api-Token") {
			writeJSON(w, `{"classic": true}`)
			return
		}
		writeJSON(w, `{"platform": true}`)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv.URL, srv.URL + tokenPath
}

func TestApiTokenClient(t *testing.T) {
	const (
		endpoint          = "/api/v2/settings/objects"
		activeGatePostfix = "/e/my-env-id"
	)

	t.Run("Correctly transforms an Active Gate URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedURL, err := url.JoinPath(activeGatePostfix, endpoint)
			require.NoError(t, err)
			require.Equal(t, expectedURL, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		activeGateURL, err := url.JoinPath(server.URL, activeGatePostfix)
		require.NoError(t, err)

		cred := rest.Credentials{ClassicEnvironmentURL: activeGateURL, Token: mockToken}
		client := rest.HybridClient(&cred)

		req := client.Get(t.Context(), endpoint)
		err = req.Finish()
		assert.NoError(t, err)
	})

	t.Run("Errors on empty env URL", func(t *testing.T) {
		cred := rest.Credentials{ClassicEnvironmentURL: "", Token: mockToken}
		client := rest.HybridClient(&cred)

		req := client.Get(t.Context(), endpoint)
		err := req.Finish()
		assert.ErrorIs(t, err, rest.NoClassicURLDefinedErr)
	})

	t.Run("Errors on invalid path", func(t *testing.T) {
		cred := rest.Credentials{ClassicEnvironmentURL: "my-url", Token: mockToken}
		client := rest.HybridClient(&cred)

		req := client.Get(t.Context(), ":/invalid-url")
		err := req.Finish()
		expectedErr := &url.Error{}
		assert.ErrorAs(t, err, &expectedErr)
	})
}
