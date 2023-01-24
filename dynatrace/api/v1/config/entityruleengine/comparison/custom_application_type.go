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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/custom_application_type"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomApplicationType Comparison for `CUSTOM_APPLICATION_TYPE` attributes.
type CustomApplicationType struct {
	BaseComparison
	Operator custom_application_type.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *custom_application_type.Value   `json:"value,omitempty"` // The value to compare to.
}

func (catc *CustomApplicationType) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.CustomApplicationType
}

func (catc *CustomApplicationType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be CUSTOM_APPLICATION_TYPE",
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
			Description: "The value to compare to. Possible values are AMAZON_ECHO, DESKTOP, EMBEDDED, IOT, MICROSOFT_HOLOLENS and UFO.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (catc *CustomApplicationType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(catc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", catc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(catc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", catc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (catc *CustomApplicationType) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), catc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &catc.Unknowns); err != nil {
			return err
		}
		delete(catc.Unknowns, "type")
		delete(catc.Unknowns, "negate")
		delete(catc.Unknowns, "operator")
		delete(catc.Unknowns, "value")
		if len(catc.Unknowns) == 0 {
			catc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		catc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		catc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		catc.Operator = custom_application_type.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		catc.Value = custom_application_type.Value(value.(string)).Ref()
	}
	return nil
}

func (catc *CustomApplicationType) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(catc.Unknowns) > 0 {
		for k, v := range catc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(catc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.CustomApplicationType)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&catc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if catc.Value != nil {
		rawMessage, err := json.Marshal(catc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (catc *CustomApplicationType) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	catc.Type = catc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &catc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &catc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &catc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		catc.Unknowns = m
	}
	return nil
}
