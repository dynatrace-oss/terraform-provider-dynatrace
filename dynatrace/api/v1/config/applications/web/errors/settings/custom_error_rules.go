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

package errors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomErrorRules []*CustomErrorRule

func (me *CustomErrorRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Description: "Configuration of the custom error in the web application",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomErrorRule).Schema()},
		},
	}
}

func (me CustomErrorRules) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("rule", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *CustomErrorRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

// CustomErrorRule represents configuration of the custom error in the web application
type CustomErrorRule struct {
	KeyPattern     *string                      `json:"keyPattern,omitempty"`   // The key of the error to look for
	KeyMatcher     *CustomErrorRuleKeyMatcher   `json:"keyMatcher,omitempty"`   // The matching operation for the **keyPattern**. Possible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`.
	ValuePattern   *string                      `json:"valuePattern,omitempty"` // The value of the error to look for
	ValueMatcher   *CustomErrorRuleValueMatcher `json:"valueMatcher,omitempty"` // The matching operation for the **valuePattern**. Possible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`.
	Capture        bool                         `json:"capture"`                // Capture (`true`) or ignore (`false`) the error
	ImpactApdex    bool                         `json:"impactApdex"`            // Include (`true`) or exclude (`false`) the error in Apdex calculation
	CustomAlerting bool                         `json:"customAlerting"`         // Include (`true`) or exclude (`false`) the error in Davis AI [problem detection and analysis](https://dt-url.net/a963kd2)
}

func (me *CustomErrorRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key_pattern": {
			Type:        schema.TypeString,
			Description: "The key of the error to look for",
			Optional:    true,
		},
		"key_matcher": {
			Type:        schema.TypeString,
			Description: "The matching operation for the **keyPattern**. Possible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`",
			Optional:    true,
		},
		"value_pattern": {
			Type:        schema.TypeString,
			Description: "The value of the error to look for",
			Optional:    true,
		},
		"value_matcher": {
			Type:        schema.TypeString,
			Description: "The matching operation for the **valuePattern**. Possible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`.",
			Optional:    true,
		},
		"capture": {
			Type:        schema.TypeBool,
			Description: "Capture (`true`) or ignore (`false`) the error",
			Optional:    true,
		},
		"impact_apdex": {
			Type:        schema.TypeBool,
			Description: "Include (`true`) or exclude (`false`) the error in Apdex calculation",
			Optional:    true,
		},
		"custom_alerting": {
			Type:        schema.TypeBool,
			Description: "Include (`true`) or exclude (`false`) the error in Davis AI [problem detection and analysis](https://dt-url.net/a963kd2)",
			Optional:    true,
		},
	}
}

func (me *CustomErrorRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key_pattern":     me.KeyPattern,
		"key_matcher":     me.KeyMatcher,
		"value_pattern":   me.ValuePattern,
		"value_matcher":   me.ValueMatcher,
		"capture":         me.Capture,
		"impact_apdex":    me.ImpactApdex,
		"custom_alerting": me.CustomAlerting,
	})
}

func (me *CustomErrorRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key_pattern":     &me.KeyPattern,
		"key_matcher":     &me.KeyMatcher,
		"value_pattern":   &me.ValuePattern,
		"value_matcher":   &me.ValueMatcher,
		"capture":         &me.Capture,
		"impact_apdex":    &me.ImpactApdex,
		"custom_alerting": &me.CustomAlerting,
	})
}
