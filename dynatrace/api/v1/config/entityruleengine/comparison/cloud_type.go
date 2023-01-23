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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/cloud_type"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CloudType Comparison for `CLOUD_TYPE` attributes.
type CloudType struct {
	BaseComparison
	Operator cloud_type.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *cloud_type.Value   `json:"value,omitempty"` // The value to compare to.
}

func (ctc *CloudType) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.CloudType
}

func (ctc *CloudType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be CLOUD_TYPE",
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
			Description: "The value to compare to. Possible values are AZURE, EC2, GOOGLE_CLOUD_PLATFORM, OPENSTACK, ORACLE and UNRECOGNIZED.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (ctc *CloudType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(ctc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", ctc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(ctc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", ctc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (ctc *CloudType) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), ctc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &ctc.Unknowns); err != nil {
			return err
		}
		delete(ctc.Unknowns, "type")
		delete(ctc.Unknowns, "negate")
		delete(ctc.Unknowns, "operator")
		delete(ctc.Unknowns, "value")
		if len(ctc.Unknowns) == 0 {
			ctc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		ctc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		ctc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		ctc.Operator = cloud_type.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		ctc.Value = cloud_type.Value(value.(string)).Ref()
	}
	return nil
}

func (ctc *CloudType) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(ctc.Unknowns) > 0 {
		for k, v := range ctc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(ctc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ctc.GetType())
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&ctc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if ctc.Value != nil {
		rawMessage, err := json.Marshal(ctc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (ctc *CloudType) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	ctc.Type = ctc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &ctc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &ctc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &ctc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		ctc.Unknowns = m
	}
	return nil
}
