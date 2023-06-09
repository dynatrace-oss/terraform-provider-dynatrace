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

package maintenancewindow

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GeneralProperties struct {
	Description                      *string         `json:"description,omitempty"`            // A short description of the maintenance purpose.
	DisableSyntheticMonitorExecution bool            `json:"disableSyntheticMonitorExecution"` // Disables the execution of the synthetic monitors that are within [the scope of this maintenance window](https://dt-url.net/0e0341m).
	MaintenanceType                  MaintenanceType `json:"maintenanceType"`                  // The type of the maintenance, possible values: `PLANNED` or `UNPLANNED`
	Name                             string          `json:"name"`                             // The name of the maintenance window, displayed in the UI
	Suppression                      SuppressionType `json:"suppression"`                      // The type of suppression of alerting and problem detection during the maintenance. Possible Values: `DETECT_PROBLEMS_AND_ALERT`, `DETECT_PROBLEMS_DONT_ALERT`, `DONT_DETECT_PROBLEMS`
}

func (me *GeneralProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the maintenance purpose.",
			Optional:    true, // nullable
		},
		"disable_synthetic": {
			Type:        schema.TypeBool,
			Description: "Disables the execution of the synthetic monitors that are within [the scope of this maintenance window](https://dt-url.net/0e0341m).",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the maintenance, possible values: `PLANNED` or `UNPLANNED`",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the maintenance window, displayed in the UI",
			Required:    true,
		},
		"suppression": {
			Type:        schema.TypeString,
			Description: "The type of suppression of alerting and problem detection during the maintenance. Possible Values: `DETECT_PROBLEMS_AND_ALERT`, `DETECT_PROBLEMS_DONT_ALERT`, `DONT_DETECT_PROBLEMS`",
			Required:    true,
		},
	}
}

func (me *GeneralProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":       me.Description,
		"disable_synthetic": me.DisableSyntheticMonitorExecution,
		"type":              me.MaintenanceType,
		"name":              me.Name,
		"suppression":       me.Suppression,
	})
}

func (me *GeneralProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description":       &me.Description,
		"disable_synthetic": &me.DisableSyntheticMonitorExecution,
		"type":              &me.MaintenanceType,
		"name":              &me.Name,
		"suppression":       &me.Suppression,
	})
}
