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

package generalparameters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CompareOperation struct {
	CaseSensitive        *bool    `json:"caseSensitive,omitempty"` // Case sensitive
	CompareOperationType string   `json:"compareOperationType"`    // Apply this comparison
	DoubleValue          *float64 `json:"doubleValue,omitempty"`   // Value
	IntValue             *int     `json:"intValue,omitempty"`      // Value
	TextValue            *string  `json:"textValue,omitempty"`     // Value
}

func (me *CompareOperation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "Case sensitive",
			Optional:    true,
		},
		"compare_operation_type": {
			Type:        schema.TypeString,
			Description: "Apply this comparison",
			Required:    true,
		},
		"double_value": {
			Type:        schema.TypeFloat,
			Description: "Value",
			Optional:    true,
		},
		"int_value": {
			Type:        schema.TypeInt,
			Description: "Value",
			Optional:    true,
		},
		"text_value": {
			Type:        schema.TypeString,
			Description: "Value",
			Optional:    true,
		},
	}
}

func (me *CompareOperation) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"case_sensitive":         me.CaseSensitive,
		"compare_operation_type": me.CompareOperationType,
		"double_value":           me.DoubleValue,
		"int_value":              me.IntValue,
		"text_value":             me.TextValue,
	})
}

func (me *CompareOperation) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"case_sensitive":         &me.CaseSensitive,
		"compare_operation_type": &me.CompareOperationType,
		"double_value":           &me.DoubleValue,
		"int_value":              &me.IntValue,
		"text_value":             &me.TextValue,
	})
	expectedValues := []string{"STRING_EQUALS", "NOT_STRING_EQUALS", "STARTS_WITH", "NOT_STARTS_WITH", "CONTAINS", "NOT_CONTAINS", "ENDS_WITH", "NOT_ENDS_WITH"}
	if me.CaseSensitive == nil && stringInSlice(me.CompareOperationType, expectedValues) {
		me.CaseSensitive = opt.NewBool(false)
	}

	intExpectedValues := []string{"INTEGER_EQUALS", "NOT_INTEGER_EQUALS", "INTEGER_GREATER_THAN", "INTEGER_GREATER_THAN_OR_EQUALS", "INTEGER_LESS_THAN", "INTEGER_LESS_THAN_OR_EQUALS"}
	if me.IntValue == nil && stringInSlice(me.CompareOperationType, intExpectedValues) {
		me.IntValue = opt.NewInt(0)
	}

	doubleExpectedValues := []string{"DOUBLE_EQUALS", "NOT_DOUBLE_EQUALS", "DOUBLE_GREATER_THAN", "DOUBLE_GREATER_THAN_OR_EQUALS", "DOUBLE_LESS_THAN", "DOUBLE_LESS_THAN_OR_EQUALS"}
	if me.DoubleValue == nil && stringInSlice(me.CompareOperationType, doubleExpectedValues) {
		me.DoubleValue = opt.NewFloat64(0)
	}

	return err
}
