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

package hosts

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HighGcActivityDetection struct {
	CustomThresholds *HighGcActivityDetectionThresholds `json:"customThresholds,omitempty"` // Alert if the GC time **or** the GC suspension is exceeded
	DetectionMode    *DetectionMode                     `json:"detectionMode,omitempty"`    // Detection mode for high GC activity. Possible Values: `auto`, `custom`
	Enabled          bool                               `json:"enabled"`                    // This setting is enabled (`true`) or disabled (`false`)
}

func (me *HighGcActivityDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_thresholds": {
			Type:        schema.TypeList,
			Description: "Alert if the GC time **or** the GC suspension is exceeded",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(HighGcActivityDetectionThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Detection mode for high GC activity. Possible Values: `auto`, `custom`",
			Optional:    true, // precondition
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
	}
}

func (me *HighGcActivityDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"custom_thresholds": me.CustomThresholds,
		"detection_mode":    me.DetectionMode,
		"enabled":           me.Enabled,
	})
}

func (me *HighGcActivityDetection) HandlePreconditions() error {
	if (me.CustomThresholds == nil) && (me.Enabled && (me.DetectionMode != nil && (string(*me.DetectionMode) == "custom"))) {
		return fmt.Errorf("'custom_thresholds' must be specified if 'enabled' is set to '%v' and 'detection_mode' is set to '%v'", me.Enabled, me.DetectionMode)
	}
	if (me.CustomThresholds != nil) && (!me.Enabled || (me.DetectionMode == nil || (me.DetectionMode != nil && string(*me.DetectionMode) != "custom"))) {
		return fmt.Errorf("'custom_thresholds' must not be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if (me.DetectionMode == nil) && (me.Enabled) {
		return fmt.Errorf("'detection_mode' must be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if (me.DetectionMode != nil) && (!me.Enabled) {
		return fmt.Errorf("'detection_mode' must not be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	return nil
}

func (me *HighGcActivityDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_thresholds": &me.CustomThresholds,
		"detection_mode":    &me.DetectionMode,
		"enabled":           &me.Enabled,
	})
}
