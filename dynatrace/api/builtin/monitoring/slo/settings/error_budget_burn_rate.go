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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ErrorBudgetBurnRate struct {
	BurnRateVisualizationEnabled bool     `json:"burnRateVisualizationEnabled"` // Burn rate visualization enabled
	FastBurnThreshold            *float64 `json:"fastBurnThreshold,omitempty"`  // The threshold defines when a burn rate is marked as fast-burning (high-emergency). Burn rates lower than this threshold (and greater than 1) are highlighted as slow-burn (low-emergency).
}

func (me *ErrorBudgetBurnRate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"burn_rate_visualization_enabled": {
			Type:        schema.TypeBool,
			Description: "Burn rate visualization enabled",
			Required:    true,
		},
		"fast_burn_threshold": {
			Type:        schema.TypeFloat,
			Description: "The threshold defines when a burn rate is marked as fast-burning (high-emergency). Burn rates lower than this threshold (and greater than 1) are highlighted as slow-burn (low-emergency).",
			Optional:    true, // precondition
		},
	}
}

func (me *ErrorBudgetBurnRate) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"burn_rate_visualization_enabled": me.BurnRateVisualizationEnabled,
		"fast_burn_threshold":             me.FastBurnThreshold,
	})
}

func (me *ErrorBudgetBurnRate) HandlePreconditions() error {
	if me.FastBurnThreshold == nil && me.BurnRateVisualizationEnabled {
		return fmt.Errorf("'fast_burn_threshold' must be specified if 'burn_rate_visualization_enabled' is set to '%v'", me.BurnRateVisualizationEnabled)
	}
	return nil
}

func (me *ErrorBudgetBurnRate) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"burn_rate_visualization_enabled": &me.BurnRateVisualizationEnabled,
		"fast_burn_threshold":             &me.FastBurnThreshold,
	})
}
