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

package slo

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ErrorBudgetBurnRate struct {
	BurnRateVisualizationEnabled *bool    `json:"burnRateVisualizationEnabled,omitempty"` // The error budget burn rate calculation is enabled (true) or disabled (false).
	FastBurnThreshold            *float64 `json:"fastBurnThreshold,omitempty"`            // The threshold between a slow and a fast burn rate.
}

func (me *ErrorBudgetBurnRate) IsEmpty() bool {
	if me == nil {
		return true
	}
	if me.BurnRateVisualizationEnabled != nil && *me.BurnRateVisualizationEnabled {
		return false
	}
	if me.FastBurnThreshold != nil && *me.FastBurnThreshold == 0.0 {
		return false
	}
	return true
}

func (me *ErrorBudgetBurnRate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"burn_rate_visualization_enabled": {
			Type:        schema.TypeBool,
			Description: "The error budget burn rate calculation is enabled (true) or disabled (false).",
			Optional:    true,
		},
		"fast_burn_threshold": {
			Type:        schema.TypeFloat,
			Description: "The threshold between a slow and a fast burn rate.",
			Optional:    true,
		},
	}
}

func (me *ErrorBudgetBurnRate) MarshalHCL(properties hcl.Properties) error {
	err := properties.EncodeAll(map[string]any{
		"burn_rate_visualization_enabled": me.BurnRateVisualizationEnabled,
		"fast_burn_threshold":             me.FastBurnThreshold,
	})
	return err
}

func (me *ErrorBudgetBurnRate) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"burn_rate_visualization_enabled": &me.BurnRateVisualizationEnabled,
		"fast_burn_threshold":             &me.FastBurnThreshold,
	})
}
