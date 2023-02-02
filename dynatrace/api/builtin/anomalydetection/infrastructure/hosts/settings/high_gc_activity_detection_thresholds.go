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

package hosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HighGcActivityDetectionThresholds struct {
	EventThresholds        *EventThresholds `json:"eventThresholds"`
	GcSuspensionPercentage int              `json:"gcSuspensionPercentage"` // Alert if the GC suspension is higher than this threshold
	GcTimePercentage       int              `json:"gcTimePercentage"`       // Alert if GC time is higher than this threshold
}

func (me *HighGcActivityDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"gc_suspension_percentage": {
			Type:        schema.TypeInt,
			Description: "Alert if the GC suspension is higher than this threshold",
			Required:    true,
		},
		"gc_time_percentage": {
			Type:        schema.TypeInt,
			Description: "Alert if GC time is higher than this threshold",
			Required:    true,
		},
	}
}

func (me *HighGcActivityDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_thresholds":         me.EventThresholds,
		"gc_suspension_percentage": me.GcSuspensionPercentage,
		"gc_time_percentage":       me.GcTimePercentage,
	})
}

func (me *HighGcActivityDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_thresholds":         &me.EventThresholds,
		"gc_suspension_percentage": &me.GcSuspensionPercentage,
		"gc_time_percentage":       &me.GcTimePercentage,
	})
}
