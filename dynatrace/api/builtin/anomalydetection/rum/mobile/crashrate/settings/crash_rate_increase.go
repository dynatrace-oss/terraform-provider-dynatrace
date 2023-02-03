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

package rummobilecrashrateincrease

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CrashRateIncrease struct {
	CrashRateIncreaseAuto  *CrashRateIncreaseAuto  `json:"crashRateIncreaseAuto,omitempty"`  // Alert crash rate increases when auto-detected baseline is exceeded by a certain number of users
	CrashRateIncreaseFixed *CrashRateIncreaseFixed `json:"crashRateIncreaseFixed,omitempty"` // Alert crash rate increases when the defined threshold is exceeded by a certain number of users
	DetectionMode          *DetectionMode          `json:"detectionMode,omitempty"`          // Possible Values: `Auto`, `Fixed`
	Enabled                bool                    `json:"enabled"`                          // This setting is enabled (`true`) or disabled (`false`)
}

func (me *CrashRateIncrease) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"crash_rate_increase_auto": {
			Type:        schema.TypeList,
			Description: "Alert crash rate increases when auto-detected baseline is exceeded by a certain number of users",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(CrashRateIncreaseAuto).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"crash_rate_increase_fixed": {
			Type:        schema.TypeList,
			Description: "Alert crash rate increases when the defined threshold is exceeded by a certain number of users",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(CrashRateIncreaseFixed).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Auto`, `Fixed`",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
	}
}

func (me *CrashRateIncrease) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"crash_rate_increase_auto":  me.CrashRateIncreaseAuto,
		"crash_rate_increase_fixed": me.CrashRateIncreaseFixed,
		"detection_mode":            me.DetectionMode,
		"enabled":                   me.Enabled,
	})
}

func (me *CrashRateIncrease) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"crash_rate_increase_auto":  &me.CrashRateIncreaseAuto,
		"crash_rate_increase_fixed": &me.CrashRateIncreaseFixed,
		"detection_mode":            &me.DetectionMode,
		"enabled":                   &me.Enabled,
	})
}
