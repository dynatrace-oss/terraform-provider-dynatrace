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

package managementzones

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DimensionalRule represents the dimensional rule of the management zone usage.
// It defines how the management zone applies.
// Each rule is evaluated independently of all other rules
type DimensionalRule struct {
	Enabled    *bool                       `json:"enabled"`    // The rule is enabled (`true`) or disabled (`false`)
	AppliesTo  Application                 `json:"appliesTo"`  // The target of the rule
	Conditions []*DimensionalRuleCondition `json:"conditions"` // A list of conditions for the management zone. \n\n The management zone applies only if **all** conditions are fulfilled
	Unknowns   map[string]json.RawMessage  `json:"-"`
}

func (me *DimensionalRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The rule is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"applies_to": {
			Type:        schema.TypeString,
			Description: "The target of the rule. Possible values are\n   - `ANY`\n   - `LOG`\n   - `METRIC`",
			Required:    true,
		},
		"condition": {
			Type:        schema.TypeList,
			Description: "A list of conditions for the management zone. The management zone applies only if **all** conditions are fulfilled",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(DimensionalRuleCondition).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *DimensionalRule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("enabled", opt.Bool(me.Enabled)); err != nil {
		return err
	}
	if err := properties.Encode("applies_to", string(me.AppliesTo)); err != nil {
		return err
	}
	if err := properties.Encode("condition", me.Conditions); err != nil {
		return err
	}

	return nil
}

func (me *DimensionalRule) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "applies_to")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "condition")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("applies_to"); ok {
		me.AppliesTo = Application(value.(string))
	}
	if result, ok := decoder.GetOk("condition.#"); ok {
		me.Conditions = []*DimensionalRuleCondition{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(DimensionalRuleCondition)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "condition", idx)); err != nil {
				return err
			}
			me.Conditions = append(me.Conditions, entry)
		}
	}
	return nil
}

func (me *DimensionalRule) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if opt.Bool(me.Enabled) {
		rawMessage, err := json.Marshal(me.Enabled)
		if err != nil {
			return nil, err
		}
		m["enabled"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.AppliesTo)
		if err != nil {
			return nil, err
		}
		m["appliesTo"] = rawMessage
	}
	if len(me.Conditions) > 0 {
		rawMessage, err := json.Marshal(me.Conditions)
		if err != nil {
			return nil, err
		}
		m["conditions"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *DimensionalRule) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["enabled"]; found {
		if err := json.Unmarshal(v, &me.Enabled); err != nil {
			return err
		}
	}
	if v, found := m["appliesTo"]; found {
		if err := json.Unmarshal(v, &me.AppliesTo); err != nil {
			return err
		}
	}
	if v, found := m["conditions"]; found {
		if err := json.Unmarshal(v, &me.Conditions); err != nil {
			return err
		}
	}
	delete(m, "appliesTo")
	delete(m, "enabled")
	delete(m, "conditions")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// Application The value to compare to.
type Application string

func (v Application) Ref() *Application {
	return &v
}

// Applications offers the known enum values
var Applications = struct {
	Any    Application
	Log    Application
	Metric Application
}{
	"ANY",
	"LOG",
	"METRIC",
}

func (v *Application) String() string {
	return string(*v)
}
