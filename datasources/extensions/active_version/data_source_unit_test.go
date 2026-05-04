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

package active_version_test

import (
	"context"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/extensions/active_version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type mockClient struct {
	getEnvironmentConfiguration func(ctx context.Context, extensionName string) (coreapi.Response, error)
}

func (m *mockClient) GetEnvironmentConfiguration(ctx context.Context, extensionName string) (coreapi.Response, error) {
	return m.getEnvironmentConfiguration(ctx, extensionName)
}

// makeResponse returns a response with a version
func makeResponse() coreapi.Response {
	return coreapi.Response{
		Data: []byte(`{"version": "2.1.1"}`),
	}
}

func newResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, active_version.DataSource().Schema, map[string]any{
		"name": "com.dynatrace.extension.test",
	})
}

// TestDataSourceReadWithClient_ReturnsActiveVersion verifies that among multiple
// installed versions the highest semver is selected.
func TestDataSourceReadWithClient_ReturnsActiveVersion(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		getEnvironmentConfiguration: func(_ context.Context, _ string) (coreapi.Response, error) {
			return makeResponse(), nil
		},
	}

	diags := active_version.DataSourceReadWithClient(t.Context(), d, client)

	require.False(t, diags.HasError(), "expected no error, got: %v", diags)
	assert.Equal(t, "2.1.1", d.Get("active_version").(string))
	assert.Equal(t, "com.dynatrace.extension.test", d.Id())
}

// TestDataSourceReadWithClient_ClientError verifies that errors from the client
// are surfaced as diagnostics.
func TestDataSourceReadWithClient_ClientError(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		getEnvironmentConfiguration: func(_ context.Context, _ string) (coreapi.Response, error) {
			return coreapi.Response{}, assert.AnError
		},
	}

	diags := active_version.DataSourceReadWithClient(t.Context(), d, client)

	require.True(t, diags.HasError(), "expected an error diagnostic, got none")
	assert.ElementsMatch(t, diags, diag.FromErr(assert.AnError), "expected diagnostic summary to contain client error message")
}

// TestDataSourceReadWithClient_InvalidJSON verifies that malformed JSON in the
// response is surfaced as an error diagnostic.
func TestDataSourceReadWithClient_InvalidJSON(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		getEnvironmentConfiguration: func(_ context.Context, _ string) (coreapi.Response, error) {
			return coreapi.Response{
				Data: []byte("not-valid-json"),
			}, nil
		},
	}

	diags := active_version.DataSourceReadWithClient(t.Context(), d, client)

	assert.True(t, diags.HasError(), "expected an error diagnostic from invalid JSON, got none")
}
