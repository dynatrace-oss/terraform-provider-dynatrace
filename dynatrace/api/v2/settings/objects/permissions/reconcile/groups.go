/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package reconcile

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

func getGroupUpserts(ctx context.Context, client permissions.AccessorClient, objectID string, currentGroups, desiredGroups permissions.Groups) ([]rest.AdminAccessRequestFn, error) {
	errs := make([]error, 0)
	updates := make([]rest.AdminAccessRequestFn, 0)
	// update and create groups
	for _, group := range desiredGroups {
		var err error
		var fn rest.AdminAccessRequestFn
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
			updates = append(updates, func(adminAccess bool) (api.Response, error) {
				return client.DeleteAccessor(ctx, objectID, permissions.Group, group.ID, adminAccess)
			})
		}
	}
	joinErr := errors.Join(errs...)
	if joinErr != nil {
		return nil, joinErr
	}
	return updates, nil
}

func getGroupCreate(ctx context.Context, client permissions.AccessorClient, objectID string, group *permissions.GroupAccessor) (rest.AdminAccessRequestFn, error) {
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

	return func(adminAccess bool) (api.Response, error) {
		return client.Create(ctx, objectID, adminAccess, body)
	}, nil
}

func getGroupUpdate(ctx context.Context, client permissions.AccessorClient, objectID string, group *permissions.GroupAccessor) (rest.AdminAccessRequestFn, error) {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(group.Access),
	})
	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) (api.Response, error) {
		return client.UpdateAccessor(ctx, objectID, permissions.Group, group.ID, adminAccess, body)
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
