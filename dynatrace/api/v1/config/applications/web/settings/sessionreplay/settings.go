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

package sessionreplay

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Settings Session replay settings
type Settings struct {
	Enabled                            bool     `json:"enabled"`                            // SessionReplay Enabled
	CostControlPercentage              int32    `json:"costControlPercentage"`              // Session replay sampling rating in percentage
	EnableCSSResourceCapturing         bool     `json:"enableCssResourceCapturing"`         // Capture (`true`) or don't capture (`false`) CSS resources from the session
	CSSResourceCapturingExclusionRules []string `json:"cssResourceCapturingExclusionRules"` // A list of URLs to be excluded from CSS resource capturing
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "SessionReplay Enabled/Disabled",
			Optional:    true,
		},
		"cost_control_percentage": {
			Type:        schema.TypeInt,
			Description: "Session replay sampling rating in percent",
			Required:    true,
		},
		"enable_css_resource_capturing": {
			Type:        schema.TypeBool,
			Description: "Capture (`true`) or don't capture (`false`) CSS resources from the session",
			Optional:    true,
		},
		"css_resource_capturing_exclusion_rules": {
			Type:        schema.TypeList,
			Description: "A list of URLs to be excluded from CSS resource capturing",
			Optional:    true,
			// MinItems: 1,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	err := properties.EncodeAll(map[string]any{
		"enabled":                                me.Enabled,
		"cost_control_percentage":                me.CostControlPercentage,
		"enable_css_resource_capturing":          me.EnableCSSResourceCapturing,
		"css_resource_capturing_exclusion_rules": me.CSSResourceCapturingExclusionRules,
	})
	if err != nil {
		return err
	}
	if len(me.CSSResourceCapturingExclusionRules) == 0 {
		me.CSSResourceCapturingExclusionRules = nil
		properties["css_resource_capturing_exclusion_rules"] = nil
		delete(properties, "css_resource_capturing_exclusion_rules")
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"enabled":                                &me.Enabled,
		"cost_control_percentage":                &me.CostControlPercentage,
		"enable_css_resource_capturing":          &me.EnableCSSResourceCapturing,
		"css_resource_capturing_exclusion_rules": &me.CSSResourceCapturingExclusionRules,
	})
	if err != nil {
		return err
	}
	if me.CSSResourceCapturingExclusionRules == nil {
		me.CSSResourceCapturingExclusionRules = []string{}
	}
	return nil
}
