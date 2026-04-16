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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers"
	su "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// mockCRUDService implements settings.CRUDService for testing
type mockCRUDService struct {
	getFunc    func(ctx context.Context, id string, v *su.ServiceUser) error
	getAllFunc func(ctx context.Context) ([]serviceusers.ServiceUserStub, error)
}

func (m *mockCRUDService) Get(ctx context.Context, id string, v *su.ServiceUser) error {
	return m.getFunc(ctx, id, v)
}

func (m *mockCRUDService) GetAll(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
	return m.getAllFunc(ctx)
}

func (m *mockCRUDService) Update(ctx context.Context, id string, v *su.ServiceUser) error {
	panic("unexpected call to Update")
}

func (m *mockCRUDService) List(ctx context.Context) (api.Stubs, error) {
	panic("unexpected call to List")
}
func (m *mockCRUDService) Create(ctx context.Context, v *su.ServiceUser) (*api.Stub, error) {
	panic("unexpected call to Create")
}
func (m *mockCRUDService) Delete(ctx context.Context, id string) error {
	panic("unexpected call to Delete")
}
func (m *mockCRUDService) SchemaID() string { return "" }

func TestDataSourceReadWithService(t *testing.T) {
	t.Run("lookup by ID succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"id": "test-uid",
		})

		mockService := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *su.ServiceUser) error {
				assert.Equal(t, "test-uid", id)
				v.Name = "Test User"
				v.Email = "test@example.com"
				v.Description = "Test description"
				v.Groups = []string{"group-1"}
				return nil
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.Empty(t, diags)
		assert.Equal(t, "test-uid", d.Id())
		assert.Equal(t, "Test User", d.Get("name"))
		assert.Equal(t, "test@example.com", d.Get("email"))
		assert.Equal(t, "Test description", d.Get("description"))

		groups := d.Get("groups").(*schema.Set)
		assert.NotNil(t, groups)
		assert.Equal(t, []any{"group-1"}, groups.List())
	})

	t.Run("lookup by ID fails if get fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"id": "test-uid",
		})

		mockService := &mockCRUDService{
			getFunc: func(ctx context.Context, id string, v *su.ServiceUser) error {
				return assert.AnError
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("lookup by name succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{
					{UID: "uid-1", Name: "Other User", Email: "other@example.com", Description: "Not the one"},
					{UID: "uid-2", Name: "Test User", Email: "test@example.com", Description: "Found user"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *su.ServiceUser) error {
				if id == "uid-2" {
					v.Name = "Test User"
					v.Email = "test@example.com"
					v.Description = "Found user"
					v.Groups = []string{"group-1"}
					return nil
				}

				assert.FailNow(t, "unexpected ID in Get:", id)
				return assert.AnError
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.Empty(t, diags)
		assert.Equal(t, "uid-2", d.Id())
		assert.Equal(t, "Test User", d.Get("name"))
		assert.Equal(t, "test@example.com", d.Get("email"))
		assert.Equal(t, "Found user", d.Get("description"))
		groups := d.Get("groups").(*schema.Set)
		assert.NotNil(t, groups)
		assert.Equal(t, []any{"group-1"}, groups.List())
	})

	t.Run("lookup by name with duplicates fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"name": "Duplicate User",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{
					{UID: "uid-1", Name: "Duplicate User", Email: "first@example.com", Description: "First duplicate"},
					{UID: "uid-2", Name: "Duplicate User", Email: "second@example.com", Description: "Second duplicate"},
				}, nil
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)
		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "multiple service users found with name")
		assert.Equal(t, "", d.Id())
	})

	t.Run("lookup by name not found", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"name": "Nonexistent User",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{
					{UID: "uid-1", Name: "Other User", Email: "other@example.com", Description: "Not the one"},
				}, nil
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "no service user found with name")
		assert.Equal(t, "", d.Id())
	})

	t.Run("lookup by name fails if get all fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return nil, assert.AnError
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("lookup by name fails if get fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{
					{UID: "uid-1", Name: "Test User", Email: "test@example.com", Description: "Found by email"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *su.ServiceUser) error {
				return assert.AnError
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("lookup by name fails if get all returns empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"name": "Test User",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{}, nil
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "no service user found with name")
	})

	t.Run("lookup by email succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"email": "test@example.com",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{
					{UID: "uid-1", Name: "Test User", Email: "test@example.com", Description: "Found by email"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *su.ServiceUser) error {
				v.Name = "Test User"
				v.Email = "test@example.com"
				v.Description = "Found by email"
				v.Groups = []string{"group-1"}
				return nil
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.Empty(t, diags)
		assert.Equal(t, "uid-1", d.Id())
		assert.Equal(t, "test@example.com", d.Get("email"))
		assert.Equal(t, "Found by email", d.Get("description"))
		groups := d.Get("groups").(*schema.Set)
		assert.NotNil(t, groups)
		assert.Equal(t, []any{"group-1"}, groups.List())
	})

	t.Run("lookup by email not found", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"email": "nonexistent@example.com",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{
					{UID: "uid-1", Name: "Other User", Email: "other@example.com", Description: "Not the one"},
				}, nil
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "no service user found with email")
		assert.Equal(t, "", d.Id())
	})

	t.Run("lookup by email fails if get all fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"email": "test@example.com",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return nil, assert.AnError
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("lookup by email fails if get fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"email": "test@example.com",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{
					{UID: "uid-1", Name: "Test User", Email: "test@example.com", Description: "Found by email"},
				}, nil
			},
			getFunc: func(ctx context.Context, id string, v *su.ServiceUser) error {
				return assert.AnError
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, assert.AnError.Error())
	})

	t.Run("lookup by email fails if get all returns empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{
			"email": "test@example.com",
		})

		mockService := &mockCRUDService{
			getAllFunc: func(ctx context.Context) ([]serviceusers.ServiceUserStub, error) {
				return []serviceusers.ServiceUserStub{}, nil
			},
		}

		diags := dataSourceReadWithService(t.Context(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "no service user found with email")
	})
}
