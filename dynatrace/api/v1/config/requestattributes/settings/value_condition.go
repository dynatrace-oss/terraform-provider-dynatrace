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

package requestattributes

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValueCondition IBM integration bus label node name condition for which the value is captured.
type ValueCondition struct {
	Negate   bool                       `json:"negate"`   // Negate the comparison.
	Operator Operator                   `json:"operator"` // Operator comparing the extracted value to the comparison value.
	Value    string                     `json:"value"`    // The value to compare to.
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *ValueCondition) IsZero() bool {
	if me.Negate {
		return false
	}
	if len(me.Operator) > 0 {
		return false
	}
	if len(me.Value) > 0 {
		return false
	}
	if len(me.Unknowns) > 0 {
		return false
	}
	return true
}

func (me *ValueCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"negate": {
			Type:        schema.TypeBool,
			Description: "Negate the comparison",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Operator comparing the extracted value to the comparison value",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *ValueCondition) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("negate", me.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(me.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}
	return nil
}

func (me *ValueCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "negate")
		delete(me.Unknowns, "operator")
		delete(me.Unknowns, "value")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("negate"); ok {
		me.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		me.Operator = Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = value.(string)
	}
	return nil
}

func (me *ValueCondition) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("negate", me.Negate); err != nil {
		return nil, err
	}
	if err := m.Marshal("operator", me.Operator); err != nil {
		return nil, err
	}
	if err := m.Marshal("value", me.Value); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *ValueCondition) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("negate", &me.Negate); err != nil {
		return err
	}
	if err := m.Unmarshal("operator", &me.Operator); err != nil {
		return err
	}
	if err := m.Unmarshal("value", &me.Value); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
