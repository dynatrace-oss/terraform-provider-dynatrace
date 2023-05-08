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

package loadactions

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Thresholds struct {
	FrustratingThresholdSeconds float64 `json:"frustratingThresholdSeconds"` // If the key performance metric is above this value, the action is assigned to the Frustrated performance zone.
	ToleratedThresholdSeconds   float64 `json:"toleratedThresholdSeconds"`   // If the key performance metric is below this value, the action is assigned to the Satisfied performance zone.
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"frustrating_threshold_seconds": {
			Type:        schema.TypeFloat,
			Description: "If the key performance metric is above this value, the action is assigned to the Frustrated performance zone.",
			Required:    true,
		},
		"tolerated_threshold_seconds": {
			Type:        schema.TypeFloat,
			Description: "If the key performance metric is below this value, the action is assigned to the Satisfied performance zone.",
			Required:    true,
		},
	}
}

func (me *Thresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"frustrating_threshold_seconds": me.FrustratingThresholdSeconds,
		"tolerated_threshold_seconds":   me.ToleratedThresholdSeconds,
	})
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"frustrating_threshold_seconds": &me.FrustratingThresholdSeconds,
		"tolerated_threshold_seconds":   &me.ToleratedThresholdSeconds,
	})
}
