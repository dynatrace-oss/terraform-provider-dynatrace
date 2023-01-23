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

// AdditionalEventHandlers Additional event handlers and wrappers
type AdditionalEventHandlers struct {
	UseMouseUpEventForClicks bool  `json:"userMouseupEventForClicks"` // Use mouseup event for clicks enabled/disabled
	ClickEventHandler        bool  `json:"clickEventHandler"`         // Click event handler enabled/disabled
	MouseUpEventHandler      bool  `json:"mouseupEventHandler"`       // Mouseup event handler enabled/disabled
	BlurEventHandler         bool  `json:"blurEventHandler"`          // Blur event handler enabled/disabled
	ChangeEventHandler       bool  `json:"changeEventHandler"`        // Change event handler enabled/disabled
	ToStringMethod           bool  `json:"toStringMethod"`            // toString method enabled/disabled
	MaxDomNodesToInstrument  int32 `json:"maxDomNodesToInstrument"`   // Max. number of DOM nodes to instrument. Valid values range from 0 to 100000.
}

func (me *AdditionalEventHandlers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"use_mouse_up_event_for_clicks": {
			Type:        schema.TypeBool,
			Description: "Use mouseup event for clicks enabled/disabled",
			Optional:    true,
		},
		"click": {
			Type:        schema.TypeBool,
			Description: "Click event handler enabled/disabled",
			Optional:    true,
		},
		"mouseup": {
			Type:        schema.TypeBool,
			Description: "Mouseup event handler enabled/disabled",
			Optional:    true,
		},
		"blur": {
			Type:        schema.TypeBool,
			Description: "Blur event handler enabled/disabled",
			Optional:    true,
		},
		"change": {
			Type:        schema.TypeBool,
			Description: "Change event handler enabled/disabled",
			Optional:    true,
		},
		"to_string_method": {
			Type:        schema.TypeBool,
			Description: "toString method enabled/disabled",
			Optional:    true,
		},
		"max_dom_nodes": {
			Type:        schema.TypeInt,
			Description: "Max. number of DOM nodes to instrument. Valid values range from 0 to 100000.",
			Required:    true,
		},
	}
}

func (me *AdditionalEventHandlers) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"use_mouse_up_event_for_clicks": me.UseMouseUpEventForClicks,
		"click":                         me.ClickEventHandler,
		"mouseup":                       me.MouseUpEventHandler,
		"blur":                          me.BlurEventHandler,
		"change":                        me.ChangeEventHandler,
		"to_string_method":              me.ToStringMethod,
		"max_dom_nodes":                 me.MaxDomNodesToInstrument,
	})
}

func (me *AdditionalEventHandlers) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"use_mouse_up_event_for_clicks": &me.UseMouseUpEventForClicks,
		"click":                         &me.ClickEventHandler,
		"mouseup":                       &me.MouseUpEventHandler,
		"blur":                          &me.BlurEventHandler,
		"change":                        &me.ChangeEventHandler,
		"to_string_method":              &me.ToStringMethod,
		"max_dom_nodes":                 &me.MaxDomNodesToInstrument,
	})
}
