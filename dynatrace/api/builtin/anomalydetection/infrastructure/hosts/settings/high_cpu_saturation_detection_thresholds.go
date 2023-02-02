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

type HighCpuSaturationDetectionThresholds struct {
	CpuSaturation   int              `json:"cpuSaturation"` // Alert if the CPU usage is higher than this threshold for the defined amount of samples
	EventThresholds *EventThresholds `json:"eventThresholds"`
}

func (me *HighCpuSaturationDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpu_saturation": {
			Type:        schema.TypeInt,
			Description: "Alert if the CPU usage is higher than this threshold for the defined amount of samples",
			Required:    true,
		},
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *HighCpuSaturationDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cpu_saturation":   me.CpuSaturation,
		"event_thresholds": me.EventThresholds,
	})
}

func (me *HighCpuSaturationDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cpu_saturation":   &me.CpuSaturation,
		"event_thresholds": &me.EventThresholds,
	})
}
