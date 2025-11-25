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

//go:build unit

package convert_test

import (
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/convert"
	permissions2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions/settings"
	"github.com/stretchr/testify/assert"
)

func TestDTOToHCL(t *testing.T) {
	t.Run("Maps DTO read permissions to HCL format", func(t *testing.T) {
		dtoPermissions := permissions2.PermissionObjects{
			Accessors: []permissions2.PermissionObject{
				{
					Accessor: permissions2.Accessor{
						Type: permissions2.Group,
						ID:   "group-id",
					},
					Permissions: []permissions2.TypePermissions{permissions2.Read},
				},
				{
					Accessor: permissions2.Accessor{
						Type: permissions2.User,
						ID:   "user-id",
					},
					Permissions: []permissions2.TypePermissions{permissions2.Read},
				},
				{
					Accessor: permissions2.Accessor{
						Type: permissions2.AllUsers,
					},
					Permissions: []permissions2.TypePermissions{permissions2.Read},
				},
			},
		}

		hclPermissions := convert.DTOToHCL(dtoPermissions, "object-id")

		assert.Equal(t, &permissions2.SettingPermissions{
			SettingsObjectID: "object-id",
			AllUsers:         permissions2.HCLAccessorRead,
			Users: permissions2.Users{
				{
					UID:    "user-id",
					Access: permissions2.HCLAccessorRead,
				},
			},
			Groups: permissions2.Groups{
				{
					ID:     "group-id",
					Access: permissions2.HCLAccessorRead,
				},
			},
		}, hclPermissions)
	})

	t.Run("Maps DTO write permissions to HCL format", func(t *testing.T) {
		dtoPermissions := permissions2.PermissionObjects{
			Accessors: []permissions2.PermissionObject{
				{
					Accessor: permissions2.Accessor{
						Type: permissions2.Group,
						ID:   "group-id",
					},
					Permissions: []permissions2.TypePermissions{permissions2.Read, permissions2.Write},
				},
				{
					Accessor: permissions2.Accessor{
						Type: permissions2.User,
						ID:   "user-id",
					},
					Permissions: []permissions2.TypePermissions{permissions2.Read, permissions2.Write},
				},
				{
					Accessor: permissions2.Accessor{
						Type: permissions2.AllUsers,
					},
					Permissions: []permissions2.TypePermissions{permissions2.Read, permissions2.Write},
				},
			},
		}

		hclPermissions := convert.DTOToHCL(dtoPermissions, "object-id")

		assert.Equal(t, &permissions2.SettingPermissions{
			SettingsObjectID: "object-id",
			AllUsers:         permissions2.HCLAccessorWrite,
			Users: permissions2.Users{
				{
					UID:    "user-id",
					Access: permissions2.HCLAccessorWrite,
				},
			},
			Groups: permissions2.Groups{
				{
					ID:     "group-id",
					Access: permissions2.HCLAccessorWrite,
				},
			},
		}, hclPermissions)
	})

	t.Run("Maps empty DTO permissions to HCL format", func(t *testing.T) {
		dtoPermissions := permissions2.PermissionObjects{
			Accessors: []permissions2.PermissionObject{},
		}

		hclPermissions := convert.DTOToHCL(dtoPermissions, "object-id")

		assert.Equal(t, &permissions2.SettingPermissions{
			SettingsObjectID: "object-id",
			AllUsers:         permissions2.HCLAccessorNone,
			Users:            permissions2.Users{},
			Groups:           permissions2.Groups{},
		}, hclPermissions)
	})
}

func TestHCLToDTOPermission(t *testing.T) {
	t.Run("Maps DTO read permissions to HCL format", func(t *testing.T) {
		hclPermission := convert.HCLToDTOPermission(permissions2.HCLAccessorRead)

		assert.Equal(t, []permissions2.TypePermissions{permissions2.Read}, hclPermission)
	})

	t.Run("Maps DTO write permissions to HCL format", func(t *testing.T) {
		hclPermission := convert.HCLToDTOPermission(permissions2.HCLAccessorWrite)

		assert.Equal(t, []permissions2.TypePermissions{permissions2.Read, permissions2.Write}, hclPermission)
	})
}
