package reconcile

import (
	"context"
	"errors"
	"slices"

	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

func CompareAndUpdate(ctx context.Context, cl permissions.PermissionUpdateClient, current *permissions.SettingPermissions, desired *permissions.SettingPermissions, adminAccess bool) error {
	groupUpserts, groupErr := getGroupUpserts(ctx, cl, current.SettingsObjectID, current.Groups, desired.Groups)
	userUpserts, userErr := getUserUpserts(ctx, cl, current.SettingsObjectID, current.Users, desired.Users)
	allUserUpsert, allUserErr := getAllUserUpsert(ctx, cl, current.SettingsObjectID, current.AllUsers, desired.AllUsers)

	prepareErrs := errors.Join(groupErr, userErr, allUserErr)

	if prepareErrs != nil {
		return prepareErrs
	}

	upserts := slices.Concat(groupUpserts, userUpserts)
	if allUserUpsert != nil {
		upserts = append(upserts, allUserUpsert)
	}
	return executeUpserts(upserts, adminAccess)
}

func executeUpserts(upserts []rest.AdminAccessRequestFn, adminAccess bool) error {
	errs := make([]error, 0, len(upserts))

	// the update must be in sync because of potential conflict errors (e.g., deleting two user permissions at the same may lead to a 409)
	for _, upsert := range upserts {
		_, err := upsert(adminAccess)
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
