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

package autotagging

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Description               *string `json:"description,omitempty"` // Description
	Name                      string  `json:"name"`                  // Tag name
	Rules                     Rules   `json:"rules,omitempty"`       // Rules
	RulesMaintainedExternally bool    `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Tag name",
			Required:    true,
		},
		"rules_maintained_externally": {
			Type:        schema.TypeBool,
			Description: "If `true` the specified rules are ignored with the assumption that they're maintained externally or via `dynatrace_autotag_rules`",
			Optional:    true,
			Default:     false,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "Rules",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Rules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":                 me.Description,
		"name":                        me.Name,
		"rules":                       me.Rules,
		"rules_maintained_externally": me.RulesMaintainedExternally,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"description":                 &me.Description,
		"name":                        &me.Name,
		"rules":                       &me.Rules,
		"rules_maintained_externally": &me.RulesMaintainedExternally,
	}); err != nil {
		return err
	}
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/895
	// Only known workaround is to ignore these blocks
	// Side node: This issue doesn't emerge upon first `terraform apply`
	//            It occurs when changes are happening (UPDATE)
	if len(me.Rules) > 0 {
		rules := Rules{}
		for _, rule := range me.Rules {
			// value_format AND rule type normally are required to contain SOMETHING
			// if they're nil/empty, we know it's because of the bug in the SDK
			if rule.ValueFormat == nil && len(rule.Type) == 0 {
				continue
			}
			rules = append(rules, rule)
		}
		me.Rules = rules
	}

	return nil
}
