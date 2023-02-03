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

type UnexpectedLowLoad struct {
	Enabled             bool     `json:"enabled"`                       // This setting is enabled (`true`) or disabled (`false`)
	ThresholdPercentage *float64 `json:"thresholdPercentage,omitempty"` // Dynatrace learns your typical application traffic over an observation period of one week. Depending on this expected value Dynatrace detects abnormal traffic drops within your application.
}

func (me *UnexpectedLowLoad) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"threshold_percentage": {
			Type:        schema.TypeFloat,
			Description: "Dynatrace learns your typical application traffic over an observation period of one week. Depending on this expected value Dynatrace detects abnormal traffic drops within your application.",
			Optional:    true,
		},
	}
}

func (me *UnexpectedLowLoad) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":              me.Enabled,
		"threshold_percentage": me.ThresholdPercentage,
	})
}

func (me *UnexpectedLowLoad) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":              &me.Enabled,
		"threshold_percentage": &me.ThresholdPercentage,
	})
}
