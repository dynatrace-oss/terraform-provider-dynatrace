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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/integer"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Integer Comparison for `INTEGER` attributes.
type Integer struct {
	BaseComparison
	Operator integer.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *int32           `json:"value,omitempty"` // The value to compare to.
}

func (ic *Integer) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.Integer
}

func (ic *Integer) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be INTEGER",
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
			Description: "Operator of the comparison. Possible values are EQUALS, EXISTS, GREATER_THAN, GREATER_THAN_OR_EQUAL, LOWER_THAN and LOWER_THAN_OR_EQUAL. You can reverse it by setting **negate** to `true`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeInt,
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

func (ic *Integer) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(ic.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", ic.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(ic.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", int(opt.Int32(ic.Value))); err != nil {
		return err
	}
	return nil
}

func (ic *Integer) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), ic); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &ic.Unknowns); err != nil {
			return err
		}
		delete(ic.Unknowns, "type")
		delete(ic.Unknowns, "negate")
		delete(ic.Unknowns, "operator")
		delete(ic.Unknowns, "value")
		if len(ic.Unknowns) == 0 {
			ic.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		ic.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		ic.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		ic.Operator = integer.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		ic.Value = opt.NewInt32(int32(value.(int)))
	}
	return nil
}

func (ic *Integer) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(ic.Unknowns) > 0 {
		for k, v := range ic.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(ic.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ic.GetType())
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&ic.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if ic.Value != nil {
		rawMessage, err := json.Marshal(ic.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (ic *Integer) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	ic.Type = ic.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &ic.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &ic.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &ic.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		ic.Unknowns = m
	}
	return nil
}
