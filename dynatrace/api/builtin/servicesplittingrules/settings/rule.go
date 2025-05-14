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

package servicesplittingrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rule struct {
	Condition                  *string   `json:"condition,omitempty"` // Limits the scope of the service splitting rule using [DQL matcher](https://dt-url.net/l603wby) conditions on resource attributes..  A rule is applied only if the condition matches, otherwise the ruleset evaluation continues.\n\nIf empty, the condition will always match.
	Description                *string   `json:"description,omitempty"`
	RuleName                   string    `json:"ruleName"`                             // Rule name
	ServiceSplittingAttributes SplitBies `json:"serviceSplittingAttributes,omitempty"` // Define the entire set of resource attributes that should split your services in the matching scope.. Each attribute that exists will contribute to the final service ID.
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "Limits the scope of the service splitting rule using [DQL matcher](https://dt-url.net/l603wby) conditions on resource attributes..  A rule is applied only if the condition matches, otherwise the ruleset evaluation continues.\n\nIf empty, the condition will always match.",
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
		"service_splitting_attributes": {
			Type:        schema.TypeList,
			Description: "Define the entire set of resource attributes that should split your services in the matching scope.. Each attribute that exists will contribute to the final service ID.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(SplitBies).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":                    me.Condition,
		"description":                  me.Description,
		"rule_name":                    me.RuleName,
		"service_splitting_attributes": me.ServiceSplittingAttributes,
	})
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":                    &me.Condition,
		"description":                  &me.Description,
		"rule_name":                    &me.RuleName,
		"service_splitting_attributes": &me.ServiceSplittingAttributes,
	})
}
