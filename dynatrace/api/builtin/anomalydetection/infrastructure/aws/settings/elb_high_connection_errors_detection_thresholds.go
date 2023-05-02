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

// ElbHighConnectionErrorsDetectionThresholds. Alert if the condition is met in 3 out of 5 samples
type ElbHighConnectionErrorsDetectionThresholds struct {
	ConnectionErrorsPerMinute int `json:"connectionErrorsPerMinute"` // Number of backend connection errors is higher than
}

func (me *ElbHighConnectionErrorsDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_errors_per_minute": {
			Type:        schema.TypeInt,
			Description: "Number of backend connection errors is higher than",
			Required:    true,
		},
	}
}

func (me *ElbHighConnectionErrorsDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"connection_errors_per_minute": me.ConnectionErrorsPerMinute,
	})
}

func (me *ElbHighConnectionErrorsDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"connection_errors_per_minute": &me.ConnectionErrorsPerMinute,
	})
}
