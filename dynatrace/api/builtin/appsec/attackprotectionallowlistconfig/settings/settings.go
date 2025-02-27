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

package attackprotectionallowlistconfig

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AttackHandling              *AttackHandling             `json:"attackHandling"`                        // Step 1: Define attack control for chosen criteria
	Enabled                     bool                        `json:"enabled"`                               // This setting is enabled (`true`) or disabled (`false`)
	Metadata                    *Metadata                   `json:"metadata"`                              // Step 4: Leave comment (optional)
	ResourceAttributeConditions ResourceAttributeConditions `json:"resourceAttributeConditions,omitempty"` // When you add multiple conditions, the rule applies if all conditions apply.\n\nIf you want the rule to apply only to a subset of your environment, provide the resource attributes that should be used to identify that part of the environment.
	RuleName                    *string                     `json:"ruleName,omitempty"`                    // Rule name
	Rules                       AgentSideCriterias          `json:"rules"`                                 // Provide conditions that must be met by the detection finding you want to allowlist.
	InsertAfter                 string                      `json:"-"`
}

func (me *Settings) Name() string {
	return uuid.NewString()
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attack_handling": {
			Type:        schema.TypeList,
			Description: "Step 1: Define attack control for chosen criteria",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(AttackHandling).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"criteria": {
			Type:        schema.TypeList,
			Description: "Step 1: Define criteria. Please specify at least one of source IP or attack pattern.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Criteria).Schema()},
			MinItems:    1,
			MaxItems:    1,
			Deprecated:  "The `criteria` attribute has been deprecated, please use the `rules` and `resource_attribute_conditions` attributes instead.",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Step 4: Leave comment (optional)",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Metadata).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"resource_attribute_conditions": {
			Type:        schema.TypeList,
			Description: "When you add multiple conditions, the rule applies if all conditions apply.\n\nIf you want the rule to apply only to a subset of your environment, provide the resource attributes that should be used to identify that part of the environment.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(ResourceAttributeConditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Optional:    true, // nullable
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "Provide conditions that must be met by the detection finding you want to allowlist.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(AgentSideCriterias).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attack_handling":               me.AttackHandling,
		"enabled":                       me.Enabled,
		"metadata":                      me.Metadata,
		"resource_attribute_conditions": me.ResourceAttributeConditions,
		"rule_name":                     me.RuleName,
		"rules":                         me.Rules,
		"insert_after":                  me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attack_handling":               &me.AttackHandling,
		"enabled":                       &me.Enabled,
		"metadata":                      &me.Metadata,
		"resource_attribute_conditions": &me.ResourceAttributeConditions,
		"rule_name":                     &me.RuleName,
		"rules":                         &me.Rules,
		"insert_after":                  &me.InsertAfter,
	})
}
