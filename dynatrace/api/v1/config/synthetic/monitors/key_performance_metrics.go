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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// KeyPerformanceMetrics The key performance metrics configuration
type KeyPerformanceMetrics struct {
	LoadActionKPM LoadActionKPM `json:"loadActionKpm"` // Defines the key performance metric for load actions
	XHRActionKPM  XHRActionKPM  `json:"xhrActionKpm"`  // Defines the key performance metric for XHR actions
}

func (me *KeyPerformanceMetrics) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load_action_kpm": {
			Type:        schema.TypeString,
			Description: "Defines the key performance metric for load actions. Supported values are `VISUALLY_COMPLETE`, `SPEED_INDEX`, `USER_ACTION_DURATION`, `TIME_TO_FIRST_BYTE`, `HTML_DOWNLOADED`, `DOM_INTERACTIVE`, `LOAD_EVENT_START` and `LOAD_EVENT_END`.",
			Required:    true,
		},
		"xhr_action_kpm": {
			Type:        schema.TypeString,
			Description: "Defines the key performance metric for XHR actions. Supported values are `VISUALLY_COMPLETE`, `USER_ACTION_DURATION`, `TIME_TO_FIRST_BYTE` and `RESPONSE_END`.",
			Required:    true,
		},
	}
}

func (me *KeyPerformanceMetrics) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("load_action_kpm", string(me.LoadActionKPM)); err != nil {
		return err
	}
	if err := properties.Encode("xhr_action_kpm", string(me.XHRActionKPM)); err != nil {
		return err
	}
	return nil
}

func (me *KeyPerformanceMetrics) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("load_action_kpm", &me.LoadActionKPM); err != nil {
		return err
	}
	if err := decoder.Decode("xhr_action_kpm", &me.XHRActionKPM); err != nil {
		return err
	}
	return nil
}

// LoadActionKPM Defines the key performance metric for load actions
type LoadActionKPM string

// LoadActionKPMs offers the known enum values
var LoadActionKPMs = struct {
	VisuallyComplete   LoadActionKPM
	SpeedIndex         LoadActionKPM
	UserActionDuration LoadActionKPM
	TimeToFirstByte    LoadActionKPM
	HTMLDownloaded     LoadActionKPM
	DOMInteractive     LoadActionKPM
	LoadEventStart     LoadActionKPM
	LoadEventEnd       LoadActionKPM
}{
	"VISUALLY_COMPLETE",
	"SPEED_INDEX",
	"USER_ACTION_DURATION",
	"TIME_TO_FIRST_BYTE",
	"HTML_DOWNLOADED",
	"DOM_INTERACTIVE",
	"LOAD_EVENT_START",
	"LOAD_EVENT_END",
}

// XHRActionKPM Defines the key performance metric for XHR actions
type XHRActionKPM string

// LoadActionKPMs offers the known enum values
var XHRActionKPMs = struct {
	VisuallyComplete   XHRActionKPM
	UserActionDuration XHRActionKPM
	TimeToFirstByte    XHRActionKPM
	ResponseEnd        XHRActionKPM
}{
	"VISUALLY_COMPLETE",
	"USER_ACTION_DURATION",
	"TIME_TO_FIRST_BYTE",
	"RESPONSE_END",
}
