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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/database_topology"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DatabaseTopology Comparison for `DATABASE_TOPOLOGY` attributes.
type DatabaseTopology struct {
	BaseComparison
	Operator database_topology.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *database_topology.Value   `json:"value,omitempty"` // The value to compare to.
}

func (dtc *DatabaseTopology) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.DatabaseTopology
}

func (dtc *DatabaseTopology) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be DATABASE_TOPOLOGY",
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
			Description: "The value to compare to. Possible values are CLUSTER, EMBEDDED, FAILOVER, IPC, LOAD_BALANCING, SINGLE_SERVER and UNSPECIFIED.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (dtc *DatabaseTopology) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(dtc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", dtc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(dtc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", dtc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (dtc *DatabaseTopology) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), dtc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &dtc.Unknowns); err != nil {
			return err
		}
		delete(dtc.Unknowns, "type")
		delete(dtc.Unknowns, "negate")
		delete(dtc.Unknowns, "operator")
		delete(dtc.Unknowns, "value")
		if len(dtc.Unknowns) == 0 {
			dtc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		dtc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		dtc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		dtc.Operator = database_topology.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		dtc.Value = database_topology.Value(value.(string)).Ref()
	}
	return nil
}

func (dtc *DatabaseTopology) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(dtc.Unknowns) > 0 {
		for k, v := range dtc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(dtc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.DatabaseTopology)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&dtc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if dtc.Value != nil {
		rawMessage, err := json.Marshal(dtc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (dtc *DatabaseTopology) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	dtc.Type = dtc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &dtc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &dtc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &dtc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		dtc.Unknowns = m
	}
	return nil
}
