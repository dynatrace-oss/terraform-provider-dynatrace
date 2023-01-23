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

// EventWrapperSettings In addition to the event handlers, events called using `addEventListener` or `attachEvent` can be captured. Be careful with this option! Event wrappers can conflict with the JavaScript code on a web page
type EventWrapperSettings struct {
	Click      bool `json:"click"`      // Click enabled/disabled
	MouseUp    bool `json:"mouseUp"`    // MouseUp enabled/disabled
	Change     bool `json:"change"`     // Change enabled/disabled
	Blur       bool `json:"blur"`       // Blur enabled/disabled
	TouchStart bool `json:"touchStart"` // TouchStart enabled/disabled
	TouchEnd   bool `json:"touchEnd"`   // TouchEnd enabled/disabled
}

func (me *EventWrapperSettings) IsDefault() bool {
	return !(me.Click || me.MouseUp || me.Change || me.Blur || me.TouchStart || me.TouchEnd)
}

func (me *EventWrapperSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"click": {
			Type:        schema.TypeBool,
			Description: "Click enabled/disabled",
			Optional:    true,
		},
		"mouseup": {
			Type:        schema.TypeBool,
			Description: "MouseUp enabled/disabled",
			Optional:    true,
		},
		"change": {
			Type:        schema.TypeBool,
			Description: "Change enabled/disabled",
			Optional:    true,
		},
		"blur": {
			Type:        schema.TypeBool,
			Description: "Blur enabled/disabled",
			Optional:    true,
		},
		"touch_start": {
			Type:        schema.TypeBool,
			Description: "TouchStart enabled/disabled",
			Optional:    true,
		},
		"touch_end": {
			Type:        schema.TypeBool,
			Description: "TouchEnd enabled/disabled",
			Optional:    true,
		},
	}
}

func (me *EventWrapperSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"click":       me.Click,
		"mouseup":     me.MouseUp,
		"change":      me.Change,
		"blur":        me.Blur,
		"touch_start": me.TouchStart,
		"touch_end":   me.TouchEnd,
	})
}

func (me *EventWrapperSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"click":       &me.Click,
		"mouseup":     &me.MouseUp,
		"change":      &me.Change,
		"blur":        &me.Blur,
		"touch_start": &me.TouchStart,
		"touch_end":   &me.TouchEnd,
	})
}
