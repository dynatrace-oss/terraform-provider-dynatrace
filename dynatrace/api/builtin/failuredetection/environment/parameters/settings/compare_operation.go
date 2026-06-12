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

package parameters

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// compareOperation. A comparison operation that tests the value of a request attribute. The type of comparison and the value to compare against depend on the data type of the request attribute.
type CompareOperation struct {
	CaseSensitive        *bool    `json:"caseSensitive,omitempty"` // If `true`, the comparison is case-sensitive. Only applicable for string comparison types. Default: `false`.
	CompareOperationType string   `json:"compareOperationType"`    // The type of comparison to apply. Available types depend on the data type of the request attribute:\n * String types support `STRING_EXISTS`, `STRING_EQUALS`, `NOT_STRING_EQUALS`, `STARTS_WITH`, `NOT_STARTS_WITH`, `CONTAINS`, `NOT_CONTAINS`, `ENDS_WITH`, `NOT_ENDS_WITH`.\n * Integer types support `INTEGER_EQUALS` and related comparisons;\n * Double types support `DOUBLE_EQUALS` and related comparisons.
	DoubleValue          *float64 `json:"doubleValue,omitempty"`   // The floating-point value to compare the request attribute against. Only applicable for double comparison types.
	IntValue             *int     `json:"intValue,omitempty"`      // The integer value to compare the request attribute against. Only applicable for integer comparison types.
	TextValue            *string  `json:"textValue,omitempty"`     // The text value to compare the request attribute against. Only applicable for string comparison types.
}

