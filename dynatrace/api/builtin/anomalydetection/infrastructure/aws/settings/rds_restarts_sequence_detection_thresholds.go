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

package aws

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RdsRestartsSequenceDetectionThresholds. Alert if the condition is met in 2 out of 20 samples
type RdsRestartsSequenceDetectionThresholds struct {
	RestartsPerMinute int `json:"restartsPerMinute"` // Number of restarts per minute is equal or higher than
}

func (me *RdsRestartsSequenceDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"restarts_per_minute": {
			Type:        schema.TypeInt,
			Description: "Number of restarts per minute is equal or higher than",
			Required:    true,
		},
	}
}

func (me *RdsRestartsSequenceDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"restarts_per_minute": me.RestartsPerMinute,
	})
}

func (me *RdsRestartsSequenceDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"restarts_per_minute": &me.RestartsPerMinute,
	})
}
