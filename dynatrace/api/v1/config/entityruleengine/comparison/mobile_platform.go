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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison/mobile_platform"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MobilePlatform Comparison for `MOBILE_PLATFORM` attributes.
type MobilePlatform struct {
	BaseComparison
	Operator mobile_platform.Operator `json:"operator"`        // Operator of the comparison. You can reverse it by setting **negate** to `true`.  Possible values depend on the **type** of the comparison. Find the list of actual models in the description of the **type** field and check the description of the model you need.
	Value    *mobile_platform.Value   `json:"value,omitempty"` // The value to compare to.
}

func (mpc *MobilePlatform) GetType() ComparisonBasicType {
	return ComparisonBasicTypes.MobilePlatform
}

func (mpc *MobilePlatform) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be MOBILE_PLATFORM",
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
			Description: "The value to compare to. Possible values are ANDROID, IOS, LINUX, MAC_OS, OTHER, TVOS and WINDOWS.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (mpc *MobilePlatform) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(mpc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", mpc.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(mpc.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", mpc.Value.String()); err != nil {
		return err
	}
	return nil
}

func (mpc *MobilePlatform) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), mpc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &mpc.Unknowns); err != nil {
			return err
		}
		delete(mpc.Unknowns, "type")
		delete(mpc.Unknowns, "negate")
		delete(mpc.Unknowns, "operator")
		delete(mpc.Unknowns, "value")
		if len(mpc.Unknowns) == 0 {
			mpc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		mpc.Type = ComparisonBasicType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		mpc.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		mpc.Operator = mobile_platform.Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		mpc.Value = mobile_platform.Value(value.(string)).Ref()
	}
	return nil
}

func (mpc *MobilePlatform) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(mpc.Unknowns) > 0 {
		for k, v := range mpc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(mpc.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ComparisonBasicTypes.MobilePlatform)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(&mpc.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	if mpc.Value != nil {
		rawMessage, err := json.Marshal(mpc.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (mpc *MobilePlatform) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	mpc.Type = mpc.GetType()
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &mpc.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &mpc.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &mpc.Value); err != nil {
			return err
		}
	}
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "type")
	if len(m) > 0 {
		mpc.Unknowns = m
	}
	return nil
}
