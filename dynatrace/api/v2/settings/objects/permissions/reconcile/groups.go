package reconcile

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func compareAndUpdateGroups(ctx context.Context, client permissions.AccessorClient, objectID string, currentGroups, desiredGroups permissions.Groups) error {
	errs := make([]error, 0)
	// update and create groups
	for _, group := range desiredGroups {
		var err error
		if exists, isEqual := containsGroups(currentGroups, group); !exists {
			err = createGroup(ctx, client, objectID, group)
		} else if !isEqual {
			err = updateGroup(ctx, client, objectID, group)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	// delete groups that are not in desiredGroups (HCL)
	for _, group := range currentGroups {
		if exists, _ := containsGroups(desiredGroups, group); !exists {
			_, err := client.DeleteAccessor(ctx, objectID, permissions.Group, group.ID, false)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func createGroup(ctx context.Context, client permissions.AccessorClient, objectID string, group *permissions.GroupAccessor) error {
	body, err := json.Marshal(permissions.PermissionObject{
		Accessor: permissions.Accessor{
			Type: permissions.Group,
			ID:   group.ID,
		},
		Permissions: convert.HCLToDTOPermission(group.Access),
	})

	if err != nil {
		return err
	}

	_, err = client.Create(ctx, objectID, false, body)
	return err
}

func updateGroup(ctx context.Context, client permissions.AccessorClient, objectID string, group *permissions.GroupAccessor) error {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(group.Access),
	})
	if err != nil {
		return err
	}

	_, err = client.UpdateAccessor(ctx, objectID, permissions.Group, group.ID, false, body)
	return err
}

func containsGroups(groups permissions.Groups, group *permissions.GroupAccessor) (exists, equals bool) {
	for _, u := range groups {
		if u.ID == group.ID {
			return true, u.Access == group.Access
		}
	}
	return false, false
}
