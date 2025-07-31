package reconcile

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func compareAndUpdateUsers(ctx context.Context, client permissions.AccessorClient, objectID string, currentUsers, desiredUsers permissions.Users) error {
	errs := make([]error, 0)
	// update and create users
	for _, user := range desiredUsers {
		var err error
		if exists, isEqual := containsUser(currentUsers, user); !exists {
			err = createUser(ctx, client, objectID, user)
		} else if !isEqual {
			err = updateUser(ctx, client, objectID, user)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	// delete users that are not in desiredUsers (HCL)
	for _, user := range currentUsers {
		if exists, _ := containsUser(desiredUsers, user); !exists {
			_, err := client.DeleteAccessor(ctx, objectID, permissions.User, user.UID, false)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func createUser(ctx context.Context, client permissions.AccessorClient, objectID string, user *permissions.UserAccessor) error {
	body, err := json.Marshal(permissions.PermissionObject{
		Accessor: permissions.Accessor{
			Type: permissions.User,
			ID:   user.UID,
		},
		Permissions: convert.HCLToDTOPermission(user.Access),
	})

	if err != nil {
		return err
	}

	_, err = client.Create(ctx, objectID, false, body)
	return err
}

func updateUser(ctx context.Context, client permissions.AccessorClient, objectID string, user *permissions.UserAccessor) error {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(user.Access),
	})
	if err != nil {
		return err
	}

	_, err = client.UpdateAccessor(ctx, objectID, permissions.User, user.UID, false, body)
	return err
}

func containsUser(users permissions.Users, user *permissions.UserAccessor) (exists, equals bool) {
	for _, u := range users {
		if u.UID == user.UID {
			return true, u.Access == user.Access
		}
	}
	return false, false
}
