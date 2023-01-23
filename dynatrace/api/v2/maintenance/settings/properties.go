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

package maintenance

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GeneralProperties struct {
	Name             string      `json:"name"`                             // The name of the maintenance window, displayed in the UI.
	Description      *string     `json:"description,omitempty"`            // A short description of the maintenance purpose.
	Type             WindowType  `json:"maintenanceType"`                  // The type of the maintenance: planned or unplanned
	Suppression      Suppression `json:"suppression"`                      // The type of suppression of alerting and problem detection during the maintenance.
	DisableSynthetic bool        `json:"disableSyntheticMonitorExecution"` // Suppress execution of synthetic monitors during the maintenance
}

func (me *GeneralProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the maintenance window, displayed in the UI",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the maintenance purpose",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the maintenance: planned or unplanned",
			Required:    true,
		},
		"suppression": {
			Type:        schema.TypeString,
			Description: "The type of suppression of alerting and problem detection during the maintenance",
			Required:    true,
		},
		"disable_synthetic": {
			Type:        schema.TypeBool,
			Description: "Suppress execution of synthetic monitors during the maintenance",
			Default:     false,
			Optional:    true,
		},
	}
}

func (me *GeneralProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":              me.Name,
		"description":       me.Description,
		"type":              me.Type,
		"suppression":       me.Suppression,
		"disable_synthetic": me.DisableSynthetic,
	})
}

func (me *GeneralProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":              &me.Name,
		"description":       &me.Description,
		"type":              &me.Type,
		"suppression":       &me.Suppression,
		"disable_synthetic": &me.DisableSynthetic,
	})
}

// WindowType The type of the maintenance: planned or unplanned.
type WindowType string

// MaintenanceWindowTypes offers the known enum values
var MaintenanceWindowTypes = struct {
	Planned   WindowType
	Unplanned WindowType
}{
	"PLANNED",
	"UNPLANNED",
}

// Suppression The type of suppression of alerting and problem detection during the maintenance.
type Suppression string

// Suppressions offers the known enum values
var Suppressions = struct {
	DetectProblemsAndAlert  Suppression
	DetectProblemsDontAlert Suppression
	DontDetectProblems      Suppression
}{
	"DETECT_PROBLEMS_AND_ALERT",
	"DETECT_PROBLEMS_DONT_ALERT",
	"DONT_DETECT_PROBLEMS",
}
