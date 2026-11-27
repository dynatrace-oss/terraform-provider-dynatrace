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

package configuration_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/synthetic/configuration"
	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/synthetic/configuration/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const configurationPath = "/platform/synthetic/v1/configuration"

const currentConfiguration = `{"bmMonitorTimeout":60000,"bmStepTimeout":30000,"maintenanceWindowsExcludedFromAvailability":false}`

func ptr[T any](v T) *T { return &v }

type recordedRequest struct {
	Method string
	Path   string
	Body   string
}

// fakeServer answers GET with the given configuration and PUT with 204, recording every request.
// A non-zero status makes every request fail with that status instead.
type fakeServer struct {
	*httptest.Server
	mu       sync.Mutex
	requests []recordedRequest
}

func newFakeServer(t *testing.T, configuration string, status int) *fakeServer {
	t.Helper()
	fake := &fakeServer{}
	fake.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		fake.mu.Lock()
		fake.requests = append(fake.requests, recordedRequest{Method: r.Method, Path: r.URL.Path, Body: string(body)})
		fake.mu.Unlock()

		if status != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":{"code":1,"message":"failure"}}`))
			return
		}
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(configuration))
	}))
	t.Cleanup(fake.Close)
	return fake
}

func (me *fakeServer) recorded() []recordedRequest {
	me.mu.Lock()
	defer me.mu.Unlock()
	return append([]recordedRequest(nil), me.requests...)
}

func platformClientSet(t *testing.T, serverURL string) *testing2.MockClientSet {
	t.Helper()
	platformClient, err := rest.CreatePlatformClient(t.Context(), serverURL, &rest.Credentials{
		Platform: rest.PlatformCredentials{EnvironmentURL: serverURL, PlatformToken: "token"},
	})
	require.NoError(t, err)
	return &testing2.MockClientSet{PlatformClientValue: platformClient}
}

func TestServiceCreationFailsIfClientUnavailable(t *testing.T) {
	service, err := configuration.Service(&testing2.MockClientSet{PlatformClientErr: assert.AnError})

	require.Nil(t, service)
	require.ErrorIs(t, err, assert.AnError)
	assert.Contains(t, err.Error(), "synthetic:environment-configuration:read")
	assert.Contains(t, err.Error(), "synthetic:environment-configuration:write")
}

func TestSchemaID(t *testing.T) {
	server := newFakeServer(t, currentConfiguration, 0)
	service, err := configuration.Service(platformClientSet(t, server.URL))
	require.NoError(t, err)

	assert.Equal(t, "platform:synthetic:configuration", service.SchemaID())
}

func TestList(t *testing.T) {
	server := newFakeServer(t, currentConfiguration, 0)
	service, err := configuration.Service(platformClientSet(t, server.URL))
	require.NoError(t, err)

	stubs, err := service.List(t.Context())

	require.NoError(t, err)
	require.Len(t, stubs, 1)
	assert.Equal(t, configuration.ID, stubs[0].ID)
	assert.Empty(t, server.recorded())
}

func TestGet(t *testing.T) {
	t.Run("Maps all fields and ignores the ID", func(t *testing.T) {
		server := newFakeServer(t, `{"bmMonitorTimeout":120000,"bmStepTimeout":45000,"maintenanceWindowsExcludedFromAvailability":true}`, 0)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		result := new(settings.Settings)
		require.NoError(t, service.Get(t.Context(), "any-id-works", result))

		assert.Equal(t, &settings.Settings{
			BmMonitorTimeout: ptr(int64(120000)),
			BmStepTimeout:    ptr(int64(45000)),
			MaintenanceWindowsExcludedFromAvailability: ptr(true),
		}, result)
		requests := server.recorded()
		require.Len(t, requests, 1)
		assert.Equal(t, http.MethodGet, requests[0].Method)
		assert.Equal(t, configurationPath, requests[0].Path)
	})

	t.Run("Propagates server errors", func(t *testing.T) {
		server := newFakeServer(t, "", http.StatusForbidden)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		require.Error(t, service.Get(t.Context(), configuration.ID, new(settings.Settings)))
	})
}

