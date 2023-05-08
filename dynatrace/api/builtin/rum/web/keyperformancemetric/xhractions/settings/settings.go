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

package xhractions

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	FallbackThresholds *FallbackThresholds `json:"fallbackThresholds,omitempty"` // (Field has overlap with `dynatrace_web_application`) If the selected key performance metric is not detected, the **User action duration** metric is used instead.
	Kpm                XhrKpm              `json:"kpm"`                          // (Field has overlap with `dynatrace_web_application`) Possible Values: `RESPONSE_END`, `RESPONSE_START`, `USER_ACTION_DURATION`, `VISUALLY_COMPLETE`
	Scope              string              `json:"-" scope:"scope"`              // The scope of this setting (APPLICATION_METHOD, APPLICATION)
	Thresholds         *Thresholds         `json:"thresholds"`                   // (Field has overlap with `dynatrace_web_application`) Set the Tolerating and Frustrated performance thresholds for this action type.
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fallback_thresholds": {
			Type:        schema.TypeList,
			Description: "If the selected key performance metric is not detected, the **User action duration** metric is used instead.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FallbackThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"kpm": {
			Type:        schema.TypeString,
			Description: "Possible Values: `RESPONSE_END`, `RESPONSE_START`, `USER_ACTION_DURATION`, `VISUALLY_COMPLETE`",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (APPLICATION_METHOD, APPLICATION)",
			Required:    true,
		},
		"thresholds": {
			Type:        schema.TypeList,
			Description: "Set the Tolerating and Frustrated performance thresholds for this action type.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Thresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"fallback_thresholds": me.FallbackThresholds,
		"kpm":                 me.Kpm,
		"scope":               me.Scope,
		"thresholds":          me.Thresholds,
	})
}

func (me *Settings) HandlePreconditions() error {
	if me.FallbackThresholds == nil && (string(me.Kpm) != "USER_ACTION_DURATION") {
		return fmt.Errorf("'fallback_thresholds' must be specified if 'kpm' is set to '%v'", me.Kpm)
	}
	if me.FallbackThresholds != nil && (string(me.Kpm) == "USER_ACTION_DURATION") {
		return fmt.Errorf("'fallback_thresholds' must not be specified if 'kpm' is set to '%v'", me.Kpm)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"fallback_thresholds": &me.FallbackThresholds,
		"kpm":                 &me.Kpm,
		"scope":               &me.Scope,
		"thresholds":          &me.Thresholds,
	})
}
