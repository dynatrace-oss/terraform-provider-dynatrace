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

package privacypreferences

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AllowListRules []*AllowListRule

func (me *AllowListRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_list_rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(AllowListRule).Schema()},
		},
	}
}

func (me AllowListRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("allow_list_rule", me)
}

func (me *AllowListRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("allow_list_rule", me)
}

type AllowListRule struct {
	AttributeExpression *string           `json:"attributeExpression,omitempty"` // Attribute masking can be applied to web applications that store data within attributes, typically data-NAME attributes in HTML5. When you define attributes, their values are masked while recording but not removed.
	CssExpression       *string           `json:"cssExpression,omitempty"`       // Content masking can be applied to webpages where personal data is displayed. When content masking is applied to parent elements, all child elements are masked by default.
	Target              MaskingTargetType `json:"target"`                        // Possible Values: `ATTRIBUTE`, `ELEMENT`
}

func (me *AllowListRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute_expression": {
			Type:        schema.TypeString,
			Description: "Attribute masking can be applied to web applications that store data within attributes, typically data-NAME attributes in HTML5. When you define attributes, their values are masked while recording but not removed.",
			Optional:    true,
		},
		"css_expression": {
			Type:        schema.TypeString,
			Description: "Content masking can be applied to webpages where personal data is displayed. When content masking is applied to parent elements, all child elements are masked by default.",
			Optional:    true,
		},
		"target": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ATTRIBUTE`, `ELEMENT`",
			Required:    true,
		},
	}
}

func (me *AllowListRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attribute_expression": me.AttributeExpression,
		"css_expression":       me.CssExpression,
		"target":               me.Target,
	})
}

func (me *AllowListRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attribute_expression": &me.AttributeExpression,
		"css_expression":       &me.CssExpression,
		"target":               &me.Target,
	})
}
