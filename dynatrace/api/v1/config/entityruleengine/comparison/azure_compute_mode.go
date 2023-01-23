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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/azure_compute_mode"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AzureComputeMode Comparison for `AZURE_COMPUTE_MODE` attributes.
type AzureComputeMode struct {
	BaseComparison
	Value    *azure_compute_mode.Value   `json:"value,omitempty"` // The value to compare to.
	Operator azure_compute_mode.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
}

func (acmc *AzureComputeMode) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.AzureComputeMode
}

func (acmc *AzureComputeMode) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be AZURE_COMPUTE_MODE",
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
			Description: "The value to compare to. Possible values are DEDICATED or SHARED.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (acmc *AzureComputeMode) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(acmc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", acmc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(acmc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", acmc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (acmc *AzureComputeMode) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), acmc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &acmc.Unknowns); err != nil {
			return err
		}
		delete(acmc.Unknowns, "type")
		delete(acmc.Unknowns, "negate")
		delete(acmc.Unknowns, "operator")
		delete(acmc.Unknowns, "value")
		if len(acmc.Unknowns) == 0 {
			acmc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		acmc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		acmc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		acmc.Operator = azure_compute_mode.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		acmc.Value = azure_compute_mode.Value(value.(string)).Ref()
	}
	return nil
}

func (acmc *AzureComputeMode) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(acmc.Unknowns) > 0 {
		for k, v := range acmc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(acmc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(acmc.GetType())
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&acmc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if acmc.Value != nil {
		rawMessage, err := json.Marshal(acmc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (acmc *AzureComputeMode) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	acmc.Type = acmc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &acmc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &acmc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &acmc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		acmc.Unknowns = m
	}
	return nil
}
