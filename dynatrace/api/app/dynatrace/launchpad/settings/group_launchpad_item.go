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

package launchpad

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GroupLaunchpadItems []*GroupLaunchpadItem

func (me *GroupLaunchpadItems) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_launchpad": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(GroupLaunchpadItem).Schema()},
		},
	}
}

func (me GroupLaunchpadItems) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("group_launchpad", me)
}

func (me *GroupLaunchpadItems) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("group_launchpad", me)
}

type GroupLaunchpadItem struct {
	IsEnabled   bool   `json:"isEnabled"`   // State
	LaunchpadID string `json:"launchpadId"` // Launchpad
	UserGroupID string `json:"userGroupId"` // User Group
}

func (me *GroupLaunchpadItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"is_enabled": {
			Type:        schema.TypeBool,
			Description: "State",
			Required:    true,
		},
		"launchpad_id": {
			Type:        schema.TypeString,
			Description: "Launchpad",
			Required:    true,
		},
		"user_group_id": {
			Type:        schema.TypeString,
			Description: "User Group",
			Required:    true,
		},
	}
}

func (me *GroupLaunchpadItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"is_enabled":    me.IsEnabled,
		"launchpad_id":  me.LaunchpadID,
		"user_group_id": me.UserGroupID,
	})
}

func (me *GroupLaunchpadItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"is_enabled":    &me.IsEnabled,
		"launchpad_id":  &me.LaunchpadID,
		"user_group_id": &me.UserGroupID,
	})
}
