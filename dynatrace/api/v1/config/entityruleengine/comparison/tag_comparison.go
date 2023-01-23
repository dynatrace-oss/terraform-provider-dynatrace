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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/tag"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Tag Comparison for `TAG` attributes.
type Tag struct {
	BaseComparison
	Value    *tag.Info    `json:"value,omitempty"` // Tag of a Dynatrace entity.
	Operator tag.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
}

func (tc *Tag) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.Tag
}

func (tc *Tag) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be TAG",
			Optional:    true,
			Deprecated:  "The value of the attribute type is implicit, therefore shouldn't get specified",
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "Reverses the operator. For example it turns the **begins with** into **does not begin with**",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Operator of the comparison. Possible values are EQUALS and TAG_KEY_EQUALS. You can reverse it by setting **negate** to `true`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Description: "Tag of a Dynatrace entity",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(tag.Info).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (tc *Tag) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(tc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", tc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(tc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", tc.Value); err != nil {
		return err
	}
	return nil
}

func (tc *Tag) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), tc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &tc.Unknowns); err != nil {
			return err
		}
		delete(tc.Unknowns, "type")
		delete(tc.Unknowns, "negate")
		delete(tc.Unknowns, "operator")
		delete(tc.Unknowns, "value")
		if len(tc.Unknowns) == 0 {
			tc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		tc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		tc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		tc.Operator = tag.Operator(value.(string))
	}
	if _, ok := decoder.GetOk("value.#"); ok {
		tc.Value = new(tag.Info)
		if err := tc.Value.UnmarshalHCL(hcl.NewDecoder(decoder, "value", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (tc *Tag) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(tc.Unknowns) > 0 {
		for k, v := range tc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(tc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.Tag)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&tc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if tc.Value != nil {
		rawMessage, err := json.Marshal(tc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (tc *Tag) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	tc.Type = tc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &tc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &tc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &tc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		tc.Unknowns = m
	}
	return nil
}