func (me *CompareOperation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "If `true`, the comparison is case-sensitive. Only applicable for string comparison types. Default: `false`.",
			Optional:    true, // precondition
		},
		"compare_operation_type": {
			Type:        schema.TypeString,
			Description: "The type of comparison to apply. Available types depend on the data type of the request attribute:\n * String types support `STRING_EXISTS`, `STRING_EQUALS`, `NOT_STRING_EQUALS`, `STARTS_WITH`, `NOT_STARTS_WITH`, `CONTAINS`, `NOT_CONTAINS`, `ENDS_WITH`, `NOT_ENDS_WITH`.\n * Integer types support `INTEGER_EQUALS` and related comparisons;\n * Double types support `DOUBLE_EQUALS` and related comparisons.",
			Required:    true,
		},
		"double_value": {
			Type:        schema.TypeFloat,
			Description: "The floating-point value to compare the request attribute against. Only applicable for double comparison types.",
			Optional:    true, // precondition
		},
		"int_value": {
			Type:        schema.TypeInt,
			Description: "The integer value to compare the request attribute against. Only applicable for integer comparison types.",
			Optional:    true, // precondition
		},
		"text_value": {
			Type:        schema.TypeString,
			Description: "The text value to compare the request attribute against. Only applicable for string comparison types.",
			Optional:    true, // precondition
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

func (me *CompareOperation) HandlePreconditions() error {
	if (me.CaseSensitive == nil) && (slices.Contains([]string{"STRING_EQUALS", "NOT_STRING_EQUALS", "STARTS_WITH", "NOT_STARTS_WITH", "CONTAINS", "NOT_CONTAINS", "ENDS_WITH", "NOT_ENDS_WITH"}, string(me.CompareOperationType))) {
		me.CaseSensitive = new(false)
	}
	if (me.DoubleValue == nil) && (slices.Contains([]string{"DOUBLE_EQUALS", "NOT_DOUBLE_EQUALS", "DOUBLE_GREATER_THAN", "DOUBLE_GREATER_THAN_OR_EQUALS", "DOUBLE_LESS_THAN", "DOUBLE_LESS_THAN_OR_EQUALS"}, string(me.CompareOperationType))) {
		me.DoubleValue = new(0.0)
	}
	if (me.IntValue == nil) && (slices.Contains([]string{"INTEGER_EQUALS", "NOT_INTEGER_EQUALS", "INTEGER_GREATER_THAN", "INTEGER_GREATER_THAN_OR_EQUALS", "INTEGER_LESS_THAN", "INTEGER_LESS_THAN_OR_EQUALS"}, string(me.CompareOperationType))) {
		me.IntValue = new(0)
	}
	if (me.CaseSensitive != nil) && (!slices.Contains([]string{"STRING_EQUALS", "NOT_STRING_EQUALS", "STARTS_WITH", "NOT_STARTS_WITH", "CONTAINS", "NOT_CONTAINS", "ENDS_WITH", "NOT_ENDS_WITH"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'case_sensitive' must not be specified unless 'compare_operation_type' is one of ['STRING_EQUALS', 'NOT_STRING_EQUALS', 'STARTS_WITH', 'NOT_STARTS_WITH', 'CONTAINS', 'NOT_CONTAINS', 'ENDS_WITH', 'NOT_ENDS_WITH']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.DoubleValue != nil) && (!slices.Contains([]string{"DOUBLE_EQUALS", "NOT_DOUBLE_EQUALS", "DOUBLE_GREATER_THAN", "DOUBLE_GREATER_THAN_OR_EQUALS", "DOUBLE_LESS_THAN", "DOUBLE_LESS_THAN_OR_EQUALS"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'double_value' must not be specified unless 'compare_operation_type' is one of ['DOUBLE_EQUALS', 'NOT_DOUBLE_EQUALS', 'DOUBLE_GREATER_THAN', 'DOUBLE_GREATER_THAN_OR_EQUALS', 'DOUBLE_LESS_THAN', 'DOUBLE_LESS_THAN_OR_EQUALS']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.IntValue != nil) && (!slices.Contains([]string{"INTEGER_EQUALS", "NOT_INTEGER_EQUALS", "INTEGER_GREATER_THAN", "INTEGER_GREATER_THAN_OR_EQUALS", "INTEGER_LESS_THAN", "INTEGER_LESS_THAN_OR_EQUALS"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'int_value' must not be specified unless 'compare_operation_type' is one of ['INTEGER_EQUALS', 'NOT_INTEGER_EQUALS', 'INTEGER_GREATER_THAN', 'INTEGER_GREATER_THAN_OR_EQUALS', 'INTEGER_LESS_THAN', 'INTEGER_LESS_THAN_OR_EQUALS']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.TextValue != nil) && (!slices.Contains([]string{"STRING_EQUALS", "NOT_STRING_EQUALS", "STARTS_WITH", "NOT_STARTS_WITH", "CONTAINS", "NOT_CONTAINS", "ENDS_WITH", "NOT_ENDS_WITH"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'text_value' must not be specified unless 'compare_operation_type' is one of ['STRING_EQUALS', 'NOT_STRING_EQUALS', 'STARTS_WITH', 'NOT_STARTS_WITH', 'CONTAINS', 'NOT_CONTAINS', 'ENDS_WITH', 'NOT_ENDS_WITH']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.TextValue == nil) && (slices.Contains([]string{"STRING_EQUALS", "NOT_STRING_EQUALS", "STARTS_WITH", "NOT_STARTS_WITH", "CONTAINS", "NOT_CONTAINS", "ENDS_WITH", "NOT_ENDS_WITH"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'text_value' must be specified when 'compare_operation_type' is one of ['STRING_EQUALS', 'NOT_STRING_EQUALS', 'STARTS_WITH', 'NOT_STARTS_WITH', 'CONTAINS', 'NOT_CONTAINS', 'ENDS_WITH', 'NOT_ENDS_WITH']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	return nil
}

func (me *CompareOperation) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"case_sensitive":         &me.CaseSensitive,
		"compare_operation_type": &me.CompareOperationType,
		"double_value":           &me.DoubleValue,
		"int_value":              &me.IntValue,
		"text_value":             &me.TextValue,
	})
}
