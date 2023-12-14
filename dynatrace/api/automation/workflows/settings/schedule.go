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

package workflows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Schedule struct {
	Active           bool                      `json:"isActive"`                     // The trigger is enabled (`true`) or not (`false`). Default is `false`
	Trigger          *ScheduleTrigger          `json:"trigger"`                      // Detailed configuration about the timing constraints that trigger the execution
	Rule             *string                   `json:"rule,omitempty" format:"uuid"` // Refers to a configured rule that determines at which days the schedule should be active. If not specified it implies that the schedule is valid every day
	FilterParameters *ScheduleFilterParameters `json:"filterParameters"`             // Advanced restrictions for the schedule to trigger executions
	Timezone         *string                   `json:"timezone,omitempty"`           // A time zone the scheduled times to align with. If not specified it will be chosen automatically based on the location of the Dynatrace Server

	// Inputs        *ScheduleInputs `json:"inputs"`                           //
	// Faulty        bool            `json:"isFaulty" flags:"readonly"`        // Signals that the configuration of this trigger is faulty
	// NextExecution string          `json:"nextExecution" format:"date-time"` //
}

func (me *Schedule) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"active": {
			Type:        schema.TypeBool,
			Description: "The trigger is enabled (`true`) or not (`false`). Default is `false`",
			Default:     false,
			Optional:    true,
		},
		"rule": {
			Type:             schema.TypeString,
			Description:      "Refers to a configured rule that determines at which days the schedule should be active. If not specified it implies that the schedule is valid every day",
			ValidateDiagFunc: ValidateUUID,
			Optional:         true,
		},
		"trigger": {
			Type:        schema.TypeList,
			Description: "Detailed configuration about the timing constraints that trigger the execution",
			MinItems:    1,
			MaxItems:    1,
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ScheduleTrigger).Schema(prefix + ".0.trigger")},
		},
		"filter_parameters": {
			Type:        schema.TypeList,
			Description: "Advanced restrictions for the schedule to trigger executions",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(ScheduleFilterParameters).Schema(prefix + ".0.filter_parameters")},
		},
		"time_zone": {
			Type:        schema.TypeString,
			Description: "A time zone the scheduled times to align with. If not specified it will be chosen automatically based on the location of the Dynatrace Server",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Schedule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"active":            me.Active,
		"rule":              me.Rule,
		"trigger":           me.Trigger,
		"filter_parameters": me.FilterParameters,
		"time_zone":         me.Timezone,
	})
}

func (me *Schedule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"active":            &me.Active,
		"rule":              &me.Rule,
		"trigger":           &me.Trigger,
		"filter_parameters": &me.FilterParameters,
		"time_zone":         &me.Timezone,
	})
}
