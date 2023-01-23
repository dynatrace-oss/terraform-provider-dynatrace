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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/bitness"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Bitness Comparison for `BITNESS` attributes.
type Bitness struct {
	BaseComparison
	Operator bitness.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *bitness.Value   `json:"value,omitempty"` // The value to compare to.
}

func (bc *Bitness) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.Bitness
}

func (bc *Bitness) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be BITNESS",
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
			Description: "Either EQUALS or EXISTS. You can reverse it by setting **negate** to `true`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to. Possible values are 32 and 64.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (bc *Bitness) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(bc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", bc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(bc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", bc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (bc *Bitness) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), bc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &bc.Unknowns); err != nil {
			return err
		}
		delete(bc.Unknowns, "type")
		delete(bc.Unknowns, "negate")
		delete(bc.Unknowns, "operator")
		delete(bc.Unknowns, "value")
		if len(bc.Unknowns) == 0 {
			bc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		bc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		bc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		bc.Operator = bitness.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		bc.Value = bitness.Value(value.(string)).Ref()
	}
	return nil
}

func (bc *Bitness) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(bc.Unknowns) > 0 {
		for k, v := range bc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(bc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.Bitness)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&bc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if bc.Value != nil {
		rawMessage, err := json.Marshal(bc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (bc *Bitness) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	bc.Type = bc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &bc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &bc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &bc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		bc.Unknowns = m
	}
	return nil
}
