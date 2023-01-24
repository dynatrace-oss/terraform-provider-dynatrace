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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service/settings/propagation"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StringRequestAttribute Comparison for `STRING_REQUEST_ATTRIBUTE` attributes, specifically of the **String** type.
type StringRequestAttribute struct {
	BaseComparisonInfo
	Source            *propagation.Source              `json:"source,omitempty"`            // Defines valid sources of request attributes for conditions or placeholders.
	Value             *string                          `json:"value,omitempty"`             // The value to compare to.
	Values            []string                         `json:"values,omitempty"`            // The values to compare to.
	CaseSensitive     *bool                            `json:"caseSensitive,omitempty"`     // The comparison is case-sensitive (`true`) or not case-sensitive (`false`).
	Comparison        StringRequestAttributeComparison `json:"comparison"`                  // Operator of the comparision. You can reverse it by setting **negate** to `true`.
	MatchOnChildCalls *bool                            `json:"matchOnChildCalls,omitempty"` // If `true`, the request attribute is matched on child service calls.   Default is `false`.
	RequestAttribute  string                           `json:"requestAttribute"`            // has no documentation
}

func (me *StringRequestAttribute) GetType() Type {
	return Types.StringRequestAttribute
}

func (me *StringRequestAttribute) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The comparison is case-sensitive (`true`) or not case-sensitive (`false`)",
		},
		"match_on_child_calls": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If `true`, the request attribute is matched on child service calls. Default is `false`",
		},
		"source": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Defines valid sources of request attributes for conditions or placeholders",
			Elem:        &schema.Resource{Schema: new(propagation.Source).Schema()},
		},
		"request_attribute": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "No documentation available for this attribute",
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
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *StringRequestAttribute) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"values":               me.Values,
		"value":                me.Value,
		"operator":             me.Comparison,
		"case_sensitive":       me.CaseSensitive,
		"match_on_child_calls": me.MatchOnChildCalls,
		"request_attribute":    me.RequestAttribute,
		"source":               me.Source,
		"unknowns":             me.Unknowns,
	})
}

func (me *StringRequestAttribute) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"values":               &me.Values,
		"value":                &me.Value,
		"operator":             &me.Comparison,
		"case_sensitive":       &me.CaseSensitive,
		"match_on_child_calls": &me.MatchOnChildCalls,
		"request_attribute":    &me.RequestAttribute,
		"source":               &me.Source,
		"unknowns":             &me.Unknowns,
	})
}

func (me *StringRequestAttribute) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"type":              me.GetType(),
		"negate":            me.Negate,
		"values":            me.Values,
		"value":             me.Value,
		"comparison":        me.Comparison,
		"caseSensitive":     me.CaseSensitive,
		"matchOnChildCalls": me.MatchOnChildCalls,
		"requestAttribute":  me.RequestAttribute,
		"source":            me.Source,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *StringRequestAttribute) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"negate":            &me.Negate,
		"values":            &me.Values,
		"value":             &me.Value,
		"comparison":        &me.Comparison,
		"caseSensitive":     &me.CaseSensitive,
		"matchOnChildCalls": &me.MatchOnChildCalls,
		"requestAttribute":  &me.RequestAttribute,
		"source":            &me.Source,
	})
}

// StringRequestAttributeComparison Operator of the comparision. You can reverse it by setting **negate** to `true`.
type StringRequestAttributeComparison string

// StringRequestAttributeComparisons offers the known enum values
var StringRequestAttributeComparisons = struct {
	BeginsWith      StringRequestAttributeComparison
	BeginsWithAnyOf StringRequestAttributeComparison
	Contains        StringRequestAttributeComparison
	EndsWith        StringRequestAttributeComparison
	EndsWithAnyOf   StringRequestAttributeComparison
	Equals          StringRequestAttributeComparison
	EqualsAnyOf     StringRequestAttributeComparison
	Exists          StringRequestAttributeComparison
	RegexMatches    StringRequestAttributeComparison
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
