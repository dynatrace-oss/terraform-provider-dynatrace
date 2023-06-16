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

package resourceattribute

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RuleItems []*RuleItem

func (me *RuleItems) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(RuleItem).Schema()},
		},
	}
}

func (me RuleItems) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *RuleItems) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

type RuleItem struct {
	AttributeKey string      `json:"attributeKey"` // Attribute key **service.name** is automatically captured by default
	Enabled      bool        `json:"enabled"`      // This setting is enabled (`true`) or disabled (`false`)
	Masking      MaskingType `json:"masking"`      // Possible Values: `MASK_ENTIRE_VALUE`, `MASK_ONLY_CONFIDENTIAL_DATA`, `NOT_MASKED`
}

func (me *RuleItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute_key": {
			Type:        schema.TypeString,
			Description: "Attribute key **service.name** is automatically captured by default",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"masking": {
			Type:        schema.TypeString,
			Description: "Possible Values: `MASK_ENTIRE_VALUE`, `MASK_ONLY_CONFIDENTIAL_DATA`, `NOT_MASKED`",
			Required:    true,
		},
	}
}

func (me *RuleItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attribute_key": me.AttributeKey,
		"enabled":       me.Enabled,
		"masking":       me.Masking,
	})
}

func (me *RuleItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attribute_key": &me.AttributeKey,
		"enabled":       &me.Enabled,
		"masking":       &me.Masking,
	})
}
