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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/application_type"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ApplicationType Comparison for `APPLICATION_TYPE` attributes.
type ApplicationType struct {
	BaseComparison
	Operator application_type.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *application_type.Value   `json:"value,omitempty"` // The value to compare to.
}

func (atc *ApplicationType) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.ApplicationType
}

func (atc *ApplicationType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be APPLICATION_TYPE",
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
			Description: "Operator of the comparison. You can reverse it by setting **negate** to `true`",
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

func (atc *ApplicationType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(atc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", atc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(atc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", atc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (atc *ApplicationType) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), atc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &atc.Unknowns); err != nil {
			return err
		}
		delete(atc.Unknowns, "type")
		delete(atc.Unknowns, "negate")
		delete(atc.Unknowns, "operator")
		delete(atc.Unknowns, "value")
		if len(atc.Unknowns) == 0 {
			atc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		atc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		atc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		atc.Operator = application_type.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		atc.Value = application_type.Value(value.(string)).Ref()
	}
	return nil
}

func (atc *ApplicationType) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(atc.Unknowns) > 0 {
		for k, v := range atc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(atc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.ApplicationType)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&atc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if atc.Value != nil {
		rawMessage, err := json.Marshal(atc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (atc *ApplicationType) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	atc.Type = atc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &atc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &atc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &atc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		atc.Unknowns = m
	}
	return nil
}
