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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HighGcActivityDetection struct {
	CustomThresholds *HighGcActivityDetectionThresholds `json:"customThresholds,omitempty"` // Alert if the GC time **or** the GC suspension is exceeded
	DetectionMode    *DetectionMode                     `json:"detectionMode,omitempty"`    // Detection mode for high GC activity
	Enabled          bool                               `json:"enabled"`                    // You may also configure high GC activity alerting for .NET processes on [extensions events page](/#settings/anomalydetection/extensionevents).
}

func (me *HighGcActivityDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_thresholds": {
			Type:        schema.TypeList,
			Description: "Alert if the GC time **or** the GC suspension is exceeded",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(HighGcActivityDetectionThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"detection_mode": {
			Type:        schema.TypeString,
			Description: "Detection mode for high GC activity",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "You may also configure high GC activity alerting for .NET processes on [extensions events page](/#settings/anomalydetection/extensionevents).",
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

func (me *HighGcActivityDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_thresholds": &me.CustomThresholds,
		"detection_mode":    &me.DetectionMode,
		"enabled":           &me.Enabled,
	})
}
