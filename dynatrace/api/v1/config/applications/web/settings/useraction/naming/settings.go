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

// Settings The settings of user action naming
type Settings struct {
	Placeholders               Placeholders `json:"placeholders"`                     // User action placeholders
	LoadActionNamingRules      Rules        `json:"loadActionNamingRules"`            // User action naming rules for loading actions
	XHRActionNamingRules       Rules        `json:"xhrActionNamingRules"`             // User action naming rules for XHR actions
	CustomActionNamingRules    Rules        `json:"customActionNamingRules"`          // User action naming rules for custom actions
	IgnoreCase                 bool         `json:"ignoreCase"`                       // Case insensitive naming
	UseFirstDetectedLoadAction bool         `json:"useFirstDetectedLoadAction"`       // First load action found under an XHR action should be used when true. Else the deepest one under the xhr action is used
	SplitUserActionsByDomain   bool         `json:"splitUserActionsByDomain"`         // Deactivate this setting if different domains should not result in separate user actions
	QueryParameterCleanups     []string     `json:"queryParameterCleanups,omitempty"` // List of parameters that should be removed from the query before using the query in the user action name
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"placeholders": {
			Type:        schema.TypeList,
			Description: "User action placeholders",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Placeholders).Schema()},
		},
		"load_action_naming_rules": {
			Type:        schema.TypeList,
			Description: "User action naming rules for loading actions",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Rules).Schema()},
		},
		"xhr_action_naming_rules": {
			Type:        schema.TypeList,
			Description: "User action naming rules for XHR actions",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Rules).Schema()},
		},
		"custom_action_naming_rules": {
			Type:        schema.TypeList,
			Description: "User action naming rules for custom actions",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Rules).Schema()},
		},
		"ignore_case": {
			Type:        schema.TypeBool,
			Description: "Case insensitive naming",
			Optional:    true,
		},
		"use_first_detected_load_action": {
			Type:        schema.TypeBool,
			Description: "First load action found under an XHR action should be used when true. Else the deepest one under the xhr action is used",
			Optional:    true,
		},
		"split_user_actions_by_domain": {
			Type:        schema.TypeBool,
			Description: "Deactivate this setting if different domains should not result in separate user actions",
			Optional:    true,
		},
		"query_parameter_cleanups": {
			Type:        schema.TypeSet,
			Description: "User action naming rules for custom actions",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"placeholders":                   me.Placeholders,
		"load_action_naming_rules":       me.LoadActionNamingRules,
		"xhr_action_naming_rules":        me.XHRActionNamingRules,
		"custom_action_naming_rules":     me.CustomActionNamingRules,
		"ignore_case":                    me.IgnoreCase,
		"use_first_detected_load_action": me.UseFirstDetectedLoadAction,
		"split_user_actions_by_domain":   me.SplitUserActionsByDomain,
		"query_parameter_cleanups":       me.QueryParameterCleanups,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"placeholders":                   &me.Placeholders,
		"load_action_naming_rules":       &me.LoadActionNamingRules,
		"xhr_action_naming_rules":        &me.XHRActionNamingRules,
		"custom_action_naming_rules":     &me.CustomActionNamingRules,
		"ignore_case":                    &me.IgnoreCase,
		"use_first_detected_load_action": &me.UseFirstDetectedLoadAction,
		"split_user_actions_by_domain":   &me.SplitUserActionsByDomain,
		"query_parameter_cleanups":       &me.QueryParameterCleanups,
	}); err != nil {
		return err
	}

	if me.Placeholders == nil {
		me.Placeholders = Placeholders{}
	}
	if me.LoadActionNamingRules == nil {
		me.LoadActionNamingRules = Rules{}
	}
	if me.XHRActionNamingRules == nil {
		me.XHRActionNamingRules = Rules{}
	}
	if me.CustomActionNamingRules == nil {
		me.CustomActionNamingRules = Rules{}
	}
	return nil
}
