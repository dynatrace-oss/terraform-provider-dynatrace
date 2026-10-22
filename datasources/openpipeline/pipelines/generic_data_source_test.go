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

package pipelines

import (
	"context"
	"errors"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	testutil "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSchemaID = "builtin:test.schema"

type testSettings struct {
	DisplayName string
	marshalErr  error
}

func (s *testSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {Type: schema.TypeString, Required: true},
	}
}

func (s *testSettings) MarshalHCL(properties hcl.Properties) error {
	if s.marshalErr != nil {
		return s.marshalErr
	}
	return properties.Encode("display_name", s.DisplayName)
}

func (s *testSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("display_name", &s.DisplayName)
}

// listFilterServiceRecorder returns a service whose ListWithFilter returns the given result
// and records the filters it was called with.
func listFilterServiceRecorder(stubs api.Stubs, err error) (*testutil.MockListFilterService[*testSettings], *[]string) {
	filters := &[]string{}
	svc := &testutil.MockListFilterService[*testSettings]{
		SchemaIDValue: testSchemaID,
		ListWithFilterFn: func(_ context.Context, filter string) (api.Stubs, error) {
			*filters = append(*filters, filter)
			return stubs, err
		},
	}
	return svc, filters
}

func serviceFnReturning(svc settings.CRUDService[*testSettings]) func(rest.ClientSet) (settings.CRUDService[*testSettings], error) {
	return func(rest.ClientSet) (settings.CRUDService[*testSettings], error) {
		return svc, nil
	}
}

func validMeta() *config.ProviderConfiguration {
	return &config.ProviderConfiguration{EnvironmentURL: "https://abc12345.live.dynatrace.com"}
}

func readDataSource(t *testing.T, resource *schema.Resource, raw map[string]any, meta any) (*schema.ResourceData, diag.Diagnostics) {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resource.Schema, raw)
	return d, resource.ReadContext(t.Context(), d, meta)
}

func TestGenericDataSource_Schema(t *testing.T) {
	resource := genericDataSource(serviceFnReturning(&testutil.MockListFilterService[*testSettings]{}))

	t.Run("custom_id conflicts with display_name and group_role", func(t *testing.T) {
		displayName := resource.Schema["display_name"]
		assert.True(t, displayName.Optional)
		assert.Equal(t, []string{"custom_id"}, displayName.ConflictsWith)

		customID := resource.Schema["custom_id"]
		assert.True(t, customID.Optional)
		assert.Equal(t, []string{"display_name", "group_role"}, customID.ConflictsWith)
	})

	t.Run("group_role conflicts with custom_id only", func(t *testing.T) {
		groupRole := resource.Schema["group_role"]
		assert.True(t, groupRole.Optional)
		assert.Equal(t, schema.TypeString, groupRole.Type)
		assert.Equal(t, []string{"custom_id"}, groupRole.ConflictsWith)
	})

	t.Run("group_role and custom_id are rejected together", func(t *testing.T) {
		cfg := terraform.NewResourceConfigRaw(map[string]any{"custom_id": "my-custom-id", "group_role": "basePipeline"})
		diags := resource.Validate(cfg)
		assert.True(t, diags.HasError())
	})

	t.Run("group_role and display_name are accepted together", func(t *testing.T) {
		cfg := terraform.NewResourceConfigRaw(map[string]any{"display_name": "my-pipeline", "group_role": "basePipeline"})
		diags := resource.Validate(cfg)
		assert.False(t, diags.HasError(), diags)
	})

	t.Run("items contain the computed settings schema and an id", func(t *testing.T) {
		items := resource.Schema["items"]
		assert.True(t, items.Computed)
		assert.Equal(t, schema.TypeList, items.Type)

		elem, ok := items.Elem.(*schema.Resource)
		require.True(t, ok)

		require.Contains(t, elem.Schema, "id")
		assert.True(t, elem.Schema["id"].Computed)

		// The settings' own "display_name" is Required; within items it must be turned into Computed.
		require.Contains(t, elem.Schema, "display_name")
		assert.True(t, elem.Schema["display_name"].Computed)
		assert.False(t, elem.Schema["display_name"].Required)
	})

	t.Run("schema is internally valid", func(t *testing.T) {
		assert.NoError(t, resource.InternalValidate(nil, false))
	})
}

func TestGetFilter(t *testing.T) {
	resource := genericDataSource(serviceFnReturning(&testutil.MockListFilterService[*testSettings]{}))

	tests := []struct {
		name     string
		raw      map[string]any
		expected string
	}{
		{
			name:     "display_name",
			raw:      map[string]any{"display_name": "my-pipeline"},
			expected: "value.displayName='my-pipeline'",
		},
		{
			name:     "custom_id",
			raw:      map[string]any{"custom_id": "my-custom-id"},
			expected: "value.customId='my-custom-id'",
		},
		{
			name:     "group_role",
			raw:      map[string]any{"group_role": "basePipeline"},
			expected: "value.groupRole='basePipeline'",
		},
		{
			name:     "display_name and group_role are combined",
			raw:      map[string]any{"display_name": "my-pipeline", "group_role": "memberPipeline"},
			expected: "value.displayName='my-pipeline' and value.groupRole='memberPipeline'",
		},
		{
			name:     "single quotes are escaped",
			raw:      map[string]any{"display_name": "it's"},
			expected: `value.displayName='it\'s'`,
		},
		{
			// Without escaping the backslash, the trailing "\" would escape the closing quote.
			name:     "backslashes are escaped",
			raw:      map[string]any{"custom_id": `a\b\`},
			expected: `value.customId='a\\b\\'`,
		},
		{
			name:     "group_role is escaped",
			raw:      map[string]any{"group_role": `it's\`},
			expected: `value.groupRole='it\'s\\'`,
		},
		{
			name:     "nothing set",
			raw:      map[string]any{},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resource.Schema, tc.raw)
			assert.Equal(t, tc.expected, getFilter(d))
		})
	}
}

