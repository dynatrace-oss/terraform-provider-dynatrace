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

// GlobalEventCaptureSettings Global event capture settings
type GlobalEventCaptureSettings struct {
	MouseUp                            bool   `json:"mouseUp"`                            // MouseUp enabled/disabled
	MouseDown                          bool   `json:"mouseDown"`                          // MouseDown enabled/disabled
	Click                              bool   `json:"click"`                              // Click enabled/disabled
	DoubleClick                        bool   `json:"doubleClick"`                        // DoubleClick enabled/disabled
	KeyUp                              bool   `json:"keyUp"`                              // KeyUp enabled/disabled
	KeyDown                            bool   `json:"keyDown"`                            // KeyDown enabled/disabled
	Scroll                             bool   `json:"scroll"`                             // Scroll enabled/disabled
	AdditionalEventCapturedAsUserInput string `json:"additionalEventCapturedAsUserInput"` // Additional events to be captured globally as user input. \n\nFor example `DragStart` or `DragEnd`. Maximum 100 characters.
}

func (me *GlobalEventCaptureSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mouseup": {
			Type:        schema.TypeBool,
			Description: "MouseUp enabled/disabled",
			Optional:    true,
		},
		"mousedown": {
			Type:        schema.TypeBool,
			Description: "MouseDown enabled/disabled",
			Optional:    true,
		},
		"click": {
			Type:        schema.TypeBool,
			Description: "Click enabled/disabled",
			Optional:    true,
		},
		"doubleclick": {
			Type:        schema.TypeBool,
			Description: "DoubleClick enabled/disabled",
			Optional:    true,
		},
		"keyup": {
			Type:        schema.TypeBool,
			Description: "KeyUp enabled/disabled",
			Optional:    true,
		},
		"keydown": {
			Type:        schema.TypeBool,
			Description: "KeyDown enabled/disabled",
			Optional:    true,
		},
		"scroll": {
			Type:        schema.TypeBool,
			Description: "Scroll enabled/disabled",
			Optional:    true,
		},
		"additional_event_captured_as_user_input": {
			Type:        schema.TypeString,
			Description: "Additional events to be captured globally as user input. \n\nFor example `DragStart` or `DragEnd`. Maximum 100 characters.",
			Optional:    true,
		},
	}
}

func (me *GlobalEventCaptureSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"mouseup":     me.MouseUp,
		"mousedown":   me.MouseDown,
		"click":       me.Click,
		"doubleclick": me.DoubleClick,
		"keyup":       me.KeyUp,
		"keydown":     me.KeyDown,
		"scroll":      me.Scroll,
		"additional_event_captured_as_user_input": me.AdditionalEventCapturedAsUserInput,
	})
}

func (me *GlobalEventCaptureSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"mouseup":     &me.MouseUp,
		"mousedown":   &me.MouseDown,
		"click":       &me.Click,
		"doubleclick": &me.DoubleClick,
		"keyup":       &me.KeyUp,
		"keydown":     &me.KeyDown,
		"scroll":      &me.Scroll,
		"additional_event_captured_as_user_input": &me.AdditionalEventCapturedAsUserInput,
	})
}
