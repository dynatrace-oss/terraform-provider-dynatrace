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

package namespace

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PodsQuotaSaturationConfig struct {
	ObservationPeriodInMinutes int `json:"observationPeriodInMinutes"` // within the last
	SamplePeriodInMinutes      int `json:"samplePeriodInMinutes"`      // of quota for at least
	Threshold                  int `json:"threshold"`                  // number of utilized namespace pods is above
}

func (me *PodsQuotaSaturationConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"observation_period_in_minutes": {
			Type:        schema.TypeInt,
			Description: "within the last",
			Required:    true,
		},
		"sample_period_in_minutes": {
			Type:        schema.TypeInt,
			Description: "of quota for at least",
			Required:    true,
		},
		"threshold": {
			Type:        schema.TypeInt,
			Description: "number of utilized namespace pods is above",
			Required:    true,
		},
	}
}

func (me *PodsQuotaSaturationConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"observation_period_in_minutes": me.ObservationPeriodInMinutes,
		"sample_period_in_minutes":      me.SamplePeriodInMinutes,
		"threshold":                     me.Threshold,
	})
}

func (me *PodsQuotaSaturationConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"observation_period_in_minutes": &me.ObservationPeriodInMinutes,
		"sample_period_in_minutes":      &me.SamplePeriodInMinutes,
		"threshold":                     &me.Threshold,
	})
}
