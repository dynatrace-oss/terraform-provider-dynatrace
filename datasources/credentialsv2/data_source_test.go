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

package credentialsv2

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const entriesPath = "/platform/credential-vault/v1/entries"

func platformClientSet(t *testing.T, serverURL string) *testing2.MockClientSet {
	platformClient, err := rest.CreatePlatformClient(t.Context(), serverURL, &rest.Credentials{
		Platform: rest.PlatformCredentials{EnvironmentURL: serverURL, PlatformToken: "token"},
	})
	require.NoError(t, err)
	return &testing2.MockClientSet{PlatformClientValue: platformClient}
}

func newServer(t *testing.T, listResponse string, receivedQuery *url.Values) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, entriesPath, r.URL.Path)
		if receivedQuery != nil {
			*receivedQuery = r.URL.Query()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(listResponse))
	}))
	t.Cleanup(server.Close)
	return server
}

func TestDataSourceRead_ReturnsAllCredentialsSharingTheName(t *testing.T) {
	server := newServer(t, `{"credentials": [
		{"id": "first-id", "name": "shared", "type": "TOKEN", "description": "the first", "scopes": ["SYNTHETIC", "APP_ENGINE"], "ownerAccessOnly": true},
		{"id": "second-id", "name": "shared", "type": "USERNAME_PASSWORD", "scopes": ["EXTENSION_AUTHENTICATION"], "ownerAccessOnly": false}]}`, nil)
	resourceData := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{"name": "shared"})

	diagnostics := dataSourceRead(t.Context(), resourceData, platformClientSet(t, server.URL))

	require.Empty(t, diagnostics)
	assert.Equal(t, "credentials-v2", resourceData.Id())
	assert.Equal(t, []any{
		map[string]any{
			"id":                "first-id",
			"name":              "shared",
			"type":              "TOKEN",
			"description":       "the first",
			"scopes":            []any{"SYNTHETIC", "APP_ENGINE"},
			"owner_access_only": true,
		},
		map[string]any{
			"id":                "second-id",
			"name":              "shared",
			"type":              "USERNAME_PASSWORD",
			"description":       "",
			"scopes":            []any{"EXTENSION_AUTHENTICATION"},
			"owner_access_only": false,
		},
	}, resourceData.Get("items"))
}

func TestDataSourceRead_SendsAllFiltersAsGiven(t *testing.T) {
	var receivedQuery url.Values
	server := newServer(t, `{"credentials": []}`, &receivedQuery)
	resourceData := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
		"name":  "my name",
		"type":  "TOKEN",
		"scope": "SYNTHETIC",
		"user":  "my-user",
	})

	diagnostics := dataSourceRead(t.Context(), resourceData, platformClientSet(t, server.URL))

	require.Empty(t, diagnostics)
	assert.Equal(t, url.Values{
		"name":      {"my name"},
		"type":      {"TOKEN"},
		"scope":     {"SYNTHETIC"},
		"user":      {"my-user"},
		"page-size": {"500"},
	}, receivedQuery)
}

func TestDataSourceRead_KeepsNamesDifferingOnlyInCase(t *testing.T) {
	server := newServer(t, `{"credentials": [
		{"id": "same-case-id", "name": "Shared"},
		{"id": "other-case-id", "name": "shared"}]}`, nil)
	resourceData := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{"name": "Shared"})

	diagnostics := dataSourceRead(t.Context(), resourceData, platformClientSet(t, server.URL))

	require.Empty(t, diagnostics)
	items := resourceData.Get("items").([]any)
	require.Len(t, items, 2)
	assert.Equal(t, "same-case-id", items[0].(map[string]any)["id"])
	assert.Equal(t, "other-case-id", items[1].(map[string]any)["id"])
}

func TestDataSourceRead_ReturnsEmptyItemsWhenNothingMatches(t *testing.T) {
	server := newServer(t, `{"credentials": []}`, nil)
	resourceData := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{"type": "TOKEN"})

	diagnostics := dataSourceRead(t.Context(), resourceData, platformClientSet(t, server.URL))

	require.Empty(t, diagnostics)
	assert.Equal(t, []any{}, resourceData.Get("items"))
}

func TestDataSourceRead_PropagatesListError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": {"code": 403, "message": "missing scope"}}`))
	}))
	t.Cleanup(server.Close)
	resourceData := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{"type": "TOKEN"})

	diagnostics := dataSourceRead(t.Context(), resourceData, platformClientSet(t, server.URL))

	assert.True(t, diagnostics.HasError())
}
