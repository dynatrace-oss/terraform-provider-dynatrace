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
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (me *ManagementZone) Deprecated() string {
	return "This resource API endpoint has been deprecated, please use `dynatrace_management_zone_v2` instead."
}

// ManagementZone The configuration of the management zone. It defines how the management zone applies.
type ManagementZone struct {
	Description              *string                    `json:"description,omitempty"`              // The description of the management zone
	Name                     string                     `json:"name"`                               // The name of the management zone.
	Rules                    []*Rule                    `json:"rules,omitempty"`                    // A list of rules for management zone usage. Each rule is evaluated independently of all other rules.
	DimensionalRules         []*DimensionalRule         `json:"dimensionalRules,omitempty"`         // A list of dimensional data rules for management zone usage. If several rules are specified, the **OR** logic applies
	EntitySelectorBasedRules []*EntitySelectorBasedRule `json:"entitySelectorBasedRules,omitempty"` // A list of entity-selector based rules for management zone usage. If several rules are specified, the **OR** logic applies
	Unknowns                 map[string]json.RawMessage `json:"-"`
}

func (mz *ManagementZone) Validate() []string {
	result := []string{}
	if len(mz.Rules) > 0 {
		for _, rule := range mz.Rules {
			result = append(result, rule.Validate()...)
		}
	}
	return result
}

func (mz *ManagementZone) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the management zone",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the management zone",
			Optional:    true,
		},
		"rules": {
			Type:        schema.TypeSet,
			Description: "A list of rules for management zone usage.  Each rule is evaluated independently of all other rules",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Rule).Schema()},
		},
		"dimensional_rule": {
			Type:        schema.TypeSet,
			Description: "A list of dimensional data rules for management zone usage. If several rules are specified, the `or` logic applies",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(DimensionalRule).Schema()},
		},
		"entity_selector_based_rule": {
			Type:        schema.TypeSet,
			Description: "A list of entity-selector based rules for management zone usage. If several rules are specified, the `or` logic applies",
			Optional:    true,
			MinItems:    1,
			Set: func(i interface{}) int {
				if m, mok := i.(map[string]any); mok {
					for k, v := range m {
						if vs, ok := v.(string); ok {
							m[k] = strings.TrimSpace(vs)
						}
					}
					dd, _ := json.Marshal(m)
					return hashcode.String(string(dd))
				}
				dd, _ := json.Marshal(i)
				return hashcode.String(string(dd))
			},
			Elem: &schema.Resource{Schema: new(EntitySelectorBasedRule).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (mz *ManagementZone) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(mz.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", mz.Name); err != nil {
		return err
	}
	if mz.Description != nil && len(*mz.Description) > 0 {
		if err := properties.Encode("description", mz.Description); err != nil {
			return err
		}
	}
	if err := properties.Encode("rules", mz.Rules); err != nil {
		return err
	}
	if err := properties.Encode("dimensional_rule", mz.DimensionalRules); err != nil {
		return err
	}
	if err := properties.Encode("entity_selector_based_rule", mz.EntitySelectorBasedRules); err != nil {
		return err
	}
	return nil
}

func (mz *ManagementZone) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), mz); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &mz.Unknowns); err != nil {
			return err
		}
		delete(mz.Unknowns, "rules")
		delete(mz.Unknowns, "dimensional_rule")
		delete(mz.Unknowns, "entity_selector_based_rule")
		delete(mz.Unknowns, "metadata")
		delete(mz.Unknowns, "name")
		delete(mz.Unknowns, "description")
		if len(mz.Unknowns) == 0 {
			mz.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		mz.Name = value.(string)
	}
	if value, ok := decoder.GetOk("description"); ok {
		mz.Description = opt.NewString(value.(string))
	}
	if err := decoder.DecodeSlice("rules", &mz.Rules); err != nil {
		return err
	}

	if err := decoder.DecodeSlice("dimensional_rule", &mz.DimensionalRules); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("entity_selector_based_rule", &mz.EntitySelectorBasedRules); err != nil {
		return err
	}
	if len(mz.EntitySelectorBasedRules) > 0 {
		var rules []*EntitySelectorBasedRule
		for _, rule := range mz.EntitySelectorBasedRules {
			if (rule.Enabled == nil || *rule.Enabled == false) && len(rule.Selector) == 0 {
				continue
			}
			rules = append(rules, rule)
		}
		mz.EntitySelectorBasedRules = rules
	}
	return nil
}

func (mz *ManagementZone) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(mz.Unknowns) > 0 {
		for k, v := range mz.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(mz.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	if mz.Description != nil {
		rawMessage, err := json.Marshal(mz.Description)
		if err != nil {
			return nil, err
		}
		m["description"] = rawMessage
	}
	if len(mz.Rules) > 0 {
		rawMessage, err := json.Marshal(mz.Rules)
		if err != nil {
			return nil, err
		}
		m["rules"] = rawMessage
	}
	if len(mz.DimensionalRules) > 0 {
		rawMessage, err := json.Marshal(mz.DimensionalRules)
		if err != nil {
			return nil, err
		}
		m["dimensionalRules"] = rawMessage
	}
	if len(mz.EntitySelectorBasedRules) > 0 {
		rawMessage, err := json.Marshal(mz.EntitySelectorBasedRules)
		if err != nil {
			return nil, err
		}
		m["entitySelectorBasedRules"] = rawMessage
	}
	return json.Marshal(m)
}

func (mz *ManagementZone) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["description"]; found {
		if err := json.Unmarshal(v, &mz.Description); err != nil {
			return err
		}
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &mz.Name); err != nil {
			return err
		}
	}
	if v, found := m["rules"]; found {
		if err := json.Unmarshal(v, &mz.Rules); err != nil {
			return err
		}
	}
	if v, found := m["dimensionalRules"]; found {
		if err := json.Unmarshal(v, &mz.DimensionalRules); err != nil {
			return err
		}
	}
	if v, found := m["entitySelectorBasedRules"]; found {
		if err := json.Unmarshal(v, &mz.EntitySelectorBasedRules); err != nil {
			return err
		}
	}
	delete(m, "id")
	delete(m, "name")
	delete(m, "description")
	delete(m, "metadata")
	delete(m, "rules")
	delete(m, "dimensionalRules")
	delete(m, "entitySelectorBasedRules")

	if len(m) > 0 {
		mz.Unknowns = m
	}
	return nil
}
