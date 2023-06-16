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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LoadSpikes struct {
	Enabled              bool     `json:"enabled"`                        // This setting is enabled (`true`) or disabled (`false`)
	LoadSpikePercent     *float64 `json:"loadSpikePercent,omitempty"`     // Threshold
	MinutesAbnormalState *int     `json:"minutesAbnormalState,omitempty"` // Time span
}

func (me *LoadSpikes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"load_spike_percent": {
			Type:        schema.TypeFloat,
			Description: "Threshold",
			Optional:    true, // precondition
		},
		"minutes_abnormal_state": {
			Type:        schema.TypeInt,
			Description: "Time span",
			Optional:    true, // precondition
		},
	}
}

func (me *LoadSpikes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":                me.Enabled,
		"load_spike_percent":     me.LoadSpikePercent,
		"minutes_abnormal_state": me.MinutesAbnormalState,
	})
}

func (me *LoadSpikes) HandlePreconditions() error {
	if me.LoadSpikePercent == nil && me.Enabled {
		return fmt.Errorf("'load_spike_percent' must be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if me.MinutesAbnormalState == nil && me.Enabled {
		return fmt.Errorf("'minutes_abnormal_state' must be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	return nil
}

func (me *LoadSpikes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":                &me.Enabled,
		"load_spike_percent":     &me.LoadSpikePercent,
		"minutes_abnormal_state": &me.MinutesAbnormalState,
	})
}
