package reconcile

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func compareAndUpdateAllUsers(ctx context.Context, client permissions.AllUsersClient, objectID string, currentAllUsers, desiredAllUsers string) error {
	if currentAllUsers == desiredAllUsers {
		return nil
	}

	if desiredAllUsers == permissions.HCLAccessorNone {
		_, err := client.DeleteAllUsersAccessor(ctx, objectID, false)
		return err
	}

	if currentAllUsers == permissions.HCLAccessorNone {
		return createAllUsers(ctx, client, objectID, desiredAllUsers)
	}

	return updateAllUsers(ctx, client, objectID, desiredAllUsers)
}

func createAllUsers(ctx context.Context, client permissions.AllUsersClient, objectID string, desiredAllUsers string) error {
	body, err := json.Marshal(permissions.PermissionObject{
		Accessor: permissions.Accessor{
			Type: permissions.AllUsers,
		},
		Permissions: convert.HCLToDTOPermission(desiredAllUsers),
	})
	if err != nil {
		return err
	}

	_, err = client.Create(ctx, objectID, false, body)
	return err
}

func updateAllUsers(ctx context.Context, client permissions.AllUsersClient, objectID string, desiredAllUsers string) error {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(desiredAllUsers),
	})
	if err != nil {
		return err
	}

	_, err = client.UpdateAllUsersAccessor(ctx, objectID, false, body)
	return err
}
