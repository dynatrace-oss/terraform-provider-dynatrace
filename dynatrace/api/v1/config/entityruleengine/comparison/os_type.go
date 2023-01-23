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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/ostype"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// OSType Comparison for `OS_TYPE` attributes.
type OSType struct {
	BaseComparison
	Operator ostype.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *ostype.Value   `json:"value,omitempty"` // The value to compare to.
}

func (otc *OSType) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.OSType
}

func (otc *OSType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be OS_TYPE",
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
			Description: "Operator of the comparison. Possible values are EQUALS and EXISTS. You can reverse it by setting **negate** to `true`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to. Possible values are AIX, DARWIN, HPUX, LINUX, SOLARIS, WINDOWS and ZOS.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (otc *OSType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(otc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", otc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(otc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", otc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (otc *OSType) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), otc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &otc.Unknowns); err != nil {
			return err
		}
		delete(otc.Unknowns, "type")
		delete(otc.Unknowns, "negate")
		delete(otc.Unknowns, "operator")
		delete(otc.Unknowns, "value")
		if len(otc.Unknowns) == 0 {
			otc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		otc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		otc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		otc.Operator = ostype.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		otc.Value = ostype.Value(value.(string)).Ref()
	}
	return nil
}

func (otc *OSType) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(otc.Unknowns) > 0 {
		for k, v := range otc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(otc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(otc.GetType())
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&otc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if otc.Value != nil {
		rawMessage, err := json.Marshal(otc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (otc *OSType) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	otc.Type = otc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &otc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &otc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &otc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		otc.Unknowns = m
	}
	return nil
}
