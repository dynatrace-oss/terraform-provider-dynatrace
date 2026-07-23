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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// platformClientSet returns a ClientSet with a real platform client pointed at the test server.
func platformClientSet(t *testing.T, serverURL string) *testing2.MockClientSet {
	credentials := &rest.Credentials{
		Platform: rest.PlatformCredentials{EnvironmentURL: serverURL, PlatformToken: "token"},
	}
	platformClient, err := rest.CreatePlatformClient(t.Context(), serverURL, credentials)
	require.NoError(t, err)
	return &testing2.MockClientSet{PlatformClientValue: platformClient, CredentialsValue: credentials}
}

// TestServiceCreationFailsIfMissingClient tests that the service creation fails if the platform client is missing.
func TestServiceCreationFailsIfMissingClient(t *testing.T) {
	service, err := openpipeline.EventsService(&testing2.MockClientSet{PlatformClientErr: assert.AnError})
	require.Nil(t, service)
	require.ErrorIs(t, err, assert.AnError)
}

func TestList(t *testing.T) {
	t.Run("Returns an error if the config doesn't exist anymore", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			_, err := w.Write([]byte("Not Found"))
			require.NoError(t, err)
		}))
		defer server.Close()
		service, err := openpipeline.EventsService(platformClientSet(t, server.URL))
		require.NoError(t, err)
		_, err = service.List(t.Context())
		assert.ErrorContains(t, err, "Not Found")
	})

	t.Run("Returns the list of configs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer server.Close()

		service, err := openpipeline.EventsService(platformClientSet(t, server.URL))
		require.NoError(t, err)
		configs, err := service.List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, configs, api2.Stubs{&api2.Stub{ID: "events", Name: "events"}})
	})
}
