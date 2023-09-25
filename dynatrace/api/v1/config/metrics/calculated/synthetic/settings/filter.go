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

package synthetic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Filter of a calculated synthetic metric.
type Filter struct {
	ActionType *string `json:"actionType,omitempty"` // Only user actions of the specified type are included in the metric calculation
	HasError   *bool   `json:"hasError,omitempty"`   // The execution status of the monitors to be included in the metric calculation: `true` or `false`
	ErrorCode  *int    `json:"errorCode,omitempty"`  // Only executions finished with the specified error code are included in the metric calculation.
	Event      *string `json:"event,omitempty"`      // Only the specified browser clickpath event is included in the metric calculation. Specify the Dynatrace entity ID of the event here.
	Location   *string `json:"location,omitempty"`   // Only executions from the specified location are included in the metric calculation. Specify the Dynatrace entity ID of the location here.
}

func (me *Filter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only user actions of the specified type are included in the metric calculation",
		},
		"has_error": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The execution status of the monitors to be included in the metric calculation: `true` or `false`",
		},
		"error_code": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Only executions finished with the specified error code are included in the metric calculation.",
		},
		"event": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only the specified browser clickpath event is included in the metric calculation. Specify the Dynatrace entity ID of the event here.",
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only executions from the specified location are included in the metric calculation. Specify the Dynatrace entity ID of the location here.",
		},
	}
}

func (me *Filter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"action_type": me.ActionType,
		"has_error":   me.HasError,
		"error_code":  me.ErrorCode,
		"event":       me.Event,
		"location":    me.Location,
	})
}

func (me *Filter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"action_type": &me.ActionType,
		"has_error":   &me.HasError,
		"error_code":  &me.ErrorCode,
		"event":       &me.Event,
		"location":    &me.Location,
	})
}
