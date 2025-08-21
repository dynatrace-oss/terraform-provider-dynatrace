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

package failuredetectionrulesets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomRules []*CustomRule

func (me *CustomRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fail_on_custom_rule": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CustomRule).Schema()},
		},
	}
}

func (me CustomRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("fail_on_custom_rule", me)
}

func (me *CustomRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("fail_on_custom_rule", me)
}

type CustomRule struct {
	DqlCondition string `json:"dqlCondition"` // Custom rule based on span attributes using [DQL matcher](https://dt-url.net/l603wby).
	Enabled      bool   `json:"enabled"`      // This setting is enabled (`true`) or disabled (`false`)
	RuleName     string `json:"ruleName"`     // Rule name
}

func (me *CustomRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dql_condition": {
			Type:        schema.TypeString,
			Description: "Custom rule based on span attributes using [DQL matcher](https://dt-url.net/l603wby).",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
	}
}

func (me *CustomRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dql_condition": me.DqlCondition,
		"enabled":       me.Enabled,
		"rule_name":     me.RuleName,
	})
}

func (me *CustomRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dql_condition": &me.DqlCondition,
		"enabled":       &me.Enabled,
		"rule_name":     &me.RuleName,
	})
}
