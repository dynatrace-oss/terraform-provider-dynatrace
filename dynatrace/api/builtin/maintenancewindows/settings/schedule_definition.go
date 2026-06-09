/**
* @license
* Copyright 2026 Dynatrace LLC
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

package maintenancewindows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ScheduleDefinition. Schedule definition for platform service schedules.
type ScheduleDefinition struct {
	Duration int      `json:"duration"`           // Duration of the maintenance window in minutes.
	Timezone *string  `json:"timezone,omitempty"` // Time zone
	Trigger  *Trigger `json:"trigger"`            // Trigger definition.
}

func (me *ScheduleDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"duration": {
			Type:        schema.TypeInt,
			Description: "Duration of the maintenance window in minutes.",
			Required:    true,
		},
		"timezone": {
			Type:        schema.TypeString,
			Description: "Time zone",
			Optional:    true, // nullable
		},
		"trigger": {
			Type:        schema.TypeList,
			Description: "Trigger definition.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Trigger).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *ScheduleDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"duration": me.Duration,
		"timezone": me.Timezone,
		"trigger":  me.Trigger,
	})
}

func (me *ScheduleDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"duration": &me.Duration,
		"timezone": &me.Timezone,
		"trigger":  &me.Trigger,
	})
}
