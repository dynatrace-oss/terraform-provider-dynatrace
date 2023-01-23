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

package comparisoninfo

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FastString Comparison for `FAST_STRING` attributes. Use it for all service property attributes.
type FastString struct {
	BaseComparisonInfo
	CaseSensitive *bool                `json:"caseSensitive,omitempty"` // The comparison is case-sensitive (`true`) or not case-sensitive (`false`).
	Comparison    FastStringComparison `json:"comparison"`              // Operator of the comparision. You can reverse it by setting **negate** to `true`.
	Value         *string              `json:"value,omitempty"`         // The value to compare to.
	Values        []string             `json:"values,omitempty"`        // The values to compare to.
}

func (me *FastString) GetType() Type {
	return Types.FastString
}

func (me *FastString) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The comparison is case-sensitive (`true`) or not case-sensitive (`false`)",
		},
		"values": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The values to compare to",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to compare to",
		},
		"operator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Operator of the comparison. You can reverse it by setting `negate` to `true`. Possible values are `EQUALS`, `EQUALS_ANY_OF` and `CONTAINS`",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *FastString) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"values":         me.Values,
		"value":          me.Value,
		"operator":       me.Comparison,
		"case_sensitive": me.CaseSensitive,
		"unknowns":       me.Unknowns,
	})
}

func (me *FastString) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"values":         &me.Values,
		"value":          &me.Value,
		"operator":       &me.Comparison,
		"case_sensitive": &me.CaseSensitive,
		"unknowns":       &me.Unknowns,
	})
}

func (me *FastString) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"type":          me.GetType(),
		"negate":        me.Negate,
		"values":        me.Values,
		"value":         me.Value,
		"comparison":    me.Comparison,
		"caseSensitive": me.CaseSensitive,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *FastString) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"negate":        &me.Negate,
		"values":        &me.Values,
		"value":         &me.Value,
		"comparison":    &me.Comparison,
		"caseSensitive": &me.CaseSensitive,
	})
}

// FastStringComparison Operator of the comparision. You can reverse it by setting **negate** to `true`.
type FastStringComparison string

// FastStringComparisons offers the known enum values
var FastStringComparisons = struct {
	Contains    FastStringComparison
	Equals      FastStringComparison
	EqualsAnyOf FastStringComparison
}{
	"CONTAINS",
	"EQUALS",
	"EQUALS_ANY_OF",
}
