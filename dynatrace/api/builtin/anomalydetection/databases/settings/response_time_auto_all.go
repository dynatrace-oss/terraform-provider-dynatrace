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

package databases

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResponseTimeAutoAll struct {
	DegradationMilliseconds float64 `json:"degradationMilliseconds"` // Absolute threshold
	DegradationPercent      float64 `json:"degradationPercent"`      // Relative threshold
}

func (me *ResponseTimeAutoAll) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"degradation_milliseconds": {
			Type:        schema.TypeFloat,
			Description: "Absolute threshold",
			Required:    true,
		},
		"degradation_percent": {
			Type:        schema.TypeFloat,
			Description: "Relative threshold",
			Required:    true,
		},
	}
}

func (me *ResponseTimeAutoAll) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"degradation_milliseconds": me.DegradationMilliseconds,
		"degradation_percent":      me.DegradationPercent,
	})
}

func (me *ResponseTimeAutoAll) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"degradation_milliseconds": &me.DegradationMilliseconds,
		"degradation_percent":      &me.DegradationPercent,
	})
}
