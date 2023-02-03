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

package rummobile

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SlowUserActionsAutoSlowest struct {
	DurationThreshold  float64 `json:"durationThreshold"`  // Absolute threshold
	SlowdownPercentage float64 `json:"slowdownPercentage"` // Relative threshold
}

func (me *SlowUserActionsAutoSlowest) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"duration_threshold": {
			Type:        schema.TypeFloat,
			Description: "Absolute threshold",
			Required:    true,
		},
		"slowdown_percentage": {
			Type:        schema.TypeFloat,
			Description: "Relative threshold",
			Required:    true,
		},
	}
}

func (me *SlowUserActionsAutoSlowest) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"duration_threshold":  me.DurationThreshold,
		"slowdown_percentage": me.SlowdownPercentage,
	})
}

func (me *SlowUserActionsAutoSlowest) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"duration_threshold":  &me.DurationThreshold,
		"slowdown_percentage": &me.SlowdownPercentage,
	})
}
