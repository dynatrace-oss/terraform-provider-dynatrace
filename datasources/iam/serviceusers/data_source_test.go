//go:build unit

/**
* @license
* Copyright 2025 Dynatrace LLC
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
	su "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// mockCRUDService implements settings.CRUDService for testing
type mockCRUDService struct {
	getFunc    func(ctx context.Context, id string, v *serviceusers.ServiceUser) error
	updateFunc func(ctx context.Context, id string, v *serviceusers.ServiceUser) error
	listFunc   func(ctx context.Context) (api.Stubs, error)
	createFunc func(ctx context.Context, v *serviceusers.ServiceUser) (*api.Stub, error)
	deleteFunc func(ctx context.Context, id string) error
}

func (m *mockCRUDService) Get(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
	return m.getFunc(ctx, id, v)
}

func (m *mockCRUDService) Update(ctx context.Context, id string, v *serviceusers.ServiceUser) error {
	return m.updateFunc(ctx, id, v)
}

func (m *mockCRUDService) List(ctx context.Context) (api.Stubs, error) {
	return m.listFunc(ctx)
}
func (m *mockCRUDService) Create(ctx context.Context, v *serviceusers.ServiceUser) (*api.Stub, error) {
	return m.createFunc(ctx, v)
}
func (m *mockCRUDService) Delete(ctx context.Context, id string) error {
	return m.deleteFunc(ctx, id)
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
				v.UID = "test-uid"
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
					v.UID = "uid-1"
					v.Name = "Other User"
					v.Email = "other@example.com"
					return nil
				case "uid-2":
					v.UID = "uid-2"
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
				v.UID = "uid-1"
				v.Name = "Other User"
				v.Email = "other@example.com"
				return nil
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "Service user not found")
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
				v.UID = "uid-1"
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
				v.UID = "uid-1"
				v.Name = "Other User"
				v.Email = "other@example.com"
				return nil
			},
		}

		diags := dataSourceReadWithService(context.Background(), d, mockService)

		assert.NotEmpty(t, diags)
		assert.Contains(t, diags[0].Summary, "Service user not found")
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
		assert.Contains(t, diags[0].Summary, "Service user not found")
	})
}

func TestSetServiceUserData(t *testing.T) {
	t.Run("sets all fields correctly", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{})

		serviceUser := &su.ServiceUser{
			UID:         "test-uid",
			Name:        "Test User",
			Email:       "test@example.com",
			Description: "Test description",
			Groups:      []string{"group-b", "group-a"},
		}

		diags := setServiceUserData(d, serviceUser)

		assert.Empty(t, diags)
		assert.Equal(t, "test-uid", d.Id())
		assert.Equal(t, "test-uid", d.Get("id"))
		assert.Equal(t, "Test User", d.Get("name"))
		assert.Equal(t, "test@example.com", d.Get("email"))
		assert.Equal(t, "Test description", d.Get("description"))

		groups := d.Get("groups").([]interface{})
		assert.Len(t, groups, 2)
		// Groups should be sorted
		assert.Equal(t, "group-a", groups[0])
		assert.Equal(t, "group-b", groups[1])
	})

	t.Run("sets empty groups as empty list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{})

		serviceUser := &su.ServiceUser{
			UID:         "test-uid",
			Name:        "Test User",
			Email:       "test@example.com",
			Description: "",
			Groups:      []string{},
		}

		diags := setServiceUserData(d, serviceUser)

		assert.Empty(t, diags)
		groups := d.Get("groups").([]interface{})
		assert.Empty(t, groups)
	})

	t.Run("sets nil groups as empty list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{})

		serviceUser := &su.ServiceUser{
			UID:         "test-uid",
			Name:        "Test User",
			Email:       "test@example.com",
			Description: "",
			Groups:      nil,
		}

		diags := setServiceUserData(d, serviceUser)

		assert.Empty(t, diags)
		groups := d.Get("groups").([]interface{})
		assert.Empty(t, groups)
	})

	t.Run("handles empty description", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]interface{}{})

		serviceUser := &su.ServiceUser{
			UID:         "test-uid",
			Name:        "Test User",
			Email:       "test@example.com",
			Description: "",
			Groups:      []string{"group-1"},
		}

		diags := setServiceUserData(d, serviceUser)

		assert.Empty(t, diags)
		assert.Equal(t, "", d.Get("description"))
	})
}
