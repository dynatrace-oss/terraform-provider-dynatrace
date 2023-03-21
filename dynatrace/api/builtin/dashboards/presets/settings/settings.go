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

type Settings struct {
	DashboardPresetsList   DashboardPresetsList `json:"dashboardPresetsList"`   // Show selected preset to respective user group only.
	EnableDashboardPresets bool                 `json:"enableDashboardPresets"` // Dashboard presets are visible to all users by default. For a pristine environment you may disable them entirely or opt to manually limit visibility to selected user groups.
}

func (me *Settings) Name() string {
	return "dashboards_preset"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dashboard_presets_list": {
			Type:        schema.TypeList,
			Description: "Show selected preset to respective user group only.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(DashboardPresetsList).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"enable_dashboard_presets": {
			Type:        schema.TypeBool,
			Description: "Dashboard presets are visible to all users by default. For a pristine environment you may disable them entirely or opt to manually limit visibility to selected user groups.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dashboard_presets_list":   me.DashboardPresetsList,
		"enable_dashboard_presets": me.EnableDashboardPresets,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dashboard_presets_list":   &me.DashboardPresetsList,
		"enable_dashboard_presets": &me.EnableDashboardPresets,
	})
}
