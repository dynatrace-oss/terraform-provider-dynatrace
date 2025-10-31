//go:build unit

/**
* @license
* Copyright 2024 Dynatrace LLC
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

package openpipeline_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("Returns an error if the client creation fails", func(t *testing.T) {
		_, err := openpipeline.EventsService(&rest.Credentials{}).List(t.Context())
		assert.ErrorIs(t, err, rest.NoPlatformCredentialsErr)
	})

	t.Run("Returns an error if the config doesn't exist anymore", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, err := w.Write([]byte("Not Found"))
			require.NoError(t, err)
		}))
		defer server.Close()
		_, err := openpipeline.EventsService(&rest.Credentials{
			OAuth: rest.OAuthCredentials{
				EnvironmentURL: server.URL,
				PlatformToken:  "token",
			},
		}).List(t.Context())
		assert.ErrorContains(t, err, "Not Found")
	})

	t.Run("Returns the list of configs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer server.Close()

		configs, err := openpipeline.EventsService(&rest.Credentials{
			OAuth: rest.OAuthCredentials{
				EnvironmentURL: server.URL,
				PlatformToken:  "token",
			},
		}).List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, configs, api2.Stubs{&api2.Stub{ID: "events", Name: "events"}})
	})
}
