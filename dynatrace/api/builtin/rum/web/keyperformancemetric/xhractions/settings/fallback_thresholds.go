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

package xhractions

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FallbackThresholds struct {
	FrustratingFallbackThresholdSeconds float64 `json:"frustratingFallbackThresholdSeconds"` // If **User action duration** is above this value, the action is assigned to the Frustrated performance zone.
	ToleratedFallbackThresholdSeconds   float64 `json:"toleratedFallbackThresholdSeconds"`   // If **User action duration** is below this value, the action is assigned to the Satisfied performance zone.
}

func (me *FallbackThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"frustrating_fallback_threshold_seconds": {
			Type:        schema.TypeFloat,
			Description: "If **User action duration** is above this value, the action is assigned to the Frustrated performance zone.",
			Required:    true,
		},
		"tolerated_fallback_threshold_seconds": {
			Type:        schema.TypeFloat,
			Description: "If **User action duration** is below this value, the action is assigned to the Satisfied performance zone.",
			Required:    true,
		},
	}
}

func (me *FallbackThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"frustrating_fallback_threshold_seconds": me.FrustratingFallbackThresholdSeconds,
		"tolerated_fallback_threshold_seconds":   me.ToleratedFallbackThresholdSeconds,
	})
}

func (me *FallbackThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"frustrating_fallback_threshold_seconds": &me.FrustratingFallbackThresholdSeconds,
		"tolerated_fallback_threshold_seconds":   &me.ToleratedFallbackThresholdSeconds,
	})
}
