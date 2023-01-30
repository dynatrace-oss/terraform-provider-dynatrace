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

// No documentation available
type FailureRateAuto struct {
	OverAlertingProtection *OverAlertingProtection `json:"overAlertingProtection"` // Avoid over-alerting
	AbsoluteIncrease       float64                 `json:"absoluteIncrease"`       // Absolute threshold
	RelativeIncrease       float64                 `json:"relativeIncrease"`       // Relative threshold
}

func (me *FailureRateAuto) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"over_alerting_protection": {
			Type:        schema.TypeList,
			Description: "Avoid over-alerting",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(OverAlertingProtection).Schema()},
			Required:    true,
		},
		"absolute_increase": {
			Type:        schema.TypeFloat,
			Description: "Absolute threshold",
			Required:    true,
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
		"over_alerting_protection": me.OverAlertingProtection,
		"absolute_increase":        me.AbsoluteIncrease,
		"relative_increase":        me.RelativeIncrease,
	})
}

func (me *FailureRateAuto) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"over_alerting_protection": &me.OverAlertingProtection,
		"absolute_increase":        &me.AbsoluteIncrease,
		"relative_increase":        &me.RelativeIncrease,
	})
}
