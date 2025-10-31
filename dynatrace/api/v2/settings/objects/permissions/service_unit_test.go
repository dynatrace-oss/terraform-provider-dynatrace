//go:build unit

/*
 * @license
 * Copyright 2025 Dynatrace LLC
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

package permissions_test

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	permissionService "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreApi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type clientStub struct {
	getAllAccessor         func(context.Context, string, bool) (coreApi.Response, error)
	create                 func(ctx context.Context, objectID string, adminAccess bool, body []byte) (coreApi.Response, error)
	updateAccessor         func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (coreApi.Response, error)
	deleteAccessor         func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (coreApi.Response, error)
	updateAllUsersAccessor func(ctx context.Context, objectID string, adminAccess bool, body []byte) (coreApi.Response, error)
	deleteAllUsersAccessor func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error)
}

func (c *clientStub) GetAllAccessors(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
	return c.getAllAccessor(ctx, objectID, adminAccess)
}

func (c *clientStub) Create(ctx context.Context, objectID string, adminAccess bool, body []byte) (coreApi.Response, error) {
	return c.create(ctx, objectID, adminAccess, body)
}

func (c *clientStub) UpdateAccessor(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (coreApi.Response, error) {
	return c.updateAccessor(ctx, objectID, accessorType, accessorID, adminAccess, body)
}

func (c *clientStub) DeleteAccessor(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (coreApi.Response, error) {
	return c.deleteAccessor(ctx, objectID, accessorType, accessorID, adminAccess)
}

func (c *clientStub) UpdateAllUsersAccessor(ctx context.Context, objectID string, adminAccess bool, body []byte) (coreApi.Response, error) {
	return c.updateAllUsersAccessor(ctx, objectID, adminAccess, body)
}

func (c *clientStub) DeleteAllUsersAccessor(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
	return c.deleteAllUsersAccessor(ctx, objectID, adminAccess)
}

type settingsClientStub struct {
	getSchemaIDsWithOwnerBasedAccessControl func(ctx context.Context) ([]string, error)
	listObjectsIDsOfSchema                  func(ctx context.Context, schemaID string) ([]string, error)
}

// Implement exporter.SettingsClient interface
func (c *settingsClientStub) GetSchemaIDsWithOwnerBasedAccessControl(ctx context.Context) ([]string, error) {
	return c.getSchemaIDsWithOwnerBasedAccessControl(ctx)
}

func (c *settingsClientStub) ListObjectsIDsOfSchema(ctx context.Context, schemaID string) ([]string, error) {
	return c.listObjectsIDsOfSchema(ctx, schemaID)
}

func TestService(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("It gets permissions", func(t *testing.T) {
			client := &clientStub{
				getAllAccessor: func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
					response := permissions.PermissionObjects{
						Accessors: []permissions.PermissionObject{
							{
								Accessor: permissions.Accessor{
									Type: permissions.AllUsers,
								},
								Permissions: []permissions.TypePermissions{permissions.Read, permissions.Write},
							},
							{
								Accessor: permissions.Accessor{
									Type: permissions.User,
									ID:   "userID",
								},
								Permissions: []permissions.TypePermissions{permissions.Read, permissions.Write},
							},
							{
								Accessor: permissions.Accessor{
									Type: permissions.Group,
									ID:   "groupID",
								},
								Permissions: []permissions.TypePermissions{permissions.Read},
							},
						},
					}
					responseBytes, err := json.Marshal(response)
					require.NoError(t, err)

					return coreApi.Response{
						StatusCode: 200,
						Data:       responseBytes,
					}, nil
				},
			}
			service := permissionService.ServiceImpl{Client: client}
			value := permissions.SettingPermissions{}
			err := service.Get(t.Context(), "objectID", &value)

			assert.NoError(t, err)
			assert.Equal(t, permissions.SettingPermissions{
				SettingsObjectID: "objectID",
				AllUsers:         permissions.HCLAccessorWrite,
				Users: permissions.Users{
					{
						UID:    "userID",
						Access: permissions.HCLAccessorWrite,
					},
				},
				Groups: permissions.Groups{
					{
						ID:     "groupID",
						Access: permissions.HCLAccessorRead,
					},
				},
			}, value)
		})

		t.Run("Errors during get", func(t *testing.T) {
			client := &clientStub{
				getAllAccessor: func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
					return coreApi.Response{}, assert.AnError
				},
			}
			service := permissionService.ServiceImpl{Client: client}
			value := permissions.SettingPermissions{
				SettingsObjectID: "objectID",
			}
			err := service.Get(t.Context(), "objectID", &value)
			assert.Error(t, err)
			assert.ErrorIs(t, assert.AnError, err)
		})

		t.Run("It errors during client creation", func(t *testing.T) {
			service := permissionService.ServiceImpl{
				Credentials: &rest.Credentials{},
			}
			err := service.Get(t.Context(), "objectID", &permissions.SettingPermissions{})
			assert.ErrorIs(t, rest.NoPlatformCredentialsErr, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("Returns stubs for all object IDs", func(t *testing.T) {
			settingsClient := &settingsClientStub{
				getSchemaIDsWithOwnerBasedAccessControl: func(ctx context.Context) ([]string, error) {
					return []string{"schemaID1", "schemaID2"}, nil
				},
				listObjectsIDsOfSchema: func(ctx context.Context, schemaID string) ([]string, error) {
					if schemaID == "schemaID1" {
						return []string{"objectID1", "objectID2"}, nil
					}
					if schemaID == "schemaID2" {
						return []string{"objectID3"}, nil
					}
					return nil, nil
				},
			}
			service := permissionService.ServiceImpl{SettingsClient: settingsClient}
			stubs, err := service.List(t.Context())
			require.NoError(t, err)
			assert.Len(t, stubs, 3)
			assert.Equal(t, &api.Stub{ID: "objectID1", Name: "objectID1"}, stubs[0])
			assert.Equal(t, &api.Stub{ID: "objectID2", Name: "objectID2"}, stubs[1])
			assert.Equal(t, &api.Stub{ID: "objectID3", Name: "objectID3"}, stubs[2])
		})

		t.Run("Returns error if GetSchemaIDsWithOwnerBasedAccessControl fails", func(t *testing.T) {
			settingsClient := &settingsClientStub{
				getSchemaIDsWithOwnerBasedAccessControl: func(ctx context.Context) ([]string, error) {
					return nil, assert.AnError
				},
			}
			service := permissionService.ServiceImpl{SettingsClient: settingsClient}
			stubs, err := service.List(t.Context())
			assert.ErrorIs(t, err, assert.AnError)
			assert.Nil(t, stubs)
		})

		t.Run("Returns error if ListObjectsIDsOfSchema fails", func(t *testing.T) {
			settingsClient := &settingsClientStub{
				getSchemaIDsWithOwnerBasedAccessControl: func(ctx context.Context) ([]string, error) {
					return []string{"schemaID1"}, nil
				},
				listObjectsIDsOfSchema: func(ctx context.Context, schemaID string) ([]string, error) {
					return nil, assert.AnError
				},
			}
			service := permissionService.ServiceImpl{SettingsClient: settingsClient}
			stubs, err := service.List(t.Context())
			assert.ErrorIs(t, err, assert.AnError)
			assert.Nil(t, stubs)
		})

		t.Run("Returns error if getSettingsClient fails", func(t *testing.T) {
			service := permissionService.ServiceImpl{
				Credentials: &rest.Credentials{},
			}
			stubs, err := service.List(t.Context())
			assert.ErrorIs(t, err, rest.NoPlatformCredentialsErr)
			assert.Nil(t, stubs)
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("It successfully upserts (no changes)", func(t *testing.T) {
			client := &clientStub{
				getAllAccessor: func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
					response := permissions.PermissionObjects{
						Accessors: []permissions.PermissionObject{},
					}
					responseBytes, err := json.Marshal(response)
					require.NoError(t, err)

					return coreApi.Response{
						StatusCode: 200,
						Data:       responseBytes,
					}, nil
				},
			}
			service := permissionService.ServiceImpl{Client: client}
			value := permissions.SettingPermissions{
				SettingsObjectID: "objectID",
				AllUsers:         permissions.HCLAccessorNone,
				Users:            permissions.Users{},
				Groups:           permissions.Groups{},
			}
			_, err := service.Create(t.Context(), &value)
			assert.NoError(t, err)
		})

		t.Run("It errors during upsert", func(t *testing.T) {
			client := &clientStub{
				updateAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (coreApi.Response, error) {
					return coreApi.Response{}, assert.AnError
				},
				getAllAccessor: func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
					response := permissions.PermissionObjects{
						Accessors: []permissions.PermissionObject{
							{
								Accessor: permissions.Accessor{
									Type: permissions.AllUsers,
								},
								Permissions: []permissions.TypePermissions{permissions.Read, permissions.Write},
							},
						},
					}
					responseBytes, err := json.Marshal(response)
					require.NoError(t, err)

					return coreApi.Response{
						StatusCode: 200,
						Data:       responseBytes,
					}, nil
				},
			}
			service := permissionService.ServiceImpl{Client: client}
			value := permissions.SettingPermissions{
				SettingsObjectID: "objectID",
				AllUsers:         permissions.HCLAccessorRead,
			}
			_, err := service.Create(t.Context(), &value)
			assert.ErrorContains(t, err, assert.AnError.Error())
		})

		t.Run("It errors during get", func(t *testing.T) {
			client := &clientStub{
				getAllAccessor: func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
					return coreApi.Response{}, assert.AnError
				},
			}
			service := permissionService.ServiceImpl{Client: client}
			value := permissions.SettingPermissions{
				SettingsObjectID: "objectID",
			}
			_, err := service.Create(t.Context(), &value)
			assert.Error(t, err)
			assert.ErrorIs(t, assert.AnError, err)
		})

		t.Run("It errors during client creation", func(t *testing.T) {
			service := permissionService.ServiceImpl{
				Credentials: &rest.Credentials{},
			}
			_, err := service.Create(t.Context(), &permissions.SettingPermissions{})
			assert.ErrorIs(t, rest.NoPlatformCredentialsErr, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("It errors during client creation", func(t *testing.T) {
			service := permissionService.ServiceImpl{
				Credentials: &rest.Credentials{},
			}
			err := service.Update(t.Context(), "", &permissions.SettingPermissions{})
			assert.ErrorIs(t, rest.NoPlatformCredentialsErr, err)
		})
	})

	t.Run("Delete deletes every permission", func(t *testing.T) {
		var deleteAllUserCalled, deleteAccessorCalled atomic.Int32
		client := &clientStub{
			getAllAccessor: func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
				response := permissions.PermissionObjects{
					Accessors: []permissions.PermissionObject{
						{
							Accessor: permissions.Accessor{
								Type: permissions.AllUsers,
							},
							Permissions: []permissions.TypePermissions{permissions.Read, permissions.Write},
						},
						{
							Accessor: permissions.Accessor{
								Type: permissions.User,
								ID:   "userID",
							},
							Permissions: []permissions.TypePermissions{permissions.Read, permissions.Write},
						},
						{
							Accessor: permissions.Accessor{
								Type: permissions.Group,
								ID:   "groupID",
							},
							Permissions: []permissions.TypePermissions{permissions.Read},
						},
					},
				}
				responseBytes, err := json.Marshal(response)
				require.NoError(t, err)

				return coreApi.Response{
					StatusCode: 200,
					Data:       responseBytes,
				}, nil
			},
			deleteAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool) (coreApi.Response, error) {
				deleteAllUserCalled.Add(1)
				return coreApi.Response{
					StatusCode: 204,
				}, nil
			},
			deleteAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (coreApi.Response, error) {
				deleteAccessorCalled.Add(1)
				return coreApi.Response{
					StatusCode: 204,
				}, nil
			},
		}
		service := permissionService.ServiceImpl{Client: client}
		err := service.Delete(t.Context(), "objectID")
		assert.NoError(t, err)
		assert.Equal(t, int32(1), deleteAllUserCalled.Load())
		assert.Equal(t, int32(2), deleteAccessorCalled.Load())
	})

	t.Run("Service returns a new instance", func(t *testing.T) {
		service := permissionService.Service(nil)
		assert.IsType(t, &permissionService.ServiceImpl{}, service)
	})

	t.Run("Returns the schema ID", func(t *testing.T) {
		service := permissionService.ServiceImpl{}
		assert.Equal(t, "settings:permissions", service.SchemaID())
	})
}
