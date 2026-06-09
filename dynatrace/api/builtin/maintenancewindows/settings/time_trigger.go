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

// Trigger. Trigger definition.
type TimeTrigger struct {
	EarliestStart *string `json:"earliestStart,omitempty"` // Earliest start date for when the first instance of this maintenance window should be created.
	Rule          *string `json:"rule,omitempty"`          // Reference to rule which specifies on which days instance of this maintenance window should be created.
	Time          string  `json:"time"`                    // Time for the trigger.
	Until         *string `json:"until,omitempty"`         // Date after which instances of this recurring maintenance window should no longer be created.
}

func (me *TimeTrigger) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"earliest_start": {
			Type:        schema.TypeString,
			Description: "Earliest start date for when the first instance of this maintenance window should be created.",
			Optional:    true, // nullable
		},
		"rule": {
			Type:        schema.TypeString,
			Description: "Reference to rule which specifies on which days instance of this maintenance window should be created.",
			Optional:    true, // nullable
		},
		"time": {
			Type:        schema.TypeString,
			Description: "Time for the trigger.",
			Required:    true,
		},
		"until": {
			Type:        schema.TypeString,
			Description: "Date after which instances of this recurring maintenance window should no longer be created.",
			Optional:    true, // nullable
		},
	}
}

func (me *TimeTrigger) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"earliest_start": me.EarliestStart,
		"rule":           me.Rule,
		"time":           me.Time,
		"until":          me.Until,
	})
}

func (me *TimeTrigger) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"earliest_start": &me.EarliestStart,
		"rule":           &me.Rule,
		"time":           &me.Time,
		"until":          &me.Until,
	})
}
