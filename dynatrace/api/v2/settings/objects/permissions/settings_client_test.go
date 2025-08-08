package permissions_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

				writer.WriteHeader(http.StatusOK)
				_, err := writer.Write([]byte(`{"items":[{"schemaId":"schemaID1","ownerBasedAccessControl":true},{"schemaId":"schemaID2","ownerBasedAccessControl":false}]}`))
				require.NoError(t, err)
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

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
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

			schemaIDs, err := client.GetSchemaIDsWithOwnerBasedAccessControl(t.Context())

			assert.Error(t, err)
			assert.Nil(t, schemaIDs)
		})

		t.Run("returns error if connection fails", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {}))

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))
			server.Close()

			schemaIDs, err := client.GetSchemaIDsWithOwnerBasedAccessControl(t.Context())

			assert.Error(t, err)
			assert.Nil(t, schemaIDs)
		})
	})

	t.Run("ListObjectsIDsOfSchema", func(t *testing.T) {

		t.Run("retries without admin access if first request fails with forbidden", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Empty(t, request.URL.Query().Get("nextPageKey"))
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				if request.URL.Query().Get("adminAccess") == "true" {
					writer.WriteHeader(http.StatusForbidden)
					return
				}

				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(`{"items":[{"objectId":"object1"}]}`))
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.NoError(t, err)
			assert.Len(t, objectIDs, 1)
			assert.Contains(t, objectIDs, "object1")
		})

		t.Run("returns object IDs of schema", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(`{"items":[{"objectId":"object1"},{"objectId":"object2"}]}`))
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.NoError(t, err)
			assert.Len(t, objectIDs, 2)
			assert.Contains(t, objectIDs, "object1")
			assert.Contains(t, objectIDs, "object2")
		})

		t.Run("returns empty slice if no objects found", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(`{"items":[]}`))
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

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
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.Error(t, err)
			assert.Nil(t, objectIDs)
		})

		t.Run("returns error if connection fails", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {}))

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))
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
					writer.WriteHeader(http.StatusOK)
					writer.Write([]byte(`{"items":[{"objectId":"object2"}]}`))
					return
				}

				require.Equal(t, request.URL.Query().Get("schemaIds"), "schemaID1")
				require.Equal(t, request.URL.Query().Get("fields"), "objectId")
				require.Empty(t, request.URL.Query().Get("nextPageKey"))
				require.Equal(t, http.MethodGet, request.Method)
				require.Equal(t, "/platform/classic/environment-api/v2/settings/objects", request.URL.Path)

				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(`{"items":[{"objectId":"object1"}],"nextPageKey":"nextPage"}`))
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.NoError(t, err)
			assert.Len(t, objectIDs, 2)
			assert.Contains(t, objectIDs, "object1")
			assert.Contains(t, objectIDs, "object2")
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

				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("{\"items\":[{\"objectId\":\"object1\"}],\"nextPageKey\":\"nextPage\"}"))
			}))
			defer server.Close()

			url, _ := url.Parse(server.URL)
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

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
			client := permissions.NewSettingsClient(rest.NewClient(url, server.Client()))

			objectIDs, err := client.ListObjectsIDsOfSchema(t.Context(), "schemaID1")

			assert.Error(t, err)
			assert.Nil(t, objectIDs)
		})

	})

}
