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

package rummobile

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SlowUserActions struct {
	DetectionMode        *DetectionMode        `json:"detectionMode,omitempty"` // Possible Values: `Auto`, `Fixed`
	Enabled              bool                  `json:"enabled"`                 // This setting is enabled (`true`) or disabled (`false`)
	SlowUserActionsAuto  *SlowUserActionsAuto  `json:"slowUserActionsAuto,omitempty"`
	SlowUserActionsFixed *SlowUserActionsFixed `json:"slowUserActionsFixed,omitempty"`
}

func (me *SlowUserActions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"slow_user_actions_auto": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(SlowUserActionsAuto).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"slow_user_actions_fixed": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(SlowUserActionsFixed).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *SlowUserActions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"detection_mode":          me.DetectionMode,
		"enabled":                 me.Enabled,
		"slow_user_actions_auto":  me.SlowUserActionsAuto,
		"slow_user_actions_fixed": me.SlowUserActionsFixed,
	})
}

func (me *SlowUserActions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"detection_mode":          &me.DetectionMode,
		"enabled":                 &me.Enabled,
		"slow_user_actions_auto":  &me.SlowUserActionsAuto,
		"slow_user_actions_fixed": &me.SlowUserActionsFixed,
	})
}
