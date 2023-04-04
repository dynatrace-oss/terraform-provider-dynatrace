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

package useractioncustommetrics

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type Filters []*Filter

func (me *Filters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Filter).Schema()},
		},
	}
}

func (me Filters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *Filters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter", me)
}

type Filter struct {
	FieldName string   `json:"fieldName"` // Field name
	Operator  Operator `json:"operator"`  // Possible Values: `EQUALS`, `GREATER_THAN`, `GREATER_THAN_OR_EQUAL_TO`, `IN`, `IS_NOT_NULL`, `IS_NULL`, `LESS_THAN`, `LESS_THAN_OR_EQUAL_TO`, `LIKE`, `NOT_EQUAL`, `NOT_LIKE`, `STARTS_WITH`
	Value     *string  `json:"value,omitempty"`
	ValueIn   []string `json:"valueIn,omitempty"` // Values
}

func (me *Filter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_name": {
			Type:        schema.TypeString,
			Description: "Field name",
			Required:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Possible Values: `EQUALS`, `GREATER_THAN`, `GREATER_THAN_OR_EQUAL_TO`, `IN`, `IS_NOT_NULL`, `IS_NULL`, `LESS_THAN`, `LESS_THAN_OR_EQUAL_TO`, `LIKE`, `NOT_EQUAL`, `NOT_LIKE`, `STARTS_WITH`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"value_in": {
			Type:        schema.TypeList,
			Description: "Values",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Filter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"field_name": me.FieldName,
		"operator":   me.Operator,
		"value":      me.Value,
		"value_in":   me.ValueIn,
	})
}

func (me *Filter) HandlePreconditions() error {
	if me.Value == nil && slices.Contains([]string{"EQUALS", "NOT_EQUAL", "LIKE", "LESS_THAN", "LESS_THAN_OR_EQUAL_TO", "GREATER_THAN", "GREATER_THAN_OR_EQUAL_TO", "NOT_LIKE", "STARTS_WITH"}, string(me.Operator)) {
		return fmt.Errorf("'value' must be specified if 'operator' is set to '%v'", me.Operator)
	}
	// ---- ValueIn []string -> {"expectedValue":"IN","property":"operator","type":"EQUALS"}
	return nil
}

func (me *Filter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"field_name": &me.FieldName,
		"operator":   &me.Operator,
		"value":      &me.Value,
		"value_in":   &me.ValueIn,
	})
}
