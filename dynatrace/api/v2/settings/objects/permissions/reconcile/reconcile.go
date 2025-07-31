package reconcile

import (
	"context"
	"errors"

	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func CompareAndUpdate(ctx context.Context, cl permissions.PermissionUpdateClient, current *permissions.SettingPermissions, desired *permissions.SettingPermissions) error {
	groupErr := compareAndUpdateGroups(ctx, cl, current.SettingsObjectID, current.Groups, desired.Groups)
	userErr := compareAndUpdateUsers(ctx, cl, current.SettingsObjectID, current.Users, desired.Users)
	allUserEr := compareAndUpdateAllUsers(ctx, cl, current.SettingsObjectID, current.AllUsers, desired.AllUsers)

	return errors.Join(groupErr, userErr, allUserEr)
}
