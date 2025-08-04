package reconcile

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func getGroupUpserts(ctx context.Context, client permissions.AccessorClient, objectID string, currentGroups, desiredGroups permissions.Groups) ([]func(adminAccess bool) error, error) {
	errs := make([]error, 0)
	updates := make([]func(adminAccess bool) error, 0)
	// update and create groups
	for _, group := range desiredGroups {
		var err error
		var fn func(adminAccess bool) error
		if exists, isEqual := containsGroups(currentGroups, group); !exists {
			fn, err = getGroupCreate(ctx, client, objectID, group)
		} else if !isEqual {
			fn, err = getGroupUpdate(ctx, client, objectID, group)
		} else {
			// group already exists and is equal, no action needed
			continue
		}

		if err != nil {
			errs = append(errs, err)
			continue
		}
		updates = append(updates, fn)
	}

	// delete groups that are not in desiredGroups (HCL)
	for _, group := range currentGroups {
		if exists, _ := containsGroups(desiredGroups, group); !exists {
			updates = append(updates, func(adminAccess bool) error {
				_, err := client.DeleteAccessor(ctx, objectID, permissions.Group, group.ID, adminAccess)
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

func getGroupCreate(ctx context.Context, client permissions.AccessorClient, objectID string, group *permissions.GroupAccessor) (func(adminAccess bool) error, error) {
	body, err := json.Marshal(permissions.PermissionObject{
		Accessor: permissions.Accessor{
			Type: permissions.Group,
			ID:   group.ID,
		},
		Permissions: convert.HCLToDTOPermission(group.Access),
	})

	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) error {
		_, err = client.Create(ctx, objectID, adminAccess, body)
		return err
	}, nil
}

func getGroupUpdate(ctx context.Context, client permissions.AccessorClient, objectID string, group *permissions.GroupAccessor) (func(adminAccess bool) error, error) {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(group.Access),
	})
	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) error {
		_, err = client.UpdateAccessor(ctx, objectID, permissions.Group, group.ID, adminAccess, body)
		return err
	}, nil
}

func containsGroups(groups permissions.Groups, group *permissions.GroupAccessor) (exists, equals bool) {
	for _, g := range groups {
		if g.ID == group.ID {
			return true, g.Access == group.Access
		}
	}
	return false, false
}
