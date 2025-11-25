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

package convert

import (
	"slices"

	permissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
)

func DTOToHCL(permissionObjects permissions.PermissionObjects, objectID string) *permissions.SettingPermissions {
	groups := make(permissions.Groups, 0)
	users := make(permissions.Users, 0)
	allUsers := permissions.HCLAccessorNone

	for _, permissionObject := range permissionObjects.Accessors {
		access := convertDTOToHCLPermission(permissionObject.Permissions)
		switch permissionObject.Accessor.Type {
		case permissions.Group:
			groups = append(groups, &permissions.GroupAccessor{
				ID:     permissionObject.Accessor.ID,
				Access: access,
			})
		case permissions.User:
			users = append(users, &permissions.UserAccessor{
				UID:    permissionObject.Accessor.ID,
				Access: access,
			})
		case permissions.AllUsers:
			allUsers = access
		}
	}
	return &permissions.SettingPermissions{
		SettingsObjectID: objectID,
		Users:            users,
		Groups:           groups,
		AllUsers:         allUsers,
	}
}

// convertDTOToHCLPermission converts the permissions from the DTO format to the HCL format.
func convertDTOToHCLPermission(pms []permissions.TypePermissions) permissions.HCLAccessor {
	if slices.Contains(pms, permissions.Write) {
		return permissions.HCLAccessorWrite
	}
	return permissions.HCLAccessorRead
}

func HCLToDTOPermission(access permissions.HCLAccessor) []permissions.TypePermissions {
	if access == permissions.HCLAccessorWrite {
		return []permissions.TypePermissions{permissions.Read, permissions.Write}
	}
	return []permissions.TypePermissions{permissions.Read}
}
