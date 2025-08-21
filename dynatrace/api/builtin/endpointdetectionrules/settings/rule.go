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

package endpointdetectionrules

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rule struct {
	Condition            *string         `json:"condition,omitempty"` // Limits the scope of the endpoint detection rule using [DQL matcher](https://dt-url.net/l603wby) conditions on span and resource attributes.. A rule is applied only if the condition matches, otherwise the ruleset evaluation continues.\n\nIf empty, the condition will always match.
	Description          *string         `json:"description,omitempty"`
	EndpointNameTemplate *string         `json:"endpointNameTemplate,omitempty"` // Specify attribute placeholders in curly braces, e.g. {http.route} or {rpc.method}.. Attribute value placeholders should be specified in curly braces, e.g. {http.route}, {rpc.method}. All attributes used in the placeholder are required for the rule to apply. If any of them is missing, the rule will not be applied and ruleset evaluation continues.\n\nIf the resolved endpoint name on a given span is empty, the request will be ignored.
	IfConditionMatches   ActionToPerform `json:"ifConditionMatches"`             // Possible Values: `DETECT_REQUEST_ON_ENDPOINT`, `SUPPRESS_REQUEST`
	RuleName             string          `json:"ruleName"`                       // Rule name
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "Limits the scope of the endpoint detection rule using [DQL matcher](https://dt-url.net/l603wby) conditions on span and resource attributes.. A rule is applied only if the condition matches, otherwise the ruleset evaluation continues.\n\nIf empty, the condition will always match.",
			Optional:    true, // nullable
		},
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"endpoint_name_template": {
			Type:        schema.TypeString,
			Description: "Specify attribute placeholders in curly braces, e.g. {http.route} or {rpc.method}.. Attribute value placeholders should be specified in curly braces, e.g. {http.route}, {rpc.method}. All attributes used in the placeholder are required for the rule to apply. If any of them is missing, the rule will not be applied and ruleset evaluation continues.\n\nIf the resolved endpoint name on a given span is empty, the request will be ignored.",
			Optional:    true, // precondition
		},
		"if_condition_matches": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DETECT_REQUEST_ON_ENDPOINT`, `SUPPRESS_REQUEST`",
			Required:    true,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":              me.Condition,
		"description":            me.Description,
		"endpoint_name_template": me.EndpointNameTemplate,
		"if_condition_matches":   me.IfConditionMatches,
		"rule_name":              me.RuleName,
	})
}

func (me *Rule) HandlePreconditions() error {
	if (me.EndpointNameTemplate == nil) && (string(me.IfConditionMatches) == "DETECT_REQUEST_ON_ENDPOINT") {
		return fmt.Errorf("'endpoint_name_template' must be specified if 'if_condition_matches' is set to '%v'", me.IfConditionMatches)
	}
	return nil
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":              &me.Condition,
		"description":            &me.Description,
		"endpoint_name_template": &me.EndpointNameTemplate,
		"if_condition_matches":   &me.IfConditionMatches,
		"rule_name":              &me.RuleName,
	})
}
