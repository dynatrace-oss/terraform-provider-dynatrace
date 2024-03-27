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

package database

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	BooleanValue *bool        `json:"booleanValue,omitempty"` // State of boolean feature flag
	Name         string       `json:"name"`                   // Name of the feature
	NumberValue  *int         `json:"numberValue,omitempty"`  // State of numeric feature flag
	StringValue  *string      `json:"stringValue,omitempty"`  // State of textual feature flag
	Type         FeatureTypes `json:"type"`                   // Possible Values: `Boolean`, `Number`, `String`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"boolean_value": {
			Type:        schema.TypeBool,
			Description: "State of boolean feature flag",
			Optional:    true, // precondition
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the feature",
			Required:    true,
		},
		"number_value": {
			Type:        schema.TypeInt,
			Description: "State of numeric feature flag",
			Optional:    true, // precondition
		},
		"string_value": {
			Type:        schema.TypeString,
			Description: "State of textual feature flag",
			Optional:    true, // precondition
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Boolean`, `Number`, `String`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"boolean_value": me.BooleanValue,
		"name":          me.Name,
		"number_value":  me.NumberValue,
		"string_value":  me.StringValue,
		"type":          me.Type,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.BooleanValue == nil) && (string(me.Type) == "boolean") {
		me.BooleanValue = opt.NewBool(false)
	}
	if (me.NumberValue == nil) && (string(me.Type) == "number") {
		me.NumberValue = opt.NewInt(0)
	}
	if (me.StringValue == nil) && (string(me.Type) == "string") {
		return fmt.Errorf("'string_value' must be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"boolean_value": &me.BooleanValue,
		"name":          &me.Name,
		"number_value":  &me.NumberValue,
		"string_value":  &me.StringValue,
		"type":          &me.Type,
	})
}
