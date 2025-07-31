package permissions_test

import (
	"context"
	"encoding/json"
	"testing"

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
		deleteAllUserCalled := 0
		deleteAccessorCalled := 0
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
				deleteAllUserCalled++
				return coreApi.Response{
					StatusCode: 204,
				}, nil
			},
			deleteAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (coreApi.Response, error) {
				deleteAccessorCalled++
				return coreApi.Response{
					StatusCode: 204,
				}, nil
			},
		}
		service := permissionService.ServiceImpl{Client: client}
		err := service.Delete(t.Context(), "objectID")
		assert.NoError(t, err)
		assert.Equal(t, 1, deleteAllUserCalled)
		assert.Equal(t, 2, deleteAccessorCalled)
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
