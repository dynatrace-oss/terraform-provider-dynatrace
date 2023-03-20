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

package customerrors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomErrorRules []*CustomErrorRule

func (me *CustomErrorRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CustomErrorRule).Schema()},
		},
	}
}

func (me CustomErrorRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("error_rule", me)
}

func (me *CustomErrorRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("error_rule", me)
}

type CustomErrorRule struct {
	CaptureSettings *CaptureSettings `json:"captureSettings"`        // Capture settings
	KeyMatcher      Matcher          `json:"keyMatcher"`             // Possible Values: `ALL`, `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`
	KeyPattern      *string          `json:"keyPattern,omitempty"`   // A case-insensitive key pattern
	ValueMatcher    Matcher          `json:"valueMatcher"`           // Possible Values: `ALL`, `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`
	ValuePattern    *string          `json:"valuePattern,omitempty"` // A case-insensitive value pattern
}

func (me *CustomErrorRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"capture_settings": {
			Type:        schema.TypeList,
			Description: "Capture settings",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(CaptureSettings).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"key_matcher": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ALL`, `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`",
			Required:    true,
		},
		"key_pattern": {
			Type:        schema.TypeString,
			Description: "A case-insensitive key pattern",
			Optional:    true,
		},
		"value_matcher": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ALL`, `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`",
			Required:    true,
		},
		"value_pattern": {
			Type:        schema.TypeString,
			Description: "A case-insensitive value pattern",
			Optional:    true,
		},
	}
}

func (me *CustomErrorRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"capture_settings": me.CaptureSettings,
		"key_matcher":      me.KeyMatcher,
		"key_pattern":      me.KeyPattern,
		"value_matcher":    me.ValueMatcher,
		"value_pattern":    me.ValuePattern,
	})
}

func (me *CustomErrorRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"capture_settings": &me.CaptureSettings,
		"key_matcher":      &me.KeyMatcher,
		"key_pattern":      &me.KeyPattern,
		"value_matcher":    &me.ValueMatcher,
		"value_pattern":    &me.ValuePattern,
	})
}
