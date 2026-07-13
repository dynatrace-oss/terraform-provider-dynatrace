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
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// platformClientSet returns a ClientSet with a real platform client pointed at the test server.
func platformClientSet(t *testing.T, serverURL string) *testing2.MockClientSet {
	platformClient, err := rest.CreatePlatformClient(t.Context(), serverURL, &rest.Credentials{
		Platform: rest.PlatformCredentials{EnvironmentURL: serverURL, PlatformToken: "token"},
	})
	require.NoError(t, err)
	return &testing2.MockClientSet{PlatformClientValue: platformClient}
}

func TestList(t *testing.T) {
	t.Run("Returns an error if the client creation fails", func(t *testing.T) {
		_, err := openpipeline.EventsService(&config.ProviderConfiguration{}).List(t.Context())
		assert.ErrorIs(t, err, rest.NoPlatformCredentialsErr)
	})

	t.Run("Returns an error if the config doesn't exist anymore", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, err := w.Write([]byte("Not Found"))
			require.NoError(t, err)
		}))
		defer server.Close()
		_, err := openpipeline.EventsService(platformClientSet(t, server.URL)).List(t.Context())
		assert.ErrorContains(t, err, "Not Found")
	})

	t.Run("Returns the list of configs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer server.Close()

		configs, err := openpipeline.EventsService(platformClientSet(t, server.URL)).List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, configs, api2.Stubs{&api2.Stub{ID: "events", Name: "events"}})
	})
}
