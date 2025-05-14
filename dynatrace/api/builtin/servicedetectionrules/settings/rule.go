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

package servicedetectionrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rule struct {
	AdditionalRequiredAttributes []string `json:"additionalRequiredAttributes,omitempty"` // Define resource attributes that should not be part of the name but are required to detect the service, e.g. service.namespace or k8s.workload.kind.. Attributes specified here are required to apply the rule. If any of them is missing, the rule will not be applied and ruleset evaluation continues.\n\nAll attribute values contribute to the final service ID.
	Condition                    *string  `json:"condition,omitempty"`                    // Limits the scope of the service detection rule using [DQL matcher](https://dt-url.net/l603wby) conditions on resource attributes.. A rule is applied only if the condition matches, otherwise the ruleset evaluation continues.\n\nIf empty, the condition will always match.
	Description                  *string  `json:"description,omitempty"`
	RuleName                     string   `json:"ruleName"`            // Rule name
	ServiceNameTemplate          string   `json:"serviceNameTemplate"` // Specify resource attribute placeholders in curly braces, e.g. {service.name} or {k8s.workload.name}.. All attributes used in the placeholder are required for the rule to apply. If any of them is missing, the rule will not be applied and ruleset evaluation continues.\n\nAll resolved attribute values contribute to the final service ID.
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"additional_required_attributes": {
			Type:        schema.TypeSet,
			Description: "Define resource attributes that should not be part of the name but are required to detect the service, e.g. service.namespace or k8s.workload.kind.. Attributes specified here are required to apply the rule. If any of them is missing, the rule will not be applied and ruleset evaluation continues.\n\nAll attribute values contribute to the final service ID.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"condition": {
			Type:        schema.TypeString,
			Description: "Limits the scope of the service detection rule using [DQL matcher](https://dt-url.net/l603wby) conditions on resource attributes.. A rule is applied only if the condition matches, otherwise the ruleset evaluation continues.\n\nIf empty, the condition will always match.",
			Optional:    true, // nullable
		},
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
		"service_name_template": {
			Type:        schema.TypeString,
			Description: "Specify resource attribute placeholders in curly braces, e.g. {service.name} or {k8s.workload.name}.. All attributes used in the placeholder are required for the rule to apply. If any of them is missing, the rule will not be applied and ruleset evaluation continues.\n\nAll resolved attribute values contribute to the final service ID.",
			Required:    true,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"additional_required_attributes": me.AdditionalRequiredAttributes,
		"condition":                      me.Condition,
		"description":                    me.Description,
		"rule_name":                      me.RuleName,
		"service_name_template":          me.ServiceNameTemplate,
	})
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"additional_required_attributes": &me.AdditionalRequiredAttributes,
		"condition":                      &me.Condition,
		"description":                    &me.Description,
		"rule_name":                      &me.RuleName,
		"service_name_template":          &me.ServiceNameTemplate,
	})
}
