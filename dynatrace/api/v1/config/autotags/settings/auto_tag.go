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

package autotags

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AutoTag Configuration of an auto-tag. It defines the conditions of tag usage and the tag value.
type AutoTag struct {
	// ID                       *string                    `json:"id,omitempty"`                       // The ID of the auto-tag.
	Name                     string                     `json:"name"`                               // The name of the auto-tag, which is applied to entities.  Additionally you can specify a **valueFormat** in the tag rule. In that case the tag is used in the `name:valueFormat` format.  For example you can extend the `Infrastructure` tag to `Infrastructure:Windows` and `Infrastructure:Linux`.
	Description              *string                    `json:"description"`                        // The description of the auto-tag
	Rules                    []*Rule                    `json:"rules,omitempty"`                    // The list of rules for tag usage. When there are multiple rules, the OR logic applies.
	EntitySelectorBasedRules []*EntitySelectorBasedRule `json:"entitySelectorBasedRules,omitempty"` // A list of entity-selector based rules for auto tagging usage.\n\nIf several rules are specified, the **OR** logic applies
	Unknowns                 map[string]json.RawMessage `json:"-"`
}

func (me *AutoTag) Validate() []string {
	result := []string{}
	if len(me.Rules) > 0 {
		for _, rule := range me.Rules {
			result = append(result, rule.Validate()...)
		}
	}
	return result
}

func (me *AutoTag) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the auto-tag, which is applied to entities.  Additionally you can specify a **valueFormat** in the tag rule. In that case the tag is used in the `name:valueFormat` format.  For example you can extend the `Infrastructure` tag to `Infrastructure:Windows` and `Infrastructure:Linux`.",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the auto-tag.",
			Optional:    true,
		},
		// "metadata": {
		// 	Type:        schema.TypeList,
		// 	MaxItems:    1,
		// 	Description: "`metadata` exists for backwards compatibility but shouldn't get specified anymore",
		// 	Deprecated:  "`metadata` exists for backwards compatibility but shouldn't get specified anymore",
		// 	Optional:    true,
		// 	Elem: &schema.Resource{
		// 		Schema: new(api.ConfigMetadata).Schema(),
		// 	},
		// },
		"rules": {
			Type:        schema.TypeSet,
			Description: "A list of rules for management zone usage.  Each rule is evaluated independently of all other rules",
			Optional:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: new(Rule).Schema(),
			},
		},
		"entity_selector_based_rule": {
			Type:        schema.TypeSet,
			Description: "A list of entity-selector based rules for management zone usage. If several rules are specified, the `or` logic applies",
			Optional:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: new(EntitySelectorBasedRule).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *AutoTag) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("description", me.Description); err != nil {
		return err
	}
	if err := properties.Encode("rules", me.Rules); err != nil {
		return err
	}
	if err := properties.Encode("entity_selector_based_rule", me.EntitySelectorBasedRules); err != nil {
		return err
	}
	return nil
}

func (me *AutoTag) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "rules")
		delete(me.Unknowns, "entity_selector_based_rule")
		delete(me.Unknowns, "metadata")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "description")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("description"); ok {
		me.Description = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("rules"); ok {
		ruleSet := value.(*schema.Set)

		for _, ruleMap := range ruleSet.List() {
			hash := ruleSet.F(ruleMap)
			entry := new(Rule)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "rules", hash)); err != nil {
				return err
			}
			me.Rules = append(me.Rules, entry)
		}
	}
	if value, ok := decoder.GetOk("entity_selector_based_rule"); ok {
		ruleSet := value.(*schema.Set)

		for _, ruleMap := range ruleSet.List() {
			hash := ruleSet.F(ruleMap)
			entry := new(EntitySelectorBasedRule)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "entity_selector_based_rule", hash)); err != nil {
				return err
			}
			me.EntitySelectorBasedRules = append(me.EntitySelectorBasedRules, entry)
		}
	}
	return nil
}

func (me *AutoTag) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Description)
		if err != nil {
			return nil, err
		}
		m["description"] = rawMessage
	}
	if len(me.Rules) > 0 {
		rawMessage, err := json.Marshal(me.Rules)
		if err != nil {
			return nil, err
		}
		m["rules"] = rawMessage
	}
	if len(me.EntitySelectorBasedRules) > 0 {
		rawMessage, err := json.Marshal(me.EntitySelectorBasedRules)
		if err != nil {
			return nil, err
		}
		m["entitySelectorBasedRules"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *AutoTag) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &me.Name); err != nil {
			return err
		}
	}
	if v, found := m["description"]; found {
		if err := json.Unmarshal(v, &me.Description); err != nil {
			return err
		}
	}
	if v, found := m["rules"]; found {
		if err := json.Unmarshal(v, &me.Rules); err != nil {
			return err
		}
	}
	if v, found := m["entitySelectorBasedRules"]; found {
		if err := json.Unmarshal(v, &me.EntitySelectorBasedRules); err != nil {
			return err
		}
	}
	delete(m, "name")
	delete(m, "metadata")
	delete(m, "rules")
	delete(m, "id")
	delete(m, "entitySelectorBasedRules")
	delete(m, "description")
	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
