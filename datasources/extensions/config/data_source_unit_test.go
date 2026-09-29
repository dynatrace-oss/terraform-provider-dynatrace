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

package extension_config_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"

	extensionconfig "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/extensions/config"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
)

// listCall records the arguments the data source passed to the client.
type listCall struct {
	extensionName string
	filter        string
}

func newMockClient(call *listCall, response coreapi.PagedListResponse, err error) *testing2.MockExtensionClient {
	return &testing2.MockExtensionClient{
		ListMonitoringConfigurationsFn: func(_ context.Context, extensionName string, filter string) (coreapi.PagedListResponse, error) {
			call.extensionName = extensionName
			call.filter = filter
			return response, err
		},
	}
}

// pagedResponse wraps the given JSON objects into a single page.
func pagedResponse(objects ...string) coreapi.PagedListResponse {
	page := coreapi.ListResponse{Response: coreapi.Response{StatusCode: 200}}
	for _, object := range objects {
		page.Objects = append(page.Objects, []byte(object))
	}
	return coreapi.PagedListResponse{page}
}

func newResourceData(t *testing.T, raw map[string]any) *schema.ResourceData {
	t.Helper()
	if raw == nil {
		raw = map[string]any{}
	}
	if _, found := raw["name"]; !found {
		raw["name"] = "com.dynatrace.extension.test"
	}
	return schema.TestResourceDataRaw(t, extensionconfig.DataSource().Schema, raw)
}

func TestDataSource_ItemSchema(t *testing.T) {
	items := extensionconfig.DataSource().Schema["items"]
	require.NotNil(t, items)
	assert.True(t, items.Computed)

	elem, ok := items.Elem.(*schema.Resource)
	require.True(t, ok, "expected the items Elem to be a *schema.Resource")

	assert.NotContains(t, elem.Schema, "name", "expected the foreign-key `name` to be removed from the item schema")
	for name, item := range elem.Schema {
		assert.True(t, item.Computed, "%s: expected Computed to be true", name)
		assert.False(t, item.Required, "%s: expected Required to be false", name)
		assert.Nil(t, item.ValidateFunc, "%s: expected ValidateFunc to be nil", name)
	}
	assert.Contains(t, elem.Schema, "id")
	assert.Contains(t, elem.Schema, "scope")
	assert.Contains(t, elem.Schema, "value")
}

func TestDataSourceReadWithClient_MapsConfigurations(t *testing.T) {
	d := newResourceData(t, nil)
	var call listCall
	client := newMockClient(&call, pagedResponse(
		`{"objectId":"vu9U3hXa3q0AAAAB","scope":"HOST-1234","value":{"version":"1.2.3","description":"first"}}`,
		`{"objectId":"vu9U3hXa3q0AAAAC","scope":"HOST-5678","value":{"version":"1.2.3","description":"second"}}`,
	), nil)

	diags := extensionconfig.DataSourceReadWithClient(t.Context(), d, client)

	require.False(t, diags.HasError(), "expected no error, got: %v", diags)
	assert.Equal(t, "com.dynatrace.extension.test", d.Id())
	assert.Equal(t, "com.dynatrace.extension.test", call.extensionName)

	items := d.Get("items").([]any)
	require.Len(t, items, 2)

	first := items[0].(map[string]any)
	assert.Equal(t, "vu9U3hXa3q0AAAAB", first["id"])
	assert.Equal(t, "HOST-1234", first["scope"])
	assert.JSONEq(t, `{"version":"1.2.3","description":"first"}`, first["value"].(string))
	assert.NotContains(t, first, "name", "expected the extension name not to be part of an item")

	second := items[1].(map[string]any)
	assert.Equal(t, "vu9U3hXa3q0AAAAC", second["id"])
	assert.Equal(t, "HOST-5678", second["scope"])
}

func TestDataSourceReadWithClient_NoConfigurations(t *testing.T) {
	d := newResourceData(t, nil)
	client := newMockClient(&listCall{}, coreapi.PagedListResponse{}, nil)

	diags := extensionconfig.DataSourceReadWithClient(t.Context(), d, client)

	require.False(t, diags.HasError(), "expected no error, got: %v", diags)
	assert.Empty(t, d.Get("items").([]any))
	assert.Equal(t, "com.dynatrace.extension.test", d.Id())
}

func TestDataSourceReadWithClient_ClientError(t *testing.T) {
	d := newResourceData(t, nil)
	client := newMockClient(&listCall{}, nil, assert.AnError)

	diags := extensionconfig.DataSourceReadWithClient(t.Context(), d, client)

	require.True(t, diags.HasError(), "expected an error diagnostic, got none")
	assert.ElementsMatch(t, diags, diag.FromErr(assert.AnError))
	assert.Empty(t, d.Id(), "expected the ID not to be set after a failed read")
}

func TestDataSourceReadWithClient_InvalidJSON(t *testing.T) {
	d := newResourceData(t, nil)
	client := newMockClient(&listCall{}, pagedResponse(`not-valid-json`), nil)

	diags := extensionconfig.DataSourceReadWithClient(t.Context(), d, client)

	assert.True(t, diags.HasError(), "expected an error diagnostic from invalid JSON, got none")
}

func TestDataSourceReadWithClient_Filter(t *testing.T) {
	tests := []struct {
		name     string
		raw      map[string]any
		expected string
	}{
		{
			name:     "version",
			raw:      map[string]any{"version": "1.2.3"},
			expected: "version='1.2.3'",
		},
		{
			name:     "description",
			raw:      map[string]any{"description": "my description"},
			expected: "description='my description'",
		},
		{
			name:     "activation context",
			raw:      map[string]any{"activation_context": "REMOTE"},
			expected: "activationContext='REMOTE'",
		},
		{
			name: "all arguments are joined in a fixed order",
			raw: map[string]any{
				"version":            "1.2.3",
				"description":        "my description",
				"activation_context": "REMOTE",
			},
			expected: "version='1.2.3' and description='my description' and activationContext='REMOTE'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newResourceData(t, tt.raw)
			var call listCall
			client := newMockClient(&call, nil, nil)

			diags := extensionconfig.DataSourceReadWithClient(t.Context(), d, client)

			require.False(t, diags.HasError(), "expected no error, got: %v", diags)
			assert.Equal(t, tt.expected, call.filter)
		})
	}
}
