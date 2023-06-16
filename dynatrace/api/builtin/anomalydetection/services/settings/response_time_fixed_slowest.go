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

type ResponseTimeFixedSlowest struct {
	SlowestDegradationMilliseconds float64 `json:"slowestDegradationMilliseconds"` // Alert if the response time of the slowest 10% degrades beyond this many ms within an observation period of 5 minutes
}

func (me *ResponseTimeFixedSlowest) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"slowest_degradation_milliseconds": {
			Type:        schema.TypeFloat,
			Description: "Alert if the response time of the slowest 10% degrades beyond this many ms within an observation period of 5 minutes",
			Required:    true,
		},
	}
}

func (me *ResponseTimeFixedSlowest) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"slowest_degradation_milliseconds": me.SlowestDegradationMilliseconds,
	})
}

func (me *ResponseTimeFixedSlowest) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"slowest_degradation_milliseconds": &me.SlowestDegradationMilliseconds,
	})
}
