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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

func getAllUserUpsert(ctx context.Context, client permissions.AllUsersClient, objectID string, currentAllUsers, desiredAllUsers string) (rest.AdminAccessRequestFn, error) {
	if currentAllUsers == desiredAllUsers {
		// No update needed, current and desired are the same
		return nil, nil
	}

	if desiredAllUsers == permissions.HCLAccessorNone {
		return func(adminAccess bool) (api.Response, error) {
			return client.DeleteAllUsersAccessor(ctx, objectID, adminAccess)
		}, nil
	}

	if currentAllUsers == permissions.HCLAccessorNone {
		return getAllUserCreate(ctx, client, objectID, desiredAllUsers)
	}

	return getAllUserUpdate(ctx, client, objectID, desiredAllUsers)
}

func getAllUserCreate(ctx context.Context, client permissions.AllUsersClient, objectID string, desiredAllUsers string) (rest.AdminAccessRequestFn, error) {
	body, err := json.Marshal(permissions.PermissionObject{
		Accessor: permissions.Accessor{
			Type: permissions.AllUsers,
		},
		Permissions: convert.HCLToDTOPermission(desiredAllUsers),
	})
	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) (api.Response, error) {
		return client.Create(ctx, objectID, adminAccess, body)
	}, nil
}

func getAllUserUpdate(ctx context.Context, client permissions.AllUsersClient, objectID string, desiredAllUsers string) (rest.AdminAccessRequestFn, error) {
	body, err := json.Marshal(permissions.PermissionObjectUpdate{
		Permissions: convert.HCLToDTOPermission(desiredAllUsers),
	})
	if err != nil {
		return nil, err
	}

	return func(adminAccess bool) (api.Response, error) {
		return client.UpdateAllUsersAccessor(ctx, objectID, adminAccess, body)
	}, nil
}
