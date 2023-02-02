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

package custom

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ErrorRateIncreaseFixed struct {
	Sensitivity       Sensitivity `json:"sensitivity"`       // Possible Values: `Low`, `Medium`, `High`
	ThresholdAbsolute float64     `json:"thresholdAbsolute"` // Absolute threshold
}

func (me *ErrorRateIncreaseFixed) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sensitivity": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Low`, `Medium`, `High`",
			Required:    true,
		},
		"threshold_absolute": {
			Type:        schema.TypeFloat,
			Description: "Absolute threshold",
			Required:    true,
		},
	}
}

func (me *ErrorRateIncreaseFixed) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"sensitivity":        me.Sensitivity,
		"threshold_absolute": me.ThresholdAbsolute,
	})
}

func (me *ErrorRateIncreaseFixed) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"sensitivity":        &me.Sensitivity,
		"threshold_absolute": &me.ThresholdAbsolute,
	})
}
