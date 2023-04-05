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

package updates

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MaintenanceWindows []*MaintenanceWindow

func (me *MaintenanceWindows) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"maintenance_window": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(MaintenanceWindow).Schema()},
		},
	}
}

func (me MaintenanceWindows) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("maintenance_window", me)
}

func (me *MaintenanceWindows) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("maintenance_window", me)
}

type MaintenanceWindow struct {
	MaintenanceWindow string `json:"maintenanceWindow"` // Select a [maintenance window for OneAgent updates](/ui/settings/builtin:deployment.management.update-windows)
}

func (me *MaintenanceWindow) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"maintenance_window": {
			Type:        schema.TypeString,
			Description: "Select a [maintenance window for OneAgent updates](/ui/settings/builtin:deployment.management.update-windows)",
			Required:    true,
		},
	}
}

func (me *MaintenanceWindow) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"maintenance_window": me.MaintenanceWindow,
	})
}

func (me *MaintenanceWindow) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"maintenance_window": &me.MaintenanceWindow,
	})
}
