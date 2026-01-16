/**
* @license
* Copyright 2026 Dynatrace LLC
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

// StringOneAgentAttribute Comparison for `STRING_ONE_AGENT_ATTRIBUTE` attributes, specifically of the **String** type.
type StringOneAgentAttribute struct {
	BaseComparisonInfo
	Value                *string                          `json:"value,omitempty"`         // The value to compare to.
	Values               []string                         `json:"values,omitempty"`        // The values to compare to.
	CaseSensitive        *bool                            `json:"caseSensitive,omitempty"` // The comparison is case-sensitive (`true`) or not case-sensitive (`false`).
	Comparison           StringRequestAttributeComparison `json:"comparison"`              // Operator of the comparision. You can reverse it by setting **negate** to `true`.
	OneAgentAttributeKey string                           `json:"oneAgentAttributeKey"`    // The One Agent attribute to extract from.
}

func (me *StringOneAgentAttribute) GetType() Type {
	return Types.StringOneAgentAttribute
}

func (me *StringOneAgentAttribute) Schema() map[string]*schema.Schema {
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
			Description: "Operator of the comparison. You can reverse it by setting `negate` to `true`. Possible values are `BEGINS_WITH`, `BEGINS_WITH_ANY_OF`, `CONTAINS`, `ENDS_WITH`, `ENDS_WITH_ANY_OF`, `EQUALS`, `EQUALS_ANY_OF`, `EXISTS` and `REGEX_MATCHES`",
		},
		"one_agent_attribute_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OneAgent attribute to extract from",
		},
	}
}

func (me *StringOneAgentAttribute) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"values":                  me.Values,
		"value":                   me.Value,
		"operator":                me.Comparison,
		"case_sensitive":          me.CaseSensitive,
		"one_agent_attribute_key": me.OneAgentAttributeKey,
	})
}

func (me *StringOneAgentAttribute) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"values":                  &me.Values,
		"value":                   &me.Value,
		"operator":                &me.Comparison,
		"case_sensitive":          &me.CaseSensitive,
		"one_agent_attribute_key": &me.OneAgentAttributeKey,
	})
}

func (me *StringOneAgentAttribute) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"type":                 me.GetType(),
		"negate":               me.Negate,
		"values":               me.Values,
		"value":                me.Value,
		"comparison":           me.Comparison,
		"caseSensitive":        me.CaseSensitive,
		"oneAgentAttributeKey": me.OneAgentAttributeKey,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *StringOneAgentAttribute) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"negate":               &me.Negate,
		"values":               &me.Values,
		"value":                &me.Value,
		"comparison":           &me.Comparison,
		"caseSensitive":        &me.CaseSensitive,
		"oneAgentAttributeKey": &me.OneAgentAttributeKey,
	})
}

// StringOneAgentAttributeComparison Operator of the comparision. You can reverse it by setting **negate** to `true`.
type StringOneAgentAttributeComparison string

// StringOneAgentAttributeComparisons offers the known enum values
var StringOneAgentAttributeComparisons = struct {
	BeginsWith      StringOneAgentAttributeComparison
	BeginsWithAnyOf StringOneAgentAttributeComparison
	Contains        StringOneAgentAttributeComparison
	EndsWith        StringOneAgentAttributeComparison
	EndsWithAnyOf   StringOneAgentAttributeComparison
	Equals          StringOneAgentAttributeComparison
	EqualsAnyOf     StringOneAgentAttributeComparison
	Exists          StringOneAgentAttributeComparison
	RegexMatches    StringOneAgentAttributeComparison
}{
	"BEGINS_WITH",
	"BEGINS_WITH_ANY_OF",
	"CONTAINS",
	"ENDS_WITH",
	"ENDS_WITH_ANY_OF",
	"EQUALS",
	"EQUALS_ANY_OF",
	"EXISTS",
	"REGEX_MATCHES",
}
