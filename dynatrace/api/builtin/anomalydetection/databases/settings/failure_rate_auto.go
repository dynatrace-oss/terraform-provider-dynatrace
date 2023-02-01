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

type FailureRateAuto struct {
	AbsoluteIncrease       float64                 `json:"absoluteIncrease"`       // Absolute threshold
	OverAlertingProtection *OverAlertingProtection `json:"overAlertingProtection"` // Avoid over-alerting
	RelativeIncrease       float64                 `json:"relativeIncrease"`       // Relative threshold
}

func (me *FailureRateAuto) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"absolute_increase": {
			Type:        schema.TypeFloat,
			Description: "Absolute threshold",
			Required:    true,
		},
		"over_alerting_protection": {
			Type:        schema.TypeList,
			Description: "Avoid over-alerting",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(OverAlertingProtection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"relative_increase": {
			Type:        schema.TypeFloat,
			Description: "Relative threshold",
			Required:    true,
		},
	}
}

func (me *FailureRateAuto) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"absolute_increase":        me.AbsoluteIncrease,
		"over_alerting_protection": me.OverAlertingProtection,
		"relative_increase":        me.RelativeIncrease,
	})
}

func (me *FailureRateAuto) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"absolute_increase":        &me.AbsoluteIncrease,
		"over_alerting_protection": &me.OverAlertingProtection,
		"relative_increase":        &me.RelativeIncrease,
	})
}
