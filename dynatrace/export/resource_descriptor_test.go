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

package export

import (
	"context"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSettingsWithInsertAfter implements settings.Settings with an InsertAfter field
type mockSettingsWithInsertAfter struct {
	InsertAfter string
	Name        string
}

func (m *mockSettingsWithInsertAfter) MarshalHCL(p hcl.Properties) error {
	return nil
}

func (m *mockSettingsWithInsertAfter) UnmarshalHCL(d hcl.Decoder) error {
	return nil
}

func (m *mockSettingsWithInsertAfter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"insert_after": {Type: schema.TypeString, Optional: true},
		"name":         {Type: schema.TypeString, Optional: true},
	}
}

// mockSettingsWithoutInsertAfter implements settings.Settings without an InsertAfter field
type mockSettingsWithoutInsertAfter struct {
	Name string
}

func (m *mockSettingsWithoutInsertAfter) MarshalHCL(p hcl.Properties) error {
	return nil
}

func (m *mockSettingsWithoutInsertAfter) UnmarshalHCL(d hcl.Decoder) error {
	return nil
}

func (m *mockSettingsWithoutInsertAfter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Optional: true},
	}
}

// mockService is a mock implementation of settings.CRUDService
type mockService[T settings.Settings] struct {
	schemaID string
}

func (m *mockService[T]) List(ctx context.Context) (api.Stubs, error) {
	return nil, nil
}

func (m *mockService[T]) Get(ctx context.Context, id string, v T) error {
	return nil
}

func (m *mockService[T]) SchemaID() string {
	return m.schemaID
}

func (m *mockService[T]) Create(ctx context.Context, v T) (*api.Stub, error) {
	return nil, nil
}

func (m *mockService[T]) Update(ctx context.Context, id string, v T) error {
	return nil
}

func (m *mockService[T]) Delete(ctx context.Context, id string) error {
	return nil
}

