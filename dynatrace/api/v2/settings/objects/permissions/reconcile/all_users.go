package reconcile

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func getAllUserUpsert(ctx context.Context, client permissions.AllUsersClient, objectID string, currentAllUsers, desiredAllUsers string) (func(adminAccess bool) error, error) {
	if currentAllUsers == desiredAllUsers {
		// No update needed, current and desired are the same
		return nil, nil
	}

	if desiredAllUsers == permissions.HCLAccessorNone {
		return func(adminAccess bool) error {
			_, err := client.DeleteAllUsersAccessor(ctx, objectID, adminAccess)
			return err
		}, nil
	}

	if currentAllUsers == permissions.HCLAccessorNone {
		return getAllUserCreate(ctx, client, objectID, desiredAllUsers)
	}

	return getAllUserUpdate(ctx, client, objectID, desiredAllUsers)
}

func getAllUserCreate(ctx context.Context, client permissions.AllUsersClient, objectID string, desiredAllUsers string) (func(adminAccess bool) error, error) {
	body, err := json.Marshal(permissions.PermissionObject{
		Accessor: permissions.Accessor{
			Type: permissions.AllUsers,
		},
		Permissions: convert.HCLToDTOPermission(desiredAllUsers),
	})
	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) error {
		_, err = client.Create(ctx, objectID, adminAccess, body)
		return err
	}, nil
}

func getAllUserUpdate(ctx context.Context, client permissions.AllUsersClient, objectID string, desiredAllUsers string) (func(adminAccess bool) error, error) {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(desiredAllUsers),
	})
	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) error {
		_, err = client.UpdateAllUsersAccessor(ctx, objectID, adminAccess, body)
		return err
	}, nil
}
