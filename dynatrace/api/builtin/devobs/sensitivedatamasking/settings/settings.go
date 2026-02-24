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

package sensitivedatamasking

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ComparisonType     *ComparisonType `json:"comparisonType,omitempty"`     // Select how the variable name should be matched. Possible values: `CONTAINS`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`
	Enabled            bool            `json:"enabled"`                      // This setting is enabled (`true`) or disabled (`false`)
	ReplacementPattern *string         `json:"replacementPattern,omitempty"` // Replacement Pattern
	ReplacementType    ReplacementType `json:"replacementType"`              // Choose how the sensitive data should be replaced. Possible values: `SHA256`, `STRING`
	RuleName           string          `json:"ruleName"`                     // Rule Name
	RuleRegex          *string         `json:"ruleRegex,omitempty"`          // Regex Pattern
	RuleType           RuleType        `json:"ruleType"`                     // Choose whether to redact by variable name or regex. Possible values: `REGEX`, `VAR_NAME`
	RuleVarName        *string         `json:"ruleVarName,omitempty"`        // Variable Name
	InsertAfter        string          `json:"-"`
}

func (me *Settings) Name() string {
	return me.RuleName
}

const defaultComparisonType = ComparisonType("EQUALS")

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"comparison_type": {
			Type:        schema.TypeString,
			Description: "Select how the variable name should be matched. Possible values: `CONTAINS`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`",
			Optional:    true, // precondition
			// new conditionally required field with default value, so we set the default value in HandlePreconditions to avoid breaking existing configurations
			// only set a default value if rule_type == "VAR_NAME". The field must be non-empty for others, therefore, "Default" can't be used.
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if d.Get("rule_type") == "VAR_NAME" {
					if old == "" {
						old = string(defaultComparisonType)
					}
					if new == "" {
						new = string(defaultComparisonType)
					}
				}
				return old == new
			},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"replacement_pattern": {
			Type:        schema.TypeString,
			Description: "Replacement Pattern",
			Optional:    true, // precondition
		},
		"replacement_type": {
			Type:        schema.TypeString,
			Description: "Choose how the sensitive data should be replaced. Possible values: `SHA256`, `STRING`",
			Required:    true,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule Name",
			Required:    true,
		},
		"rule_regex": {
			Type:        schema.TypeString,
			Description: "Regex Pattern",
			Optional:    true, // precondition
		},
		"rule_type": {
			Type:        schema.TypeString,
			Description: "Choose whether to redact by variable name or regex. Possible values: `REGEX`, `VAR_NAME`",
			Required:    true,
		},
		"rule_var_name": {
			Type:        schema.TypeString,
			Description: "Variable Name",
			Optional:    true, // precondition
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
		"comparison_type":     me.ComparisonType,
		"enabled":             me.Enabled,
		"replacement_pattern": me.ReplacementPattern,
		"replacement_type":    me.ReplacementType,
		"rule_name":           me.RuleName,
		"rule_regex":          me.RuleRegex,
		"rule_type":           me.RuleType,
		"rule_var_name":       me.RuleVarName,
		"insert_after":        me.InsertAfter,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.ComparisonType == nil) && (string(me.RuleType) == "VAR_NAME") {
		// new conditionally required field with default value, so we set the default value here to avoid breaking existing configurations
		defaultType := defaultComparisonType
		me.ComparisonType = &defaultType
	}
	if (me.ComparisonType != nil) && (string(me.RuleType) != "VAR_NAME") {
		return fmt.Errorf("'comparison_type' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.ReplacementPattern == nil) && (string(me.ReplacementType) == "STRING") {
		return fmt.Errorf("'replacement_pattern' must be specified if 'replacement_type' is set to '%v'", me.ReplacementType)
	}
	if (me.ReplacementPattern != nil) && (string(me.ReplacementType) != "STRING") {
		return fmt.Errorf("'replacement_pattern' must not be specified if 'replacement_type' is set to '%v'", me.ReplacementType)
	}
	if (me.RuleRegex == nil) && (string(me.RuleType) == "REGEX") {
		return fmt.Errorf("'rule_regex' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.RuleRegex != nil) && (string(me.RuleType) != "REGEX") {
		return fmt.Errorf("'rule_regex' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.RuleVarName == nil) && (string(me.RuleType) == "VAR_NAME") {
		return fmt.Errorf("'rule_var_name' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.RuleVarName != nil) && (string(me.RuleType) != "VAR_NAME") {
		return fmt.Errorf("'rule_var_name' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"comparison_type":     &me.ComparisonType,
		"enabled":             &me.Enabled,
		"replacement_pattern": &me.ReplacementPattern,
		"replacement_type":    &me.ReplacementType,
		"rule_name":           &me.RuleName,
		"rule_regex":          &me.RuleRegex,
		"rule_type":           &me.RuleType,
		"rule_var_name":       &me.RuleVarName,
		"insert_after":        &me.InsertAfter,
	})
}