func TestCreate(t *testing.T) {
	t.Run("Reads then writes all fields and returns the fixed ID", func(t *testing.T) {
		server := newFakeServer(t, currentConfiguration, 0)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		stub, err := service.Create(t.Context(), &settings.Settings{
			BmMonitorTimeout: ptr(int64(120000)),
			BmStepTimeout:    ptr(int64(45000)),
			MaintenanceWindowsExcludedFromAvailability: ptr(true),
		})

		require.NoError(t, err)
		assert.Equal(t, configuration.ID, stub.ID)
		requests := server.recorded()
		require.Len(t, requests, 2)
		assert.Equal(t, http.MethodGet, requests[0].Method)
		assert.Equal(t, configurationPath, requests[0].Path)
		assert.Equal(t, http.MethodPut, requests[1].Method)
		assert.Equal(t, configurationPath, requests[1].Path)
		assert.JSONEq(t, `{"bmMonitorTimeout":120000,"bmStepTimeout":45000,"maintenanceWindowsExcludedFromAvailability":true}`, requests[1].Body)
	})

	t.Run("Keeps unset fields from the tenant", func(t *testing.T) {
		server := newFakeServer(t, currentConfiguration, 0)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		_, err = service.Create(t.Context(), &settings.Settings{MaintenanceWindowsExcludedFromAvailability: ptr(true)})

		require.NoError(t, err)
		requests := server.recorded()
		require.Len(t, requests, 2)
		assert.JSONEq(t, `{"bmMonitorTimeout":60000,"bmStepTimeout":30000,"maintenanceWindowsExcludedFromAvailability":true}`, requests[1].Body)
	})

	t.Run("Propagates a failing read without writing", func(t *testing.T) {
		server := newFakeServer(t, "", http.StatusForbidden)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		stub, err := service.Create(t.Context(), &settings.Settings{MaintenanceWindowsExcludedFromAvailability: ptr(true)})

		require.Error(t, err)
		assert.Nil(t, stub)
		assert.Len(t, server.recorded(), 1)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Writes all fields", func(t *testing.T) {
		server := newFakeServer(t, currentConfiguration, 0)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		require.NoError(t, service.Update(t.Context(), configuration.ID, &settings.Settings{
			BmMonitorTimeout: ptr(int64(900000)),
			BmStepTimeout:    ptr(int64(10000)),
			MaintenanceWindowsExcludedFromAvailability: ptr(true),
		}))

		requests := server.recorded()
		require.Len(t, requests, 2)
		assert.Equal(t, http.MethodGet, requests[0].Method)
		assert.Equal(t, http.MethodPut, requests[1].Method)
		assert.JSONEq(t, `{"bmMonitorTimeout":900000,"bmStepTimeout":10000,"maintenanceWindowsExcludedFromAvailability":true}`, requests[1].Body)
	})

	t.Run("Keeps unset fields from the tenant", func(t *testing.T) {
		server := newFakeServer(t, `{"bmMonitorTimeout":60000,"bmStepTimeout":30000,"maintenanceWindowsExcludedFromAvailability":true}`, 0)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		require.NoError(t, service.Update(t.Context(), configuration.ID, &settings.Settings{BmStepTimeout: ptr(int64(90000))}))

		requests := server.recorded()
		require.Len(t, requests, 2)
		assert.JSONEq(t, `{"bmMonitorTimeout":60000,"bmStepTimeout":90000,"maintenanceWindowsExcludedFromAvailability":true}`, requests[1].Body)
	})

	t.Run("Sends an explicit false", func(t *testing.T) {
		server := newFakeServer(t, `{"bmMonitorTimeout":60000,"bmStepTimeout":30000,"maintenanceWindowsExcludedFromAvailability":true}`, 0)
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		require.NoError(t, service.Update(t.Context(), configuration.ID, &settings.Settings{MaintenanceWindowsExcludedFromAvailability: ptr(false)}))

		requests := server.recorded()
		require.Len(t, requests, 2)
		assert.JSONEq(t, `{"bmMonitorTimeout":60000,"bmStepTimeout":30000,"maintenanceWindowsExcludedFromAvailability":false}`, requests[1].Body)
	})

	t.Run("Propagates a failing write", func(t *testing.T) {
		// GET succeeds, PUT is rejected with 400.
		var puts int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				puts++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error":{"code":400,"message":"bad"}}`))
				return
			}
			_, _ = w.Write([]byte(currentConfiguration))
		}))
		defer server.Close()
		service, err := configuration.Service(platformClientSet(t, server.URL))
		require.NoError(t, err)

		require.Error(t, service.Update(t.Context(), configuration.ID, &settings.Settings{MaintenanceWindowsExcludedFromAvailability: ptr(true)}))
		assert.Equal(t, 1, puts)
	})
}

func TestDelete(t *testing.T) {
	server := newFakeServer(t, currentConfiguration, 0)
	service, err := configuration.Service(platformClientSet(t, server.URL))
	require.NoError(t, err)

	require.NoError(t, service.Delete(t.Context(), configuration.ID))

	assert.Empty(t, server.recorded())
}
