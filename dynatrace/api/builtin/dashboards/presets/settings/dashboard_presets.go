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

package presets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DashboardPresetsList []*DashboardPresets

func (me *DashboardPresetsList) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dashboard_presets": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DashboardPresets).Schema()},
		},
	}
}

func (me DashboardPresetsList) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("dashboard_presets", me)
}

func (me *DashboardPresetsList) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("dashboard_presets", me)
}

type DashboardPresets struct {
	DashboardPreset string `json:"DashboardPreset"` // Dashboard preset to limit visibility for
	UserGroup       string `json:"UserGroup"`       // User group to show selected dashboard preset to
}

func (me *DashboardPresets) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dashboard_preset": {
			Type:        schema.TypeString,
			Description: "Dashboard preset to limit visibility for",
			Required:    true,
		},
		"user_group": {
			Type:        schema.TypeString,
			Description: "User group to show selected dashboard preset to",
			Required:    true,
		},
	}
}

func (me *DashboardPresets) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dashboard_preset": me.DashboardPreset,
		"user_group":       me.UserGroup,
	})
}

func (me *DashboardPresets) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dashboard_preset": &me.DashboardPreset,
		"user_group":       &me.UserGroup,
	})
}
