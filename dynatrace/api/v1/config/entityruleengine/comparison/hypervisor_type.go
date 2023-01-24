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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/hypervisor_type"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HypervisorType Comparison for `HYPERVISOR_TYPE` attributes.
type HypervisorType struct {
	BaseComparison
	Operator hypervisor_type.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *hypervisor_type.Value   `json:"value,omitempty"` // The value to compare to.
}

func (htc *HypervisorType) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.HypervisorType
}

func (htc *HypervisorType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be HYPERVISOR_TYPE",
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
			Description: "The value to compare to. Possible values are AHV, HYPER_V, KVM, LPAR, QEMU, VIRTUAL_BOX, VMWARE, WPAR and XEN.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (htc *HypervisorType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(htc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", htc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(htc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", htc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (htc *HypervisorType) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), htc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &htc.Unknowns); err != nil {
			return err
		}
		delete(htc.Unknowns, "type")
		delete(htc.Unknowns, "negate")
		delete(htc.Unknowns, "operator")
		delete(htc.Unknowns, "value")
		if len(htc.Unknowns) == 0 {
			htc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		htc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		htc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		htc.Operator = hypervisor_type.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		htc.Value = hypervisor_type.Value(value.(string)).Ref()
	}
	return nil
}

func (htc *HypervisorType) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(htc.Unknowns) > 0 {
		for k, v := range htc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(htc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(htc.GetType())
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&htc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if htc.Value != nil {
		rawMessage, err := json.Marshal(htc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (htc *HypervisorType) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	htc.Type = htc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &htc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &htc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &htc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		htc.Unknowns = m
	}
	return nil
}
