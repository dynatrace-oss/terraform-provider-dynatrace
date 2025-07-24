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

package settings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type SettingPermissions struct {
	SettingsObjectID string
	AllUsers         string
	Users            Users
	Groups           Groups
}

func (_ *SettingPermissions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"settings_object_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"all_users": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"read", "write", "none"}, false),
		},
		"users": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem:     &schema.Resource{Schema: new(Users).Schema()},
		},
		"groups": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem:     &schema.Resource{Schema: new(Groups).Schema()},
		},
	}
}

func (sp *SettingPermissions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"settings_object_id": sp.SettingsObjectID,
		"all_users":          sp.AllUsers,
		"users":              sp.Users,
		"groups":             sp.Groups,
	})
}

func (sp *SettingPermissions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"settings_object_id": &sp.SettingsObjectID,
		"all_users":          &sp.AllUsers,
		"users":              &sp.Users,
		"groups":             &sp.Groups,
	})
}
