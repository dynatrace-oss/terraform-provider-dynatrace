package reconcile_test

import (
	"context"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/reconcile"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type clientStub struct {
	create                 func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error)
	updateAccessor         func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error)
	deleteAccessor         func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (api.Response, error)
	updateAllUsersAccessor func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error)
	deleteAllUsersAccessor func(ctx context.Context, objectID string, adminAccess bool) (api.Response, error)
}

func (c clientStub) Create(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
	return c.create(ctx, objectID, adminAccess, body)
}

func (c clientStub) UpdateAccessor(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error) {
	return c.updateAccessor(ctx, objectID, accessorType, accessorID, adminAccess, body)
}

func (c clientStub) DeleteAccessor(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (api.Response, error) {
	return c.deleteAccessor(ctx, objectID, accessorType, accessorID, adminAccess)
}

func (c clientStub) UpdateAllUsersAccessor(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
	return c.updateAllUsersAccessor(ctx, objectID, adminAccess, body)
}

func (c clientStub) DeleteAllUsersAccessor(ctx context.Context, objectID string, adminAccess bool) (api.Response, error) {
	return c.deleteAllUsersAccessor(ctx, objectID, adminAccess)
}

func TestCompareAndUpdate(t *testing.T) {
	t.Run("Successful users operations", func(t *testing.T) {
		t.Run("Create new user", func(t *testing.T) {
			createCalls := 0
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					// Simulate successful creation
					createCalls++
					return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
				},
			}

			current := &permissions.SettingPermissions{
				Users: permissions.Users{},
			}
			desired := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorWrite},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, createCalls, "Expected one user to be created")
		})

		t.Run("Update existing user", func(t *testing.T) {
			updateCalls := 0
			client := clientStub{
				updateAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error) {
					// Simulate successful update
					updateCalls++
					return api.Response{StatusCode: 200, Data: []byte("{}")}, nil
				},
			}

			current := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorWrite},
				},
			}
			desired := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorRead},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, updateCalls, "Expected one user to be updated")
		})

		t.Run("Delete existing user", func(t *testing.T) {
			deleteCalls := 0
			client := clientStub{
				deleteAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (api.Response, error) {
					deleteCalls++
					return api.Response{StatusCode: 204, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorWrite},
				},
			}
			desired := &permissions.SettingPermissions{
				Users: permissions.Users{},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, deleteCalls, "Expected one user to be deleted")
		})

		t.Run("Does not update user if not changed", func(t *testing.T) {
			client := clientStub{
				updateAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error) {
					t.FailNow() // Should not be called
					return api.Response{}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorWrite},
				},
			}
			desired := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorWrite},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
		})
	})

	t.Run("Successful groups operations", func(t *testing.T) {
		t.Run("Create new group", func(t *testing.T) {
			createCalls := 0
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					createCalls++
					return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Groups: permissions.Groups{},
			}
			desired := &permissions.SettingPermissions{
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorWrite},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, createCalls, "Expected one group to be created")
		})

		t.Run("Update existing group", func(t *testing.T) {
			updateCalls := 0
			client := clientStub{
				updateAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error) {
					updateCalls++
					return api.Response{StatusCode: 200, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorWrite},
				},
			}
			desired := &permissions.SettingPermissions{
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorRead},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, updateCalls, "Expected one group to be updated")
		})

		t.Run("Delete existing group", func(t *testing.T) {
			deleteCalls := 0
			client := clientStub{
				deleteAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (api.Response, error) {
					deleteCalls++
					return api.Response{StatusCode: 204, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorWrite},
				},
			}
			desired := &permissions.SettingPermissions{
				Groups: permissions.Groups{},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, deleteCalls, "Expected one group to be deleted")
		})

		t.Run("Does not update group if not changed", func(t *testing.T) {
			client := clientStub{
				updateAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error) {
					t.FailNow() // Should not be called
					return api.Response{}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorWrite},
				},
			}
			desired := &permissions.SettingPermissions{
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorWrite},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
		})
	})

	t.Run("Successful allUser operations", func(t *testing.T) {
		t.Run("Create allUser", func(t *testing.T) {
			createCalls := 0
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					createCalls++
					return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorNone,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, createCalls, "Expected allUser to be created")
		})

		t.Run("Update allUser", func(t *testing.T) {
			updateCalls := 0
			client := clientStub{
				updateAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					updateCalls++
					return api.Response{StatusCode: 200, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorRead,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, updateCalls, "Expected allUser to be updated")
		})

		t.Run("Delete allUser", func(t *testing.T) {
			deleteCalls := 0
			client := clientStub{
				deleteAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool) (api.Response, error) {
					deleteCalls++
					return api.Response{StatusCode: 204, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorNone,
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
			assert.Equal(t, 1, deleteCalls, "Expected allUser to be deleted")
		})

		t.Run("Does not update allUser if not changed", func(t *testing.T) {
			client := clientStub{
				updateAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					t.FailNow() // Should not be called
					return api.Response{}, nil
				},
			}
			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.NoError(t, err)
		})
	})

	t.Run("Handles errors correctly", func(t *testing.T) {
		t.Run("Continues if one user errors", func(t *testing.T) {
			createCalls := 0
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					createCalls++
					if createCalls == 1 {
						return api.Response{}, assert.AnError
					}
					return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Users: permissions.Users{},
			}
			desired := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorWrite},
					&permissions.UserAccessor{UID: "user2", Access: permissions.HCLAccessorWrite},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.Error(t, err)
			assert.Equal(t, 2, createCalls, "Should attempt both users")
		})

		t.Run("Continues if one group errors", func(t *testing.T) {
			createCalls := 0
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					createCalls++
					if createCalls == 1 {
						return api.Response{}, assert.AnError
					}
					return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
				},
			}
			current := &permissions.SettingPermissions{
				Groups: permissions.Groups{},
			}
			desired := &permissions.SettingPermissions{
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorWrite},
					&permissions.GroupAccessor{ID: "group2", Access: permissions.HCLAccessorWrite},
				},
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.Error(t, err)
			assert.Equal(t, 2, createCalls, "Should attempt both groups")
		})

		t.Run("Continues if one allUser errors", func(t *testing.T) {
			calls := 0
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					calls++
					return api.Response{}, assert.AnError
				},
			}
			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorNone,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.Error(t, err)
			assert.Equal(t, 1, calls, "Should attempt allUser")
		})

		t.Run("Returns error if all operations fail", func(t *testing.T) {
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					return api.Response{}, assert.AnError
				},
				updateAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error) {
					return api.Response{}, assert.AnError
				},
				deleteAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (api.Response, error) {
					return api.Response{}, assert.AnError
				},
				updateAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					return api.Response{}, assert.AnError
				},
				deleteAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool) (api.Response, error) {
					return api.Response{}, assert.AnError
				},
			}
			current := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorRead},
				},
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorRead},
				},
				AllUsers: permissions.HCLAccessorRead,
			}
			desired := &permissions.SettingPermissions{
				Users: permissions.Users{
					&permissions.UserAccessor{UID: "user2", Access: permissions.HCLAccessorWrite},
				},
				Groups: permissions.Groups{
					&permissions.GroupAccessor{ID: "group2", Access: permissions.HCLAccessorWrite},
				},
				AllUsers: permissions.HCLAccessorWrite,
			}
			err := reconcile.CompareAndUpdate(t.Context(), client, current, desired, false)
			assert.Error(t, err)
		})
	})

	t.Run("Handles admin access correctly for accessors", func(t *testing.T) {
		client := clientStub{
			create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
				require.True(t, adminAccess, "Expected admin access to be true")
				return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
			},
			updateAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool, body []byte) (api.Response, error) {
				require.True(t, adminAccess, "Expected admin access to be true")
				return api.Response{StatusCode: 200, Data: []byte("{}")}, nil
			},
			deleteAccessor: func(ctx context.Context, objectID string, accessorType string, accessorID string, adminAccess bool) (api.Response, error) {
				require.True(t, adminAccess, "Expected admin access to be true")
				return api.Response{StatusCode: 204, Data: []byte("{}")}, nil
			},
		}

		current := &permissions.SettingPermissions{
			Users: permissions.Users{
				&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorRead},
				&permissions.UserAccessor{UID: "user2", Access: permissions.HCLAccessorRead},
			},
			Groups: permissions.Groups{
				&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorRead},
				&permissions.GroupAccessor{ID: "group2", Access: permissions.HCLAccessorRead},
			},
		}
		desired := &permissions.SettingPermissions{
			Users: permissions.Users{
				// update
				&permissions.UserAccessor{UID: "user1", Access: permissions.HCLAccessorWrite},
				// delete user2
				// create
				&permissions.UserAccessor{UID: "user3", Access: permissions.HCLAccessorWrite},
			},
			Groups: permissions.Groups{
				// update
				&permissions.GroupAccessor{ID: "group1", Access: permissions.HCLAccessorWrite},
				// delete group2
				// create
				&permissions.GroupAccessor{ID: "group3", Access: permissions.HCLAccessorWrite},
			},
		}
		err := reconcile.CompareAndUpdate(context.Background(), client, current, desired, true)
		assert.NoError(t, err, "Expected no error during admin access operations")
	})

	t.Run("Handle adminAccess correctly for allUsers", func(t *testing.T) {
		t.Run("Updates allUsers with admin access", func(t *testing.T) {
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					require.True(t, adminAccess, "Expected admin access to be true")
					return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
				},
				updateAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					require.True(t, adminAccess, "Expected admin access to be true")
					return api.Response{StatusCode: 200, Data: []byte("{}")}, nil
				},
				deleteAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool) (api.Response, error) {
					require.True(t, adminAccess, "Expected admin access to be true")
					return api.Response{StatusCode: 204, Data: []byte("{}")}, nil
				},
			}

			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorNone,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			err := reconcile.CompareAndUpdate(context.Background(), client, current, desired, true)
			assert.NoError(t, err, "Expected no error during admin access operations")
		})

		t.Run("Deletes allUsers with admin access", func(t *testing.T) {
			client := clientStub{
				deleteAllUsersAccessor: func(ctx context.Context, objectID string, adminAccess bool) (api.Response, error) {
					require.True(t, adminAccess, "Expected admin access to be true")
					return api.Response{StatusCode: 204, Data: []byte("{}")}, nil
				},
			}

			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorNone,
			}
			err := reconcile.CompareAndUpdate(context.Background(), client, current, desired, true)
			assert.NoError(t, err, "Expected no error during admin access operations")
		})

		t.Run("Creates allUsers with admin access", func(t *testing.T) {
			client := clientStub{
				create: func(ctx context.Context, objectID string, adminAccess bool, body []byte) (api.Response, error) {
					require.True(t, adminAccess, "Expected admin access to be true")
					return api.Response{StatusCode: 201, Data: []byte("{}")}, nil
				},
			}

			current := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorNone,
			}
			desired := &permissions.SettingPermissions{
				AllUsers: permissions.HCLAccessorWrite,
			}
			err := reconcile.CompareAndUpdate(context.Background(), client, current, desired, true)
			assert.NoError(t, err, "Expected no error during admin access operations")
		})
	})
}
