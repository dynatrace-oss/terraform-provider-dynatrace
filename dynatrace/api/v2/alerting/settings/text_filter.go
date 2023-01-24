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

package alerting

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TextFilter struct {
	Operator      Operator `json:"operator"`      // Operator of the comparison
	Value         string   `json:"value"`         // The value to compare with
	Negate        bool     `json:"negate"`        // Negate the operator
	Enabled       bool     `json:"enabled"`       // Enable this filter
	CaseSensitive bool     `json:"caseSensitive"` // Case Sensitive comparison of text
}

func (me *TextFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The filter is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "Reverses the comparison **operator**. For example it turns the **begins with** into **does not begin with**",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Operator of the comparison.   You can reverse it by setting **negate** to `true`. Possible values are `BEGINS_WITH`, `CONTAINS`, `CONTAINS_REGEX`, `ENDS_WITH` and `EQUALS`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to",
			Required:    true,
		},
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "The condition is case sensitive (`false`) or case insensitive (`true`).   If not set, then `false` is used, making the condition case sensitive",
			Optional:    true,
		},
	}
}

func (me *TextFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":        me.Enabled,
		"negate":         me.Negate,
		"operator":       string(me.Operator),
		"value":          me.Value,
		"case_sensitive": me.CaseSensitive,
	})
}

func (me *TextFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("negate"); ok {
		me.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("case_sensitive"); ok {
		me.CaseSensitive = (value.(bool))
	}
	if value, ok := decoder.GetOk("operator"); ok {
		me.Operator = Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = value.(string)
	}
	return nil
}
