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

type OverAlertingProtection struct {
	MinutesAbnormalState int     `json:"minutesAbnormalState"` // Only alert if the abnormal state remains for at least
	RequestsPerMinute    float64 `json:"requestsPerMinute"`    // Only alert if there are at least
}

func (me *OverAlertingProtection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"minutes_abnormal_state": {
			Type:        schema.TypeInt,
			Description: "Only alert if the abnormal state remains for at least",
			Required:    true,
		},
		"requests_per_minute": {
			Type:        schema.TypeFloat,
			Description: "Only alert if there are at least",
			Required:    true,
		},
	}
}

func (me *OverAlertingProtection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"minutes_abnormal_state": me.MinutesAbnormalState,
		"requests_per_minute":    me.RequestsPerMinute,
	})
}

func (me *OverAlertingProtection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"minutes_abnormal_state": &me.MinutesAbnormalState,
		"requests_per_minute":    &me.RequestsPerMinute,
	})
}
