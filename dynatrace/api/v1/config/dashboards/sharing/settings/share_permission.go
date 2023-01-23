/**
* @license
* Copyright 2020 Dynatrace LLC
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

package sharing

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SharePermission represents access permissions of the dashboard
type SharePermission struct {
	ID         *string        `json:"id,omitempty"` // The ID of the user or group to whom the permission is granted.\n\nNot applicable if the **type** is set to `ALL`
	Type       PermissionType `json:"type"`         // The type of the permission: \n\n* `USER`: The dashboard is shared with the specified user. \n* `GROUP`: The dashboard is shared with all users of the specified group. \n* `ALL`: The dashboard is shared via link. Any authenticated user with the link can view the dashboard
	Permission Permission     `json:"permission"`   // The level of the permission: \n \n* `VIEW`: The dashboard is shared with read permission. \n* `EDIT`: The dashboard is shared with edit permission
}

func (me *SharePermission) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the user or group to whom the permission is granted.\n\nNot applicable if the **type** is set to `ALL`",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the permission: \n\n* `USER`: The dashboard is shared with the specified user. \n* `GROUP`: The dashboard is shared with all users of the specified group. \n* `ALL`: The dashboard is shared via link. Any authenticated user with the link can view the dashboard",
		},
		"level": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The level of the permission: \n \n* `VIEW`: The dashboard is shared with read permission. \n* `EDIT`: The dashboard is shared with edit permission",
		},
	}
}

// MarshalHCL has no documentation
func (me *SharePermission) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("id", me.ID); err != nil {
		return err
	}
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("level", string(me.Permission)); err != nil {
		return err
	}
	return nil
}

// UnmarshalHCL has no documentation
func (me *SharePermission) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("id"); ok {
		me.ID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = PermissionType(value.(string))
	}
	if value, ok := decoder.GetOk("level"); ok {
		me.Permission = Permission(value.(string))
	}
	return nil
}

type SharePermissions []*SharePermission

func (me *SharePermissions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"permission": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(SharePermission).Schema()},
			Description: "Access permissions of the dashboard",
		},
	}
}

// MarshalHCL has no documentation
func (me SharePermissions) MarshalHCL(properties hcl.Properties) error {
	props := hcl.Properties{}
	return props.EncodeSlice("permission", me)
}

// UnmarshalHCL has no documentation
func (me *SharePermissions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("permission", me)
}

type PermissionType string

var PermissionTypes = struct {
	All   PermissionType
	Group PermissionType
	User  PermissionType
}{
	All:   PermissionType("ALL"),
	Group: PermissionType("GROUP"),
	User:  PermissionType("USER"),
}

type Permission string

var Permissions = struct {
	Edit Permission
	View Permission
}{
	Edit: Permission("EDIT"),
	View: Permission("VIEW"),
}
