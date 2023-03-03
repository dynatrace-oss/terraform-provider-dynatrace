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

type HighSystemLoadDetectionThresholds struct {
	EventThresholds *EventThresholds `json:"eventThresholds"`
	SystemLoad      float64          `json:"systemLoad"` // Alert if the System Load / Logical cpu core is higher than this threshold for the defined amount of samples
}

func (me *HighSystemLoadDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"system_load": {
			Type:        schema.TypeFloat,
			Description: "Alert if the System Load / Logical cpu core is higher than this threshold for the defined amount of samples",
			Required:    true,
		},
	}
}

func (me *HighSystemLoadDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_thresholds": me.EventThresholds,
		"system_load":      me.SystemLoad,
	})
}

func (me *HighSystemLoadDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_thresholds": &me.EventThresholds,
		"system_load":      &me.SystemLoad,
	})
}
