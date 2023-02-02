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

package rumcustomcrashrateincrease

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CrashRateIncreaseAuto struct {
	BaselineViolationPercentage float64     `json:"baselineViolationPercentage"` // Dynatrace learns the typical crash rate for all app versions and will create an alert if the baseline is violated by more than a specified threshold. Analysis happens based on a sliding window of 10 minutes.
	ConcurrentUsers             float64     `json:"concurrentUsers"`             // Amount of users
	Sensitivity                 Sensitivity `json:"sensitivity"`                 // Possible Values: `Low`, `Medium`, `High`
}

func (me *CrashRateIncreaseAuto) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"baseline_violation_percentage": {
			Type:        schema.TypeFloat,
			Description: "Dynatrace learns the typical crash rate for all app versions and will create an alert if the baseline is violated by more than a specified threshold. Analysis happens based on a sliding window of 10 minutes.",
			Required:    true,
		},
		"concurrent_users": {
			Type:        schema.TypeFloat,
			Description: "Amount of users",
			Required:    true,
		},
		"sensitivity": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Low`, `Medium`, `High`",
			Required:    true,
		},
	}
}

func (me *CrashRateIncreaseAuto) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"baseline_violation_percentage": me.BaselineViolationPercentage,
		"concurrent_users":              me.ConcurrentUsers,
		"sensitivity":                   me.Sensitivity,
	})
}

func (me *CrashRateIncreaseAuto) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"baseline_violation_percentage": &me.BaselineViolationPercentage,
		"concurrent_users":              &me.ConcurrentUsers,
		"sensitivity":                   &me.Sensitivity,
	})
}
