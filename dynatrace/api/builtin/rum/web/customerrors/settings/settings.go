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

package customerrors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ErrorRules                           CustomErrorRules `json:"errorRules"`                           // (Field has overlap with `dynatrace_application_error_rules`)
	IgnoreCustomErrorsInApdexCalculation bool             `json:"ignoreCustomErrorsInApdexCalculation"` // (Field has overlap with `dynatrace_application_error_rules`) This setting overrides Apdex settings for individual rules listed below
	Scope                                string           `json:"-" scope:"scope"`                      // The scope of this setting (APPLICATION)
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_rules": {
			Type:        schema.TypeList,
			Description: "(Field has overlap with `dynatrace_application_error_rules`)",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(CustomErrorRules).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"ignore_custom_errors_in_apdex_calculation": {
			Type:        schema.TypeBool,
			Description: "(Field has overlap with `dynatrace_application_error_rules`) This setting overrides Apdex settings for individual rules listed below",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (APPLICATION)",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"error_rules": me.ErrorRules,
		"ignore_custom_errors_in_apdex_calculation": me.IgnoreCustomErrorsInApdexCalculation,
		"scope": me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"error_rules": &me.ErrorRules,
		"ignore_custom_errors_in_apdex_calculation": &me.IgnoreCustomErrorsInApdexCalculation,
		"scope": &me.Scope,
	})
}
