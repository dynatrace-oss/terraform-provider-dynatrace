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
	Enabled            bool            `json:"enabled"` // This setting is enabled (`true`) or disabled (`false`)
	ReplacementPattern *string         `json:"replacementPattern,omitempty"`
	ReplacementType    ReplacementType `json:"replacementType"` // Possible Values: `SHA256`, `STRING`
	RuleName           string          `json:"ruleName"`        // Rule Name
	RuleRegex          *string         `json:"ruleRegex,omitempty"`
	RuleType           RuleType        `json:"ruleType"` // Possible Values: `REGEX`, `VAR_NAME`
	RuleVarName        *string         `json:"ruleVarName,omitempty"`
	InsertAfter        string          `json:"-"`
}

func (me *Settings) Name() string {
	return me.RuleName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"replacement_pattern": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"replacement_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `SHA256`, `STRING`",
			Required:    true,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule Name",
			Required:    true,
		},
		"rule_regex": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"rule_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `REGEX`, `VAR_NAME`",
			Required:    true,
		},
		"rule_var_name": {
			Type:        schema.TypeString,
			Description: "no documentation available",
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
	if (me.ReplacementPattern == nil) && (string(me.ReplacementType) == "STRING") {
		return fmt.Errorf("'replacement_pattern' must be specified if 'replacement_type' is set to '%v'", me.ReplacementType)
	}
	if (me.RuleRegex == nil) && (string(me.RuleType) == "REGEX") {
		return fmt.Errorf("'rule_regex' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.RuleVarName == nil) && (string(me.RuleType) == "VAR_NAME") {
		return fmt.Errorf("'rule_var_name' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
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
