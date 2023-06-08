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

type FailureRate struct {
	AutoDetection  *FailureRateAuto  `json:"autoDetection,omitempty"`  // Alert if the percentage of failing service calls increases by **both** the absolute and relative thresholds:
	DetectionMode  *DetectionMode    `json:"detectionMode,omitempty"`  // Possible Values: `Auto`, `Fixed`
	Enabled        bool              `json:"enabled"`                  // This setting is enabled (`true`) or disabled (`false`)
	FixedDetection *FailureRateFixed `json:"fixedDetection,omitempty"` // Alert if a given failure rate is exceeded during any 5-minute-period
}

func (me *FailureRate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_detection": {
			Type:        schema.TypeList,
			Description: "Alert if the percentage of failing service calls increases by **both** the absolute and relative thresholds:",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FailureRateAuto).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Auto`, `Fixed`",
			Optional:    true, // precondition
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"fixed_detection": {
			Type:        schema.TypeList,
			Description: "Alert if a given failure rate is exceeded during any 5-minute-period",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FailureRateFixed).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *FailureRate) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_detection":  me.AutoDetection,
		"detection_mode":  me.DetectionMode,
		"enabled":         me.Enabled,
		"fixed_detection": me.FixedDetection,
	})
}

func (me *FailureRate) HandlePreconditions() error {
	if me.AutoDetection == nil && me.Enabled && (me.DetectionMode != nil && string(*me.DetectionMode) == "auto") {
		return fmt.Errorf("'auto_detection' must be specified if 'enabled' is set to '%v' and 'detection_mode' is set to '%v'", me.Enabled, me.DetectionMode)
	}
	if me.AutoDetection != nil && (!me.Enabled || me.DetectionMode == nil || (me.DetectionMode != nil && string(*me.DetectionMode) != "auto")) {
		return fmt.Errorf("'auto_detection' must not be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if me.DetectionMode == nil && me.Enabled {
		return fmt.Errorf("'detection_mode' must be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if me.DetectionMode != nil && !me.Enabled {
		return fmt.Errorf("'detection_mode' must not be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if me.FixedDetection == nil && me.Enabled && (me.DetectionMode != nil && string(*me.DetectionMode) == "fixed") {
		return fmt.Errorf("'fixed_detection' must be specified if 'enabled' is set to '%v' and 'detection_mode' is set to '%v'", me.Enabled, me.DetectionMode)
	}
	if me.FixedDetection != nil && (!me.Enabled || me.DetectionMode == nil || (me.DetectionMode != nil && string(*me.DetectionMode) != "fixed")) {
		return fmt.Errorf("'fixed_detection' must not be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	return nil
}

func (me *FailureRate) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_detection":  &me.AutoDetection,
		"detection_mode":  &me.DetectionMode,
		"enabled":         &me.Enabled,
		"fixed_detection": &me.FixedDetection,
	})
}
