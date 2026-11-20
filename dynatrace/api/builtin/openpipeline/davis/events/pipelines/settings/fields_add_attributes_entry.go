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

package pipelines

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FieldsAddAttributesEntries []*FieldsAddAttributesEntry

func (me *FieldsAddAttributesEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(FieldsAddAttributesEntry).Schema()},
		},
	}
}

func (me FieldsAddAttributesEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("field", me)
}

func (me *FieldsAddAttributesEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("field", me)
}

type FieldsAddAttributesEntry struct {
	MultiValue []string            `json:"multiValue,omitempty"` // Multi value
	Name       string              `json:"name"`                 // Fields's name
	Value      *string             `json:"value,omitempty"`      // Field's value
	ValueType  *FieldsAddValueType `json:"valueType,omitempty"`  // Field value type. Possible values: `arrayString`, `boolean`, `double`, `long`, `string`
}

func (me *FieldsAddAttributesEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"multi_value": {
			Type:        schema.TypeList,
			Description: "Multi value",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Fields's name",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Field's value",
			Optional:    true, // precondition
		},
		"value_type": {
			Type:        schema.TypeString,
			Description: "Field value type. Possible values: `arrayString`, `boolean`, `double`, `long`, `string`",
			Optional:    true, // nullable
		},
	}
}

func (me *FieldsAddAttributesEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"multi_value": me.MultiValue,
		"name":        me.Name,
		"value":       me.Value,
		"value_type":  me.ValueType,
	})
}

func (me *FieldsAddAttributesEntry) HandlePreconditions() error {
	if (me.Value == nil) && ((me.ValueType == nil || string(*me.ValueType) != "arrayString") || (me.ValueType == nil)) {
		me.Value = new("")
	}
	if (me.Value != nil) && ((me.ValueType != nil && string(*me.ValueType) == "arrayString") && (me.ValueType != nil)) {
		return fmt.Errorf("'value' must not be specified unless ('value_type' is not set to 'arrayString' or 'value_type' is not set); got 'value_type'='%v'", opt.ValOrNil(me.ValueType))
	}
	// ---- MultiValue []string -> {"expectedValue":"arrayString","property":"valueType","type":"EQUALS"}
	return nil
}

func (me *FieldsAddAttributesEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"multi_value": &me.MultiValue,
		"name":        &me.Name,
		"value":       &me.Value,
		"value_type":  &me.ValueType,
	})
}
