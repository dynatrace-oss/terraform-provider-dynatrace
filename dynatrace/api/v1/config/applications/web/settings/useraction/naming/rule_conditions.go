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

type RuleConditions []*RuleCondition

func (me *RuleConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeList,
			Description: "Defines the conditions when the naming rule should apply",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(RuleCondition).Schema()},
		},
	}
}

func (me RuleConditions) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("condition", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *RuleConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("condition", me)
}

// RuleCondition The settings of conditions for user action naming
type RuleCondition struct {
	Operand1 string   `json:"operand1"`           // Must be a defined placeholder wrapped in curly braces
	Operand2 *string  `json:"operand2,omitempty"` // Must be null if operator is `IS_EMPTY`, a regex if operator is `MATCHES_REGULAR_ERPRESSION`. In all other cases the value can be a freetext or a placeholder wrapped in curly braces
	Operator Operator `json:"operator"`           // The operator of the condition. Possible values are `CONTAINS`, `ENDS_WITH`, `EQUALS`, `IS_EMPTY`, `IS_NOT_EMPTY`, `MATCHES_REGULAR_EXPRESSION`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUALS`, `NOT_MATCHES_REGULAR_EXPRESSION`, `NOT_STARTS_WITH` and `STARTS_WITH`.
}

func (me *RuleCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"operand1": {
			Type:        schema.TypeString,
			Description: "Must be a defined placeholder wrapped in curly braces",
			Required:    true,
		},
		"operand2": {
			Type:        schema.TypeString,
			Description: "Must be null if operator is `IS_EMPTY`, a regex if operator is `MATCHES_REGULAR_ERPRESSION`. In all other cases the value can be a freetext or a placeholder wrapped in curly braces",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "The operator of the condition. Possible values are `CONTAINS`, `ENDS_WITH`, `EQUALS`, `IS_EMPTY`, `IS_NOT_EMPTY`, `MATCHES_REGULAR_EXPRESSION`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUALS`, `NOT_MATCHES_REGULAR_EXPRESSION`, `NOT_STARTS_WITH` and `STARTS_WITH`.",
			Required:    true,
		},
	}
}

func (me *RuleCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"operand1": me.Operand1,
		"operand2": me.Operand2,
		"operator": me.Operator,
	})
}

func (me *RuleCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"operand1": &me.Operand1,
		"operand2": &me.Operand2,
		"operator": &me.Operator,
	})
}
