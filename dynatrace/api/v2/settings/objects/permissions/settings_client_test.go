//go:build unit

/*
 * @license
 * Copyright 2025 Dynatrace LLC
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

package permissions_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettingsClient(t *testing.T) {
	t.Run("GetSchemaIDsWithOwnerBasedAccessControl", func(t *testing.T) {
		t.Run("Returns schema IDs with owner-based access control", func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, request.URL.Query().Get("fields"), "schemaId,ownerBasedAccessControl")
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/schemas", request.URL.Path)

				_, err := writer.Write([]byte(`{"items":[{"schemaId":"schemaID1","ownerBasedAccessControl":true},{"schemaId":"schemaID2","ownerBasedAccessControl":false}]}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			schemaIDs, err := client.GetSchemaIDsWithOwnerBasedAccessControl(t.Context())

			assert.NoError(t, err)
			assert.Len(t, schemaIDs, 1)
			assert.Contains(t, schemaIDs, "schemaID1")
		})

		t.Run("returns error if response is not 200", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusInternalServerError)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			schemaIDs, err := client.GetSchemaIDsWithOwnerBasedAccessControl(t.Context())

			assert.Error(t, err)
			assert.Nil(t, schemaIDs)
		})

		t.Run("returns error if response is not JSON", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, err := writer.Write([]byte(`{`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			schemaIDs, err := client.GetSchemaIDsWithOwnerBasedAccessControl(t.Context())

			syntaxErr := &json.SyntaxError{}
			assert.ErrorAs(t, err, &syntaxErr)
			assert.Nil(t, schemaIDs)
		})

		t.Run("returns error if connection fails", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {}))

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))
			server.Close()

			schemaIDs, err := client.GetSchemaIDsWithOwnerBasedAccessControl(t.Context())

			assert.Error(t, err)
			assert.Nil(t, schemaIDs)
		})
	})

	t.Run("ListObjectsIDsOfSchema", func(t *testing.T) {

		t.Run("retries without admin access if first request fails with forbidden", func(t *testing.T) {

			var triedWithAdminAccess atomic.Bool

			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Empty(t, request.URL.Query().Get("nextPageKey"))
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				if request.URL.Query().Get("adminAccess") == "true" {
					triedWithAdminAccess.Store(true)
					writer.WriteHeader(http.StatusForbidden)
					return
				}

				_, err := writer.Write([]byte(`{"items":[{"objectId":"object1"}]}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.NoError(t, err)
			assert.Len(t, objectIDs, 1)
			assert.Contains(t, objectIDs, "object1")

			assert.True(t, triedWithAdminAccess.Load(), "Missing request with adminAccess=true")
		})

		t.Run("returns object IDs of schema", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				_, err := writer.Write([]byte(`{"items":[{"objectId":"object1"},{"objectId":"object2"}]}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.NoError(t, err)
			assert.Len(t, objectIDs, 2)
			assert.Contains(t, objectIDs, "object1")
			assert.Contains(t, objectIDs, "object2")
		})

		t.Run("returns empty slice if no objects found", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, err := writer.Write([]byte(`{"items":[]}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.NoError(t, err)
			assert.Empty(t, objectIDs)
		})

		t.Run("returns error if response is not 200", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusInternalServerError)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.Error(t, err)
			assert.Nil(t, objectIDs)
		})

		t.Run("returns errors if first page is not JSON", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, err := writer.Write([]byte(`{items:[]}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			syntaxErr := &json.SyntaxError{}
			assert.ErrorAs(t, err, &syntaxErr)
			assert.Empty(t, objectIDs)
		})

		t.Run("returns error if connection fails", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {}))

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))
			server.Close()

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.Error(t, err)
			assert.Nil(t, objectIDs)
		})

		t.Run("returns object IDs with pagination", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.Query().Get("nextPageKey") != "" {
					require.Equal(t, request.URL.Query().Get("nextPageKey"), "nextPage")
					require.Len(t, request.URL.Query(), 1)

					_, err := writer.Write([]byte(`{"items":[{"objectId":"object2"}]}`))
					require.NoError(t, err)
					return
				}

				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Empty(t, request.URL.Query().Get("nextPageKey"))
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				_, err := writer.Write([]byte(`{"items":[{"objectId":"object1"}],"nextPageKey":"nextPage"}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.NoError(t, err)
			assert.Len(t, objectIDs, 2)
			assert.Contains(t, objectIDs, "object1")
			assert.Contains(t, objectIDs, "object2")
		})

		t.Run("Returns error if second page is not JSON", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.Query().Get("nextPageKey") != "" {
					require.Equal(t, request.URL.Query().Get("nextPageKey"), "nextPage")
					require.Len(t, request.URL.Query(), 1)

					_, err := writer.Write([]byte(`{items:[{objectId:object2}]}`))
					require.NoError(t, err)
					return
				}

				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Empty(t, request.URL.Query().Get("nextPageKey"))
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				_, err := writer.Write([]byte(`{"items":[{"objectId":"object1"}],"nextPageKey":"nextPage"}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			syntaxErr := &json.SyntaxError{}
			assert.ErrorAs(t, err, &syntaxErr)
			assert.Empty(t, objectIDs, 2)
		})

		t.Run("returns error if pagination fails", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.Query().Get("nextPageKey") != "" {
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}

				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Empty(t, request.URL.Query().Get("nextPageKey"))
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				_, err := writer.Write([]byte("{\"items\":[{\"objectId\":\"object1\"}],\"nextPageKey\":\"nextPage\"}"))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.Error(t, err)
			assert.Nil(t, objectIDs)
		})

		t.Run("returns error if first pagination response is not 200", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Empty(t, request.URL.Query().Get("nextPageKey"))
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				writer.WriteHeader(http.StatusInternalServerError)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewPlatformSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.Error(t, err)
			assert.Nil(t, objectIDs)
		})

	})

}