func TestAddInsertAfterWeakIDDependencies(t *testing.T) {
	t.Run("adds_weak_id_dependency_for_settings20_with_insert_after", func(t *testing.T) {
		// Create a resource with Settings 2.0 schema that has InsertAfter
		testResType := ResourceType("test_builtin_resource")
		resources := map[ResourceType]ResourceDescriptor{
			testResType: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithInsertAfter]{
						Service: &mockService[*mockSettingsWithInsertAfter]{schemaID: "builtin:test.schema"},
					}
				},
				protoType:    &mockSettingsWithInsertAfter{},
				Dependencies: []Dependency{},
			},
		}

		// Execute
		AddInsertAfterWeakIDDependencies(resources)

		// Verify
		descriptor := resources[testResType]
		require.Len(t, descriptor.Dependencies, 1, "Expected one dependency to be added")
		assert.True(t, descriptor.HasWeakIDDependencyTo(testResType), "Expected weak ID dependency to be added")
	})

	t.Run("adds_weak_id_dependency_for_app_schema_with_insert_after", func(t *testing.T) {
		// Create a resource with app: schema that has InsertAfter
		testResType := ResourceType("test_app_resource")
		resources := map[ResourceType]ResourceDescriptor{
			testResType: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithInsertAfter]{
						Service: &mockService[*mockSettingsWithInsertAfter]{schemaID: "app:dynatrace.test"},
					}
				},
				protoType:    &mockSettingsWithInsertAfter{},
				Dependencies: []Dependency{},
			},
		}

		// Execute
		AddInsertAfterWeakIDDependencies(resources)

		// Verify
		descriptor := resources[testResType]
		require.Len(t, descriptor.Dependencies, 1, "Expected one dependency to be added")
		assert.True(t, descriptor.HasWeakIDDependencyTo(testResType), "Expected weak ID dependency to be added")
	})

	t.Run("skips_non_settings20_schema", func(t *testing.T) {
		// Create a resource with non-Settings 2.0 schema
		testResType := ResourceType("test_legacy_resource")
		resources := map[ResourceType]ResourceDescriptor{
			testResType: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithInsertAfter]{
						Service: &mockService[*mockSettingsWithInsertAfter]{schemaID: "legacy:test.schema"},
					}
				},
				protoType:    &mockSettingsWithInsertAfter{},
				Dependencies: []Dependency{},
			},
		}

		// Execute
		AddInsertAfterWeakIDDependencies(resources)

		// Verify - no dependency should be added for non-Settings 2.0 schema
		descriptor := resources[testResType]
		assert.Empty(t, descriptor.Dependencies, "Expected no dependencies for non-Settings 2.0 schema")
	})

	t.Run("skips_schema_without_insert_after_attribute", func(t *testing.T) {
		// Create a resource without InsertAfter field
		testResType := ResourceType("test_no_insert_after")
		resources := map[ResourceType]ResourceDescriptor{
			testResType: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithoutInsertAfter]{
						Service: &mockService[*mockSettingsWithoutInsertAfter]{schemaID: "builtin:test.schema"},
					}
				},
				protoType:    &mockSettingsWithoutInsertAfter{},
				Dependencies: []Dependency{},
			},
		}

		// Execute
		AddInsertAfterWeakIDDependencies(resources)

		// Verify - no dependency should be added for schema without InsertAfter
		descriptor := resources[testResType]
		assert.Empty(t, descriptor.Dependencies, "Expected no dependencies for schema without InsertAfter")
	})

	t.Run("skips_when_weak_id_dependency_already_exists", func(t *testing.T) {
		// Create a resource that already has a weak ID dependency
		testResType := ResourceType("test_existing_weak_id")
		existingDep := Dependencies.WeakID(testResType)
		resources := map[ResourceType]ResourceDescriptor{
			testResType: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithInsertAfter]{
						Service: &mockService[*mockSettingsWithInsertAfter]{schemaID: "builtin:test.schema"},
					}
				},
				protoType:    &mockSettingsWithInsertAfter{},
				Dependencies: []Dependency{existingDep},
			},
		}

		// Execute
		AddInsertAfterWeakIDDependencies(resources)

		// Verify - should still have only one dependency
		descriptor := resources[testResType]
		require.Len(t, descriptor.Dependencies, 1, "Expected only the existing dependency")
	})

	t.Run("processes_multiple_resources_independently", func(t *testing.T) {
		// Create multiple resources with different configurations
		resTypeWithInsertAfter := ResourceType("test_with_insert_after")
		resTypeWithoutInsertAfter := ResourceType("test_without_insert_after")
		resTypeNonSettings20 := ResourceType("test_non_settings20")

		resources := map[ResourceType]ResourceDescriptor{
			resTypeWithInsertAfter: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithInsertAfter]{
						Service: &mockService[*mockSettingsWithInsertAfter]{schemaID: "builtin:with.insertafter"},
					}
				},
				protoType:    &mockSettingsWithInsertAfter{},
				Dependencies: []Dependency{},
			},
			resTypeWithoutInsertAfter: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithoutInsertAfter]{
						Service: &mockService[*mockSettingsWithoutInsertAfter]{schemaID: "builtin:without.insertafter"},
					}
				},
				protoType:    &mockSettingsWithoutInsertAfter{},
				Dependencies: []Dependency{},
			},
			resTypeNonSettings20: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithInsertAfter]{
						Service: &mockService[*mockSettingsWithInsertAfter]{schemaID: "v1:legacy.schema"},
					}
				},
				protoType:    &mockSettingsWithInsertAfter{},
				Dependencies: []Dependency{},
			},
		}

		// Execute
		AddInsertAfterWeakIDDependencies(resources)

		// Verify each resource was processed correctly
		assert.True(t, resources[resTypeWithInsertAfter].HasWeakIDDependencyTo(resTypeWithInsertAfter),
			"Resource with InsertAfter should have weak ID dependency")
		assert.Empty(t, resources[resTypeWithoutInsertAfter].Dependencies,
			"Resource without InsertAfter should have no dependencies")
		assert.Empty(t, resources[resTypeNonSettings20].Dependencies,
			"Non-Settings 2.0 resource should have no dependencies")
	})

	t.Run("preserves_existing_non_weak_id_dependencies", func(t *testing.T) {
		// Create a resource with existing dependencies that are not weak ID
		testResType := ResourceType("test_with_existing_deps")
		otherResType := ResourceType("other_resource")
		existingDep := Dependencies.ID(otherResType)

		resources := map[ResourceType]ResourceDescriptor{
			testResType: {
				Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
					return &settings.GenericCRUDService[*mockSettingsWithInsertAfter]{
						Service: &mockService[*mockSettingsWithInsertAfter]{schemaID: "builtin:test.schema"},
					}
				},
				protoType:    &mockSettingsWithInsertAfter{},
				Dependencies: []Dependency{existingDep},
			},
		}

		// Execute
		AddInsertAfterWeakIDDependencies(resources)

		// Verify - should have both the existing dep and the new weak ID dep
		descriptor := resources[testResType]
		require.Len(t, descriptor.Dependencies, 2, "Expected two dependencies")
		assert.True(t, descriptor.HasWeakIDDependencyTo(testResType), "Expected weak ID dependency to be added")
	})
}