func TestDataSourceRead(t *testing.T) {
	t.Run("sets items and id from the filtered stubs", func(t *testing.T) {
		svc, filters := listFilterServiceRecorder(api.Stubs{
			{ID: "id-1", Value: &testSettings{DisplayName: "pipeline"}},
			{ID: "id-2", Value: &testSettings{DisplayName: "pipeline"}},
		}, nil)

		d, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{"display_name": "pipeline"}, validMeta())

		require.False(t, diags.HasError(), diags)
		assert.Equal(t, []string{"value.displayName='pipeline'"}, *filters)
		assert.Equal(t, "value.displayName='pipeline'", d.Id())
		assert.Equal(t, []any{
			map[string]any{"id": "id-1", "display_name": "pipeline"},
			map[string]any{"id": "id-2", "display_name": "pipeline"},
		}, d.Get("items"))
	})

	t.Run("uses the combined filter for display_name and group_role", func(t *testing.T) {
		svc, filters := listFilterServiceRecorder(api.Stubs{
			{ID: "id-1", Value: &testSettings{DisplayName: "pipeline"}},
		}, nil)

		d, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{"display_name": "pipeline", "group_role": "compositionPipeline"}, validMeta())

		expectedFilter := "value.displayName='pipeline' and value.groupRole='compositionPipeline'"
		require.False(t, diags.HasError(), diags)
		assert.Equal(t, []string{expectedFilter}, *filters)
		assert.Equal(t, expectedFilter, d.Id())
		assert.Equal(t, []any{
			map[string]any{"id": "id-1", "display_name": "pipeline"},
		}, d.Get("items"))
	})

	t.Run("sets empty items if nothing matches", func(t *testing.T) {
		svc, filters := listFilterServiceRecorder(api.Stubs{}, nil)

		d, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{"custom_id": "unknown"}, validMeta())

		require.False(t, diags.HasError(), diags)
		assert.Equal(t, []string{"value.customId='unknown'"}, *filters)
		assert.Equal(t, "value.customId='unknown'", d.Id())
		assert.Empty(t, d.Get("items"))
	})

	t.Run("lists all and uses the schema ID as id if no filter is set", func(t *testing.T) {
		svc, filters := listFilterServiceRecorder(api.Stubs{
			{ID: "id-1", Value: &testSettings{DisplayName: "a"}},
			{ID: "id-2", Value: &testSettings{DisplayName: "b"}},
		}, nil)

		d, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{}, validMeta())

		require.False(t, diags.HasError(), diags)
		assert.Equal(t, []string{""}, *filters)
		// An empty ID would make Terraform treat the data source result as null.
		assert.Equal(t, testSchemaID, d.Id())
		assert.Len(t, d.Get("items"), 2)
	})

	t.Run("returns error if credentials are invalid", func(t *testing.T) {
		svc, filters := listFilterServiceRecorder(api.Stubs{}, nil)

		_, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{"display_name": "pipeline"}, &config.ProviderConfiguration{})

		require.True(t, diags.HasError())
		assert.Contains(t, diags[0].Summary, "No Environment URL has been specified")
		// The service must not be used without valid credentials.
		assert.Empty(t, *filters)
	})

	t.Run("returns error if the service cannot be created", func(t *testing.T) {
		serviceErr := errors.New("service error")
		serviceFn := func(rest.ClientSet) (settings.CRUDService[*testSettings], error) {
			return nil, serviceErr
		}

		_, diags := readDataSource(t, genericDataSource(serviceFn), map[string]any{"display_name": "pipeline"}, validMeta())

		require.True(t, diags.HasError())
		assert.Equal(t, serviceErr.Error(), diags[0].Summary)
	})

	t.Run("returns error if the service does not support list filters", func(t *testing.T) {
		svc := &testutil.MockCRUDService[*testSettings]{SchemaIDValue: testSchemaID}

		_, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{"display_name": "pipeline"}, validMeta())

		require.True(t, diags.HasError())
		assert.Equal(t, settings.ErrListFilterNotImplemented(testSchemaID).Error(), diags[0].Summary)
	})

	t.Run("returns error if listing fails", func(t *testing.T) {
		listErr := errors.New("list error")
		svc, _ := listFilterServiceRecorder(nil, listErr)

		d, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{"display_name": "pipeline"}, validMeta())

		require.True(t, diags.HasError())
		assert.Equal(t, listErr.Error(), diags[0].Summary)
		assert.Empty(t, d.Id())
	})

	t.Run("returns error if marshalling an item fails", func(t *testing.T) {
		marshalErr := errors.New("marshal error")
		svc, _ := listFilterServiceRecorder(api.Stubs{
			{ID: "id-1", Value: &testSettings{marshalErr: marshalErr}},
		}, nil)

		d, diags := readDataSource(t, genericDataSource(serviceFnReturning(svc)), map[string]any{"display_name": "pipeline"}, validMeta())

		require.True(t, diags.HasError())
		assert.Equal(t, marshalErr.Error(), diags[0].Summary)
		assert.Empty(t, d.Id())
	})
}
