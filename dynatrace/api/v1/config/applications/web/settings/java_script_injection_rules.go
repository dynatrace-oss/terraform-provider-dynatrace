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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JavaScriptInjectionRules []*JavaScriptInjectionRule

func (me *JavaScriptInjectionRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Description: "Java script injection rule",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(JavaScriptInjectionRule).Schema()},
		},
	}
}

func (me JavaScriptInjectionRules) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("rule", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *JavaScriptInjectionRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

type JavaScriptInjectionRule struct {
	Enabled     bool             `json:"enabled"`               // The enable or disable rule of the java script injection
	URLOperator URLOperator      `json:"urlOperator"`           // The url operator of the java script injection. Possible values are `ALL_PAGES`, `CONTAINS`, `ENDS_WITH`, `EQUALS` and `STARTS_WITH`.
	URLPattern  *string          `json:"urlPattern,omitempty"`  // The url pattern of the java script injection
	Rule        JSInjectionRule  `json:"rule"`                  // The url rule of the java script injection. Possible values are `AFTER_SPECIFIC_HTML`, `AUTOMATIC_INJECTION`, `BEFORE_SPECIFIC_HTML` and `DO_NOT_INJECT`.
	HTMLPattern *string          `json:"htmlPattern,omitempty"` // The HTML pattern of the java script injection
	Target      *InjectionTarget `json:"target,omitempty"`      // The target against which the rule of the java script injection should be matched. Possible values are `PAGE_QUERY` and `URL`.
}

func (me *JavaScriptInjectionRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "`fetch()` request capture enabled/disabled",
			Optional:    true,
		},
		"url_operator": {
			Type:        schema.TypeString,
			Description: "The url operator of the java script injection. Possible values are `ALL_PAGES`, `CONTAINS`, `ENDS_WITH`, `EQUALS` and `STARTS_WITH`.",
			Required:    true,
		},
		"url_pattern": {
			Type:        schema.TypeString,
			Description: "The url pattern of the java script injection",
			Optional:    true,
		},
		"rule": {
			Type:        schema.TypeString,
			Description: "The url rule of the java script injection. Possible values are `AFTER_SPECIFIC_HTML`, `AUTOMATIC_INJECTION`, `BEFORE_SPECIFIC_HTML` and `DO_NOT_INJECT`.",
			Required:    true,
		},
		"html_pattern": {
			Type:        schema.TypeString,
			Description: "The HTML pattern of the java script injection",
			Optional:    true,
		},
		"target": {
			Type:        schema.TypeString,
			Description: "The target against which the rule of the java script injection should be matched. Possible values are `PAGE_QUERY` and `URL`.",
			Optional:    true,
		},
	}
}

func (me *JavaScriptInjectionRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":      me.Enabled,
		"url_operator": me.URLOperator,
		"url_pattern":  me.URLPattern,
		"rule":         me.Rule,
		"html_pattern": me.HTMLPattern,
		"target":       me.Target,
	})
}

func (me *JavaScriptInjectionRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":      &me.Enabled,
		"url_operator": &me.URLOperator,
		"url_pattern":  &me.URLPattern,
		"rule":         &me.Rule,
		"html_pattern": &me.HTMLPattern,
		"target":       &me.Target,
	})
}
