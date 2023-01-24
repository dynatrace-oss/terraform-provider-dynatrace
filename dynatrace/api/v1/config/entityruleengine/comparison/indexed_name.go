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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/indexed_name"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IndexedName Comparison for `INDEXED_NAME` attributes.
type IndexedName struct {
	BaseComparison
	Operator indexed_name.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *string               `json:"value,omitempty"` // The value to compare to.
}

func (inc *IndexedName) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.IndexedName
}

func (inc *IndexedName) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be INDEXED_NAME",
			Optional:    true,
			Deprecated:  "The value of the attribute type is implicit, therefore shouldn't get specified",
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "Reverses the operator. For example it turns EQUALS into DOES NOT EQUAL",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Either EQUALS, CONTAINS or EXISTS. You can reverse it by setting **negate** to `true`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (inc *IndexedName) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(inc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", inc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(inc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", inc.Value); err != nil {
		return err
	}
	return nil
}

func (inc *IndexedName) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), inc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &inc.Unknowns); err != nil {
			return err
		}
		delete(inc.Unknowns, "type")
		delete(inc.Unknowns, "negate")
		delete(inc.Unknowns, "operator")
		delete(inc.Unknowns, "value")
		if len(inc.Unknowns) == 0 {
			inc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		inc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		inc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		inc.Operator = indexed_name.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		inc.Value = opt.NewString(value.(string))
	}
	return nil
}

func (inc *IndexedName) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(inc.Unknowns) > 0 {
		for k, v := range inc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(inc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(inc.GetType())
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&inc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if inc.Value != nil {
		rawMessage, err := json.Marshal(inc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (inc *IndexedName) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	inc.Type = inc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &inc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &inc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &inc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		inc.Unknowns = m
	}
	return nil
}
