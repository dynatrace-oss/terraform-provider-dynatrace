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

package comparison

import (
	"encoding/json"

	"github.com/dlclark/regexp2"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/stringc"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// String Comparison for `STRING` attributes.
type String struct {
	BaseComparison
	CaseSensitive bool             `json:"caseSensitive,omitempty"` // The comparison is case-sensitive (`true`) or insensitive (`false`).
	Operator      stringc.Operator `json:"operator"`                // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value         *string          `json:"value,omitempty"`         // The value to compare to.
}

func (sc *String) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.String
}

func (sc *String) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be STRING",
			Optional:    true,
			Deprecated:  "The value of the attribute type is implicit, therefore shouldn't get specified",
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "Reverses the operator. For example it turns the **begins with** into **does not begin with**",
			Optional:    true,
		},
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "The comparison is case-sensitive (`true`) or insensitive (`false`)",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Operator of the comparison. Possible values are BEGINS_WITH, CONTAINS, ENDS_WITH, EQUALS, EXISTS and REGEX_MATCHES. You can reverse it by setting **negate** to `true`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (sc *String) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(sc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", sc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(sc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("case_sensitive", sc.CaseSensitive); err != nil {
		return err
	}
	if err := properties.Encode("value", sc.Value); err != nil {
		return err
	}
	return nil
}

func (sc *String) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), sc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &sc.Unknowns); err != nil {
			return err
		}
		delete(sc.Unknowns, "type")
		delete(sc.Unknowns, "negate")
		delete(sc.Unknowns, "operator")
		delete(sc.Unknowns, "value")
		if len(sc.Unknowns) == 0 {
			sc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		sc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		sc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		sc.Operator = stringc.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		sc.Value = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("case_sensitive"); ok {
		sc.CaseSensitive = value.(bool)
	}

	return nil
}

func (sc *String) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(sc.Unknowns) > 0 {
		for k, v := range sc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(sc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.String)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&sc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if sc.Value != nil {
		rawMessage, err := json.Marshal(sc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	if sc.Operator != stringc.Operators.Exists {
		rawMessage, err := json.Marshal(sc.CaseSensitive)
		if err != nil {
			return nil, err
		}
		m["caseSensitive"] = rawMessage
	}
	return json.Marshal(m)
}

func (sc *String) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	sc.Type = sc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &sc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &sc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &sc.Value); err != nil {
			return err
		}
	}
	if v, found := m["caseSensitive"]; found {
		if err := json.Unmarshal(v, &sc.CaseSensitive); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	delete(m, "caseSensitive")
	if len(m) > 0 {
		sc.Unknowns = m
	}
	return nil
}

func (sc *String) Validate() []string {
	if sc.Operator == stringc.Operators.RegexMatches && sc.Value != nil {
		r, _ := regexCheck(*sc.Value)
		return r
	}
	return []string{}
}

// Please use possessive or lazy quantifiers within your capture group
const QUANTIFIERS_WITHIN_CAPTURING_GROUP = `(?<!\\)(?:\\\\)*([^\*\+]*[\+\*](?![\+\?]).*(?<!\\)(?:\\\\)*)`

// Please do not use a greedy all match
const ALL_MATCH_NO_CAPTURING_GROUP = `(?<!\()\.[\*\+](?![\)\+\?])`

// Greedy or lazy character classes are not allowed, please use a possessive quantifier instead
const CHARACTER_CLASS_WITH_GREEDY_OR_LAZY_QUANTIFIER = `(?<!\\)(?:\\\\)*](?>\*|\+)(?!\+)`

func regexCheck(s string) ([]string, error) {
	result := []string{}
	match := false
	var err error

	if match, err = regexp2.MustCompile(QUANTIFIERS_WITHIN_CAPTURING_GROUP, 0).MatchString(s); err != nil {
		return nil, err
	} else if match {
		result = append(result, "FLAWED SETTINGS Please use possessive or lazy quantifiers within your capture group")
	}
	if match, err = regexp2.MustCompile(ALL_MATCH_NO_CAPTURING_GROUP, 0).MatchString(s); err != nil {
		return nil, err
	} else if match {
		result = append(result, "FLAWED SETTINGS Please do not use a greedy all match")
	}
	if match, err = regexp2.MustCompile(CHARACTER_CLASS_WITH_GREEDY_OR_LAZY_QUANTIFIER, 0).MatchString(s); err != nil {
		return nil, err
	} else if match {
		result = append(result, "FLAWED SETTINGS Greedy or lazy character classes are not allowed, please use a possessive quantifier instead")
	}
	return result, nil
}
