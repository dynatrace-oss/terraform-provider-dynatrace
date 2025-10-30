/**
* @license
* Copyright 2025 Dynatrace LLC
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

package ingestsources

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GenericValueAssignment struct {
	Constant           *string                        `json:"constant,omitempty"`           // Constant value
	Field              *ValueAssignmentFromFieldEntry `json:"field,omitempty"`              // Value from field
	MultiValueConstant []string                       `json:"multiValueConstant,omitempty"` // Constant multi value
	Type               AssignmentType                 `json:"type"`                         // Type of value assignment. Possible Values: `constant`, `field`, `multiValueConstant`.
}

func (me *GenericValueAssignment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"constant": {
			Type:        schema.TypeString,
			Description: "Constant value",
			Optional:    true, // precondition
		},
		"field": {
			Type:        schema.TypeList,
			Description: "Value from field",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ValueAssignmentFromFieldEntry).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"multi_value_constant": {
			Type:        schema.TypeList,
			Description: "Constant multi value",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Type of value assignment. Possible Values: `constant`, `field`, `multiValueConstant`.",
			Required:    true,
		},
	}
}

func (me *GenericValueAssignment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"constant":             me.Constant,
		"field":                me.Field,
		"multi_value_constant": me.MultiValueConstant,
		"type":                 me.Type,
	})
}

func (me *GenericValueAssignment) HandlePreconditions() error {
	if (me.Constant == nil) && (string(me.Type) == "constant") {
		me.Constant = opt.NewString("")
	}
	if (me.Field == nil) && (string(me.Type) == "field") {
		return fmt.Errorf("'field' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Field != nil) && (string(me.Type) != "field") {
		return fmt.Errorf("'field' must not be specified if 'type' is set to '%v'", me.Type)
	}
	// ---- MultiValueConstant []string -> {"expectedValue":"multiValueConstant","property":"type","type":"EQUALS"}
	return nil
}

func (me *GenericValueAssignment) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"constant":             &me.Constant,
		"field":                &me.Field,
		"multi_value_constant": &me.MultiValueConstant,
		"type":                 &me.Type,
	})
}
