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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/service_type"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ServiceType Comparison for `SERVICE_TYPE` attributes.
type ServiceType struct {
	BaseComparison
	Operator service_type.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *service_type.Value   `json:"value,omitempty"` // The value to compare to.
}

func (stc *ServiceType) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.ServiceType
}

func (stc *ServiceType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be SERVICE_TYPE",
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
			Description: "The value to compare to. Possible values are BACKGROUND_ACTIVITY, CICS_SERVICE, CUSTOM_SERVICE, DATABASE_SERVICE, ENTERPRISE_SERVICE_BUS_SERVICE, EXTERNAL, IBM_INTEGRATION_BUS_SERVICE, IMS_SERVICE, MESSAGING_SERVICE, QUEUE_LISTENER_SERVICE, RMI_SERVICE, RPC_SERVICE, WEB_REQUEST_SERVICE and WEB_SERVICE.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (stc *ServiceType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(stc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", stc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(stc.Operator)); err != nil {
		return err
	}
	if stc.Value != nil {
		if err := properties.Encode("value", stc.Value.String()); err != nil {
			return err
		}
	}
	return nil
}

func (stc *ServiceType) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), stc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &stc.Unknowns); err != nil {
			return err
		}
		delete(stc.Unknowns, "type")
		delete(stc.Unknowns, "negate")
		delete(stc.Unknowns, "operator")
		delete(stc.Unknowns, "value")
		if len(stc.Unknowns) == 0 {
			stc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		stc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		stc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		stc.Operator = service_type.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		stc.Value = service_type.Value(value.(string)).Ref()
	}
	return nil
}

func (stc *ServiceType) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(stc.Unknowns) > 0 {
		for k, v := range stc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(stc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(stc.GetType())
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&stc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if stc.Value != nil {
		rawMessage, err := json.Marshal(stc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (stc *ServiceType) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	stc.Type = stc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &stc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &stc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &stc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		stc.Unknowns = m
	}
	return nil
}
