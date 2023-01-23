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

package naming

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rules []*Rule

func (me *Rules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Description: "The settings of naming rule",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Rule).Schema()},
		},
	}
}

func (me Rules) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("rule", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *Rules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

// Rule The settings of naming rule
type Rule struct {
	Template        string         `json:"template"`             // Naming pattern. Use Curly brackets `{}` to select placeholders
	Conditions      RuleConditions `json:"conditions,omitempty"` // Defines the conditions when the naming rule should apply
	UseOrConditions bool           `json:"useOrConditions"`      // If set to `true` the conditions will be connected by logical OR instead of logical AND
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Description: "Naming pattern. Use Curly brackets `{}` to select placeholders",
			Required:    true,
		},
		"conditions": {
			Type:        schema.TypeList,
			Description: "Defines the conditions when the naming rule should apply",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(RuleConditions).Schema()},
		},
		"use_or_conditions": {
			Type:        schema.TypeBool,
			Description: "If set to `true` the conditions will be connected by logical OR instead of logical AND",
			Optional:    true,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"template":          me.Template,
		"conditions":        me.Conditions,
		"use_or_conditions": me.UseOrConditions,
	})
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"template":          &me.Template,
		"conditions":        &me.Conditions,
		"use_or_conditions": &me.UseOrConditions,
	})
}
