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

package latest_version_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/extensions/latest_version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type mockClient struct {
	listExtensionVersions func(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error)
}

func (m *mockClient) ListExtensionVersions(ctx context.Context, extensionName string) (coreapi.PagedListResponse, error) {
	return m.listExtensionVersions(ctx, extensionName)
}

// makePagedListResponse builds a PagedListResponse where each version string
// is marshalled into the {"version": "..."} JSON format expected by the data source.
func makePagedListResponse(versions ...string) coreapi.PagedListResponse {
	var objects [][]byte
	for _, v := range versions {
		b, _ := json.Marshal(map[string]string{"version": v})
		objects = append(objects, b)
	}
	return coreapi.PagedListResponse{
		{Objects: objects},
	}
}

func newResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, latest_version.DataSource().Schema, map[string]any{
		"name": "com.dynatrace.extension.test",
	})
}

// TestDataSourceReadWithClient_ReturnsLatestVersion verifies that among multiple
// installed versions the highest semver is selected.
func TestDataSourceReadWithClient_ReturnsLatestVersion(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		listExtensionVersions: func(_ context.Context, _ string) (coreapi.PagedListResponse, error) {
			return makePagedListResponse("1.0.0", "1.2.0", "1.1.5"), nil
		},
	}

	diags := latest_version.DataSourceReadWithClient(t.Context(), d, client)

	require.False(t, diags.HasError(), "expected no error, got: %v", diags)
	assert.Equal(t, "1.2.0", d.Get("latest_version").(string))
	assert.Equal(t, "com.dynatrace.extension.test", d.Id())
}

// TestDataSourceReadWithClient_SingleVersion verifies that a single version is
// returned correctly.
func TestDataSourceReadWithClient_SingleVersion(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		listExtensionVersions: func(_ context.Context, _ string) (coreapi.PagedListResponse, error) {
			return makePagedListResponse("2.3.1"), nil
		},
	}

	diags := latest_version.DataSourceReadWithClient(t.Context(), d, client)

	require.False(t, diags.HasError(), "expected no error, got: %v", diags)
	assert.Equal(t, "2.3.1", d.Get("latest_version").(string))
}

// TestDataSourceReadWithClient_EmptyResponse verifies that when no versions are
// returned the resource ID is cleared and no error is returned.
func TestDataSourceReadWithClient_EmptyResponse(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		listExtensionVersions: func(_ context.Context, _ string) (coreapi.PagedListResponse, error) {
			return coreapi.PagedListResponse{}, nil
		},
	}

	diags := latest_version.DataSourceReadWithClient(t.Context(), d, client)

	require.False(t, diags.HasError(), "expected no error, got: %v", diags)
	assert.Empty(t, d.Id())
}

// TestDataSourceReadWithClient_ClientError verifies that errors from the client
// are surfaced as diagnostics.
func TestDataSourceReadWithClient_ClientError(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		listExtensionVersions: func(_ context.Context, _ string) (coreapi.PagedListResponse, error) {
			return nil, assert.AnError
		},
	}

	diags := latest_version.DataSourceReadWithClient(t.Context(), d, client)

	require.True(t, diags.HasError(), "expected an error diagnostic, got none")
	assert.ElementsMatch(t, diags, diag.FromErr(assert.AnError), "expected diagnostic summary to contain client error message")
}

// TestDataSourceReadWithClient_SkipsInvalidSemver verifies that version strings
// that are not valid semver are silently skipped while valid ones are still used.
func TestDataSourceReadWithClient_SkipsInvalidSemver(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		listExtensionVersions: func(_ context.Context, _ string) (coreapi.PagedListResponse, error) {
			return makePagedListResponse("not-a-version", "1.0.0", "also-bad"), nil
		},
	}

	diags := latest_version.DataSourceReadWithClient(t.Context(), d, client)

	require.False(t, diags.HasError(), "expected no error, got: %v", diags)
	assert.Equal(t, "1.0.0", d.Get("latest_version").(string))
}

// TestDataSourceReadWithClient_InvalidJSON verifies that malformed JSON in the
// response is surfaced as an error diagnostic.
func TestDataSourceReadWithClient_InvalidJSON(t *testing.T) {
	d := newResourceData(t)
	client := &mockClient{
		listExtensionVersions: func(_ context.Context, _ string) (coreapi.PagedListResponse, error) {
			return coreapi.PagedListResponse{
				{Objects: [][]byte{[]byte("not-valid-json")}},
			}, nil
		},
	}

	diags := latest_version.DataSourceReadWithClient(t.Context(), d, client)

	assert.True(t, diags.HasError(), "expected an error diagnostic from invalid JSON, got none")
}
