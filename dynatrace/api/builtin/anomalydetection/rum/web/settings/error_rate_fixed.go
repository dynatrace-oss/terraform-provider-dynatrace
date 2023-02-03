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

package rumweb

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ErrorRateFixed struct {
	ErrorRateReqPerMin     float64     `json:"errorRateReqPerMin"`     // To avoid over-alerting for low traffic applications
	ErrorRateSensitivity   Sensitivity `json:"errorRateSensitivity"`   // Possible Values: `Low`, `Medium`, `High`
	MaxFailureRateIncrease float64     `json:"maxFailureRateIncrease"` // Alert if this custom error rate threshold is exceeded during any 5-minute-period
	MinutesAbnormalState   float64     `json:"minutesAbnormalState"`   // Amount of minutes the observed traffic has to stay in abnormal state before alert
}

func (me *ErrorRateFixed) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_rate_req_per_min": {
			Type:        schema.TypeFloat,
			Description: "To avoid over-alerting for low traffic applications",
			Required:    true,
		},
		"error_rate_sensitivity": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Low`, `Medium`, `High`",
			Required:    true,
		},
		"max_failure_rate_increase": {
			Type:        schema.TypeFloat,
			Description: "Alert if this custom error rate threshold is exceeded during any 5-minute-period",
			Required:    true,
		},
		"minutes_abnormal_state": {
			Type:        schema.TypeFloat,
			Description: "Amount of minutes the observed traffic has to stay in abnormal state before alert",
			Required:    true,
		},
	}
}

func (me *ErrorRateFixed) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"error_rate_req_per_min":    me.ErrorRateReqPerMin,
		"error_rate_sensitivity":    me.ErrorRateSensitivity,
		"max_failure_rate_increase": me.MaxFailureRateIncrease,
		"minutes_abnormal_state":    me.MinutesAbnormalState,
	})
}

func (me *ErrorRateFixed) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"error_rate_req_per_min":    &me.ErrorRateReqPerMin,
		"error_rate_sensitivity":    &me.ErrorRateSensitivity,
		"max_failure_rate_increase": &me.MaxFailureRateIncrease,
		"minutes_abnormal_state":    &me.MinutesAbnormalState,
	})
}
