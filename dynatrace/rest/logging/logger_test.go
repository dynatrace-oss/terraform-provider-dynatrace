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

package logging_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/github/connector/connection"
	set "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/github/connector/connection/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/boundaries"
	setboundaries "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/boundaries/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger(t *testing.T) *strings.Builder {
	t.Helper()
	var builder strings.Builder
	logging.Logger = &logging.RESTLogger{
		Log: log.New(&builder, "", log.LstdFlags),
	}
	return &builder
}

func assertLoggedOnce(t *testing.T, output string) {
	t.Helper()
	assert.Equal(t, 1, strings.Count(output, "[REQUEST]"))
	assert.Equal(t, 1, strings.Count(output, "[RESPONSE]"))
}

func TestProviderLogging(t *testing.T) {
	t.Setenv("DYNATRACE_LOG_HTTP", "true")
	t.Setenv("DYNATRACE_HTTP_RESPONSE", "true")

	connectionTests := []struct {
		name        string
		credentials func(url string) *rest.Credentials
	}{
		{
			name: "log HTTP requests and responses only once for platform token",
			credentials: func(url string) *rest.Credentials {
				return &rest.Credentials{
					URL: url,
					OAuth: rest.OAuthCredentials{
						PlatformToken: "my-token",
					},
				}
			},
		},
		{
			name: "log HTTP requests and responses only once for classic",
			credentials: func(url string) *rest.Credentials {
				return &rest.Credentials{
					URL:   url,
					Token: "my-token",
				}
			},
		},
	}
	for _, tc := range connectionTests {
		t.Run(tc.name, func(t *testing.T) {
			builder := newTestLogger(t)
			mux := http.ServeMux{}
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"objectId": "test-id", "value": {}}`))
				require.NoError(t, err)
			})
			httpServer := httptest.NewServer(&mux)
			defer httpServer.Close()

			service := connection.Service(tc.credentials(httpServer.URL))
			err := service.Get(t.Context(), "test-id", new(set.Settings))
			require.NoError(t, err)
			assertLoggedOnce(t, builder.String())
		})
	}

	t.Run("log HTTP requests and responses only once for IAM", func(t *testing.T) {
		builder := newTestLogger(t)
		mux := http.ServeMux{}
		mux.HandleFunc("/sso", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"access_token": "tok","token_type":"Bearer","expires_in":3600}`))
			require.NoError(t, err)
		})
		mux.HandleFunc("GET /iam/v1/repo/account/account-id/boundaries/test-id", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			data, err := json.Marshal(setboundaries.PolicyBoundary{})
			require.NoError(t, err)
			_, err = w.Write(data)
			require.NoError(t, err)
		})
		httpServer := httptest.NewServer(&mux)
		defer httpServer.Close()

		service := boundaries.Service(&rest.Credentials{
			URL: httpServer.URL,
			IAM: struct {
				ClientID     string
				AccountID    string
				ClientSecret string
				TokenURL     string
				EndpointURL  string
			}{
				ClientID:     "id",
				AccountID:    "account-id",
				ClientSecret: "secret",
				TokenURL:     httpServer.URL + "/sso",
				EndpointURL:  httpServer.URL,
			},
		})
		err := service.Get(t.Context(), "test-id", new(setboundaries.PolicyBoundary))
		require.NoError(t, err)
		assertLoggedOnce(t, builder.String())
	})
}
