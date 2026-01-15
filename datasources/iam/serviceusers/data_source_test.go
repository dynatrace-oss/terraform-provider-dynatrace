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

package serviceusers

import (
	"context"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	serviceusers "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// mockCRUDService implements settings.CRUDService for testing
type mockCRUDService struct {
	getFunc  func(ctx context.Context, id string, v *serviceusers.ServiceUser) error
	listFunc func(ctx context.Context) (api.Stubs, error)
}

func (m *mockCRUDService) Get(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
	return m.getFunc(ctx, id, v)
}

func (m *mockCRUDService) Update(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
	panic("unexpected call to Update")
}

func (m *mockCRUDService) List(ctx context.Context) (api.Stubs, error) {
	return m.listFunc(ctx)
}
func (m *mockCRUDService) Create(ctx context.Context, v *serviceusers.ServiceUser) (*api.Stub, error) {
	panic("unexpected call to Create")
}
func (m *mockCRUDService) Delete(ctx context.Context, id string) error {
	panic("unexpected call to Delete")
}
func (m *mockCRUDService) SchemaID() string { return "" }

func TestDataSourceReadWithService(t *testing.T) {
	t.Run("lookup by ID succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"id": "test-uid",
		})

		mockService := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				assert.Equal(t, "test-uid", id)
				v.Name = "Test User"
				v.Email = "test@example.com"
				v.Description = "Test description"
				v.Groups = []string{"group-1"}
				return nil
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.Empty(t, diags)
		assert.Equal(t, "test-uid", d.Id())
		assert.Equal(t, "Test User", d.Get("name"))
		assert.Equal(t, "test@example.com", d.Get("email"))
	})

	t.Run("lookup by ID fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"id": "test-uid",
		})

		mockService := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				return assert.AnError
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("lookup by name succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return api.Stubs{
					{ID: "uid-1", Name: "Other User"},
					{ID: "uid-2", Name: "Test User"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				switch id {
				case "uid-1":
					v.Name = "Other User"
					v.Email = "other@example.com"
					return nil
				case "uid-2":
					v.Name = "Test User"
					v.Email = "test@example.com"
					v.Description = "Found user"
					v.Groups = []string{"group-1"}
					return nil

				default:
					assert.FailNow(t, "unexpected ID in Get:", id)
					return assert.AnError
				}
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.Empty(t, diags)
		assert.Equal(t, "uid-2", d.Id())
		assert.Equal(t, "Test User", d.Get("name"))
		assert.Equal(t, "test@example.com", d.Get("email"))
		assert.Equal(t, "Found user", d.Get("description"))
	})

	t.Run("lookup by name with duplicates fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"name": "Duplicate User",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return api.Stubs{
					{ID: "uid-1", Name: "Duplicate User"},
					{ID: "uid-2", Name: "Duplicate User"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				switch id {
				case "uid-1":
					v.Name = "Duplicate User"
					v.Email = "first@example.com"
					v.Description = "First duplicate"
					return nil
				case "uid-2":
					v.Name = "Duplicate User"
					v.Email = "second@example.com"
					v.Description = "Second duplicate"
					return nil
				default:
					assert.FailNow(t, "unexpected ID in Get:", id)
					return assert.AnError
				}
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)
		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "multiple service users found with name")
		assert.Equal(t, "", d.Id())
	})

	t.Run("lookup by name not found", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"name": "Nonexistent User",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return api.Stubs{
					{ID: "uid-1", Name: "Other User"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				v.Name = "Other User"
				v.Email = "other@example.com"
				return nil
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "no service user found with name")
		assert.Equal(t, "", d.Id())
	})

	t.Run("lookup by email succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"email": "test@example.com",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return api.Stubs{
					{ID: "uid-1", Name: "Test User"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				v.Name = "Test User"
				v.Email = "test@example.com"
				v.Description = "Found by email"
				return nil
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.Empty(t, diags)
		assert.Equal(t, "uid-1", d.Id())
		assert.Equal(t, "test@example.com", d.Get("email"))
		assert.Equal(t, "Found by email", d.Get("description"))
	})

	t.Run("lookup by email not found", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"email": "nonexistent@example.com",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return api.Stubs{
					{ID: "uid-1", Name: "Other User"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				v.Name = "Other User"
				v.Email = "other@example.com"
				return nil
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "no service user found with email")
		assert.Equal(t, "", d.Id())
	})

	t.Run("list fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return nil, assert.AnError
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("get fails during iteration", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return api.Stubs{
					{ID: "uid-1", Name: "Test User"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
				return assert.AnError
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("empty list returns not found", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			listFunc: func(ctx context.Context) (api.Stubs, error) {
				return api.Stubs{}, nil
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "no service user found with name")
	})
}
