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

package services

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResponseTimeAutoSlowest struct {
	SlowestDegradationMilliseconds float64 `json:"slowestDegradationMilliseconds"` // Absolute threshold
	SlowestDegradationPercent      float64 `json:"slowestDegradationPercent"`      // Relative threshold
}

func (me *ResponseTimeAutoSlowest) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"slowest_degradation_milliseconds": {
			Type:        schema.TypeFloat,
			Description: "Absolute threshold",
			Required:    true,
		},
		"slowest_degradation_percent": {
			Type:        schema.TypeFloat,
			Description: "Relative threshold",
			Required:    true,
		},
	}
}

func (me *ResponseTimeAutoSlowest) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"slowest_degradation_milliseconds": me.SlowestDegradationMilliseconds,
		"slowest_degradation_percent":      me.SlowestDegradationPercent,
	})
}

func (me *ResponseTimeAutoSlowest) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"slowest_degradation_milliseconds": &me.SlowestDegradationMilliseconds,
		"slowest_degradation_percent":      &me.SlowestDegradationPercent,
	})
}
