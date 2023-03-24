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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No documentation available
type LoadDrops struct {
	Enabled              bool     `json:"enabled"`                        // Detect service load drops
	LoadDropPercent      *float64 `json:"loadDropPercent,omitempty"`      // Threshold
	MinutesAbnormalState *int     `json:"minutesAbnormalState,omitempty"` // Time span
}

func (me *LoadDrops) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Detect service load drops",
			Required:    true,
		},
		"load_drop_percent": {
			Type:        schema.TypeFloat,
			Description: "Threshold",
			Optional:    true,
		},
		"minutes_abnormal_state": {
			Type:        schema.TypeInt,
			Description: "Time span",
			Optional:    true,
		},
	}
}

func (me *LoadDrops) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":                me.Enabled,
		"load_drop_percent":      me.LoadDropPercent,
		"minutes_abnormal_state": me.MinutesAbnormalState,
	})
}

func (me *LoadDrops) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"enabled":                &me.Enabled,
		"load_drop_percent":      &me.LoadDropPercent,
		"minutes_abnormal_state": &me.MinutesAbnormalState,
	})
	if me.LoadDropPercent == nil && me.Enabled {
		me.LoadDropPercent = opt.NewFloat64(0)
	}
	return err
}
