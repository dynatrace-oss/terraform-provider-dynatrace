package reconcile

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func getUserUpserts(ctx context.Context, client permissions.AccessorClient, objectID string, currentUsers, desiredUsers permissions.Users) ([]func(adminAccess bool) error, error) {
	errs := make([]error, 0)
	updates := make([]func(adminAccess bool) error, 0)
	// update and create users
	for _, user := range desiredUsers {
		var err error
		var fn func(adminAccess bool) error
		if exists, isEqual := containsUser(currentUsers, user); !exists {
			fn, err = getUserCreate(ctx, client, objectID, user)
		} else if !isEqual {
			fn, err = getUserUpdate(ctx, client, objectID, user)
		} else {
			// user already exists and is equal, no action needed
			continue
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}
		updates = append(updates, fn)
	}

	// delete users that are not in desiredUsers (HCL)
	for _, user := range currentUsers {
		if exists, _ := containsUser(desiredUsers, user); !exists {
			updates = append(updates, func(adminAccess bool) error {
				_, err := client.DeleteAccessor(ctx, objectID, permissions.User, user.UID, adminAccess)
				return err
			})
		}
	}
	joinErr := errors.Join(errs...)
	if joinErr != nil {
		return nil, joinErr
	}
	return updates, nil
}

func getUserCreate(ctx context.Context, client permissions.AccessorClient, objectID string, user *permissions.UserAccessor) (func(adminAccess bool) error, error) {
	body, err := json.Marshal(permissions.PermissionObject{
		Accessor: permissions.Accessor{
			Type: permissions.User,
			ID:   user.UID,
		},
		Permissions: convert.HCLToDTOPermission(user.Access),
	})

	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) error {
		_, err = client.Create(ctx, objectID, adminAccess, body)
		return err
	}, nil
}

func getUserUpdate(ctx context.Context, client permissions.AccessorClient, objectID string, user *permissions.UserAccessor) (func(adminAccess bool) error, error) {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(user.Access),
	})
	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) error {
		_, err = client.UpdateAccessor(ctx, objectID, permissions.User, user.UID, adminAccess, body)
		return err
	}, nil
}

func containsUser(users permissions.Users, user *permissions.UserAccessor) (exists, equals bool) {
	for _, u := range users {
		if u.UID == user.UID {
			return true, u.Access == user.Access
		}
	}
	return false, false
}
