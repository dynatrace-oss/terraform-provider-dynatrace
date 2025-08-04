package reconcile

import (
	"context"
	"errors"
	"slices"
	"sync"

	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func CompareAndUpdate(ctx context.Context, cl permissions.PermissionUpdateClient, current *permissions.SettingPermissions, desired *permissions.SettingPermissions) error {
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
	return executeUpserts(upserts)
}

func executeUpserts(upserts []func(adminAccess bool) error) error {
	errChan := make(chan error, len(upserts))

	var wg sync.WaitGroup
	wg.Add(len(upserts))
	for _, upsert := range upserts {
		go func() {
			defer wg.Done()
			errChan <- upsert(false)
		}()
	}
	wg.Wait()
	close(errChan)

	errs := make([]error, len(upserts))
	for err := range errChan {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
