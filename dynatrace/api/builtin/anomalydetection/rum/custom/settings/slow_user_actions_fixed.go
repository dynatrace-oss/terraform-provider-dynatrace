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

package custom

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SlowUserActionsFixed struct {
	DurationAvoidOveralerting *SlowUserActionsAvoidOveralerting `json:"durationAvoidOveralerting"` // To avoid over-alerting do not alert for low traffic applications with less than
	DurationThresholdAllFixed *SlowUserActionsManualAll         `json:"durationThresholdAllFixed"` // Alert if the action duration of all user actions degrades beyond the absolute threshold:
	DurationThresholdSlowest  *SlowUserActionsManualSlowest     `json:"durationThresholdSlowest"`  // Alert if the action duration of the slowest 10% of user actions degrades beyond the absolute threshold:
	Sensitivity               Sensitivity                       `json:"sensitivity"`               // Possible Values: `Medium`, `High`, `Low`
}

func (me *SlowUserActionsFixed) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"duration_avoid_overalerting": {
			Type:        schema.TypeList,
			Description: "To avoid over-alerting do not alert for low traffic applications with less than",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SlowUserActionsAvoidOveralerting).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"duration_threshold_all_fixed": {
			Type:        schema.TypeList,
			Description: "Alert if the action duration of all user actions degrades beyond the absolute threshold:",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SlowUserActionsManualAll).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"duration_threshold_slowest": {
			Type:        schema.TypeList,
			Description: "Alert if the action duration of the slowest 10% of user actions degrades beyond the absolute threshold:",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SlowUserActionsManualSlowest).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"sensitivity": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Medium`, `High`, `Low`",
			Required:    true,
		},
	}
}

func (me *SlowUserActionsFixed) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"duration_avoid_overalerting":  me.DurationAvoidOveralerting,
		"duration_threshold_all_fixed": me.DurationThresholdAllFixed,
		"duration_threshold_slowest":   me.DurationThresholdSlowest,
		"sensitivity":                  me.Sensitivity,
	})
}

func (me *SlowUserActionsFixed) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"duration_avoid_overalerting":  &me.DurationAvoidOveralerting,
		"duration_threshold_all_fixed": &me.DurationThresholdAllFixed,
		"duration_threshold_slowest":   &me.DurationThresholdSlowest,
		"sensitivity":                  &me.Sensitivity,
	})
}
