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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ContentCapture contains settings for content capture
type ContentCapture struct {
	ResourceTimingSettings        *ResourceTimingSettings   `json:"resourceTimingSettings"`        // Settings for resource timings capture
	JavaScriptErrors              bool                      `json:"javaScriptErrors"`              // JavaScript errors monitoring enabled/disabled
	TimeoutSettings               *TimeoutSettings          `json:"timeoutSettings"`               // Settings for timed action capture
	VisuallyCompleteAndSpeedIndex bool                      `json:"visuallyCompleteAndSpeedIndex"` // Visually complete and Speed index support enabled/disabled
	VisuallyCompleteSettings      *VisuallyCompleteSettings `json:"visuallyComplete2Settings"`     // Settings for VisuallyComplete
}

func (me *ContentCapture) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_timing_settings": {
			Type:        schema.TypeList,
			Description: "Settings for resource timings capture",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ResourceTimingSettings).Schema()},
		},
		"javascript_errors": {
			Type:        schema.TypeBool,
			Description: "JavaScript errors monitoring enabled/disabled",
			Optional:    true,
		},
		"timeout_settings": {
			Type:        schema.TypeList,
			Description: "Settings for timed action capture",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(TimeoutSettings).Schema()},
		},
		"visually_complete_and_speed_index": {
			Type:        schema.TypeBool,
			Description: "Visually complete and Speed index support enabled/disabled",
			Optional:    true,
		},
		"visually_complete_settings": {
			Type:        schema.TypeList,
			Description: "Settings for VisuallyComplete",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(VisuallyCompleteSettings).Schema()},
		},
	}
}

func (me *ContentCapture) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"resource_timing_settings":          me.ResourceTimingSettings,
		"javascript_errors":                 me.JavaScriptErrors,
		"timeout_settings":                  me.TimeoutSettings,
		"visually_complete_and_speed_index": me.VisuallyCompleteAndSpeedIndex,
		"visually_complete_settings":        me.VisuallyCompleteSettings,
	})
}

func (me *ContentCapture) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"resource_timing_settings":          &me.ResourceTimingSettings,
		"javascript_errors":                 &me.JavaScriptErrors,
		"timeout_settings":                  &me.TimeoutSettings,
		"visually_complete_and_speed_index": &me.VisuallyCompleteAndSpeedIndex,
		"visually_complete_settings":        &me.VisuallyCompleteSettings,
	})
}
