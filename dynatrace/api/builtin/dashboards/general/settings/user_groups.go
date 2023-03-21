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

package general

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type UserGroupss []*UserGroups

func (me *UserGroupss) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_dashboard": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(UserGroups).Schema()},
		},
	}
}

func (me UserGroupss) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("default_dashboard", me)
}

func (me *UserGroupss) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("default_dashboard", me)
}

type UserGroups struct {
	Dashboard string `json:"Dashboard"` // Preset dashboard to show as default landing page
	UserGroup string `json:"UserGroup"` // Show selected dashboard by default for this user group
}

func (me *UserGroups) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dashboard": {
			Type:        schema.TypeString,
			Description: "Preset dashboard to show as default landing page",
			Required:    true,
		},
		"user_group": {
			Type:        schema.TypeString,
			Description: "Show selected dashboard by default for this user group",
			Required:    true,
		},
	}
}

func (me *UserGroups) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dashboard":  me.Dashboard,
		"user_group": me.UserGroup,
	})
}

func (me *UserGroups) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dashboard":  &me.Dashboard,
		"user_group": &me.UserGroup,
	})
}
