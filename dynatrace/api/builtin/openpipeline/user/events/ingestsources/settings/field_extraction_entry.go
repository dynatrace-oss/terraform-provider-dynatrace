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

type FieldExtractionEntries []*FieldExtractionEntry

func (me *FieldExtractionEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_extraction_entry": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(FieldExtractionEntry).Schema()},
		},
	}
}

func (me FieldExtractionEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("field_extraction_entry", me)
}

func (me *FieldExtractionEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("field_extraction_entry", me); err != nil {
		return err
	}
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/895
	// Only known workaround is to ignore these blocks
	newEntries := FieldExtractionEntries{}
	var emptyItem FieldExtractionEntry
	for _, entry := range *me {
		if *entry != emptyItem {
			newEntries = append(newEntries, entry)
		}
	}
	*me = newEntries
	return nil
}

type FieldExtractionEntry struct {
	ConstantFieldName    *string                  `json:"constantFieldName,omitempty"`    // Destination field name
	ConstantValue        *string                  `json:"constantValue,omitempty"`        // Constant value to be assigned to field
	DefaultValue         *string                  `json:"defaultValue,omitempty"`         // Default value
	DestinationFieldName *string                  `json:"destinationFieldName,omitempty"` // Destination field name
	ExtractionType       FieldValueExtractionType `json:"extractionType"`                 // Field value extraction type. Possible values: `constant`, `field`
	SourceFieldName      *string                  `json:"sourceFieldName,omitempty"`      // Source field name
}

func (me *FieldExtractionEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"constant_field_name": {
			Type:        schema.TypeString,
			Description: "Destination field name",
			Optional:    true, // precondition
		},
		"constant_value": {
			Type:        schema.TypeString,
			Description: "Constant value to be assigned to field",
			Optional:    true, // precondition
		},
		"default_value": {
			Type:        schema.TypeString,
			Description: "Default value",
			Optional:    true, // nullable & precondition
		},
		"destination_field_name": {
			Type:        schema.TypeString,
			Description: "Destination field name",
			Optional:    true, // nullable & precondition
		},
		"extraction_type": {
			Type:        schema.TypeString,
			Description: "Field value extraction type. Possible values: `constant`, `field`",
			Optional:    true,
			Default:     "field", // new required attribute. Fallback to "field"
		},
		"source_field_name": {
			Type:        schema.TypeString,
			Description: "Source field name",
			Optional:    true, // precondition
		},
	}
}

func (me *FieldExtractionEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"constant_field_name":    me.ConstantFieldName,
		"constant_value":         me.ConstantValue,
		"default_value":          me.DefaultValue,
		"destination_field_name": me.DestinationFieldName,
		"extraction_type":        me.ExtractionType,
		"source_field_name":      me.SourceFieldName,
	})
}

func (me *FieldExtractionEntry) HandlePreconditions() error {
	if (me.ConstantFieldName == nil) && (string(me.ExtractionType) == "constant") {
		me.ConstantFieldName = opt.NewString("")
	}
	if (me.ConstantValue == nil) && (string(me.ExtractionType) == "constant") {
		me.ConstantValue = opt.NewString("")
	}
	if (me.SourceFieldName == nil) && (string(me.ExtractionType) == "field") {
		me.SourceFieldName = opt.NewString("")
	}
	if (me.DefaultValue != nil) && (string(me.ExtractionType) != "field") {
		return fmt.Errorf("'default_value' must not be specified if 'extraction_type' is set to '%v'", me.ExtractionType)
	}
	if (me.DestinationFieldName != nil) && (string(me.ExtractionType) != "field") {
		return fmt.Errorf("'destination_field_name' must not be specified if 'extraction_type' is set to '%v'", me.ExtractionType)
	}
	return nil
}

func (me *FieldExtractionEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"constant_field_name":    &me.ConstantFieldName,
		"constant_value":         &me.ConstantValue,
		"default_value":          &me.DefaultValue,
		"destination_field_name": &me.DestinationFieldName,
		"extraction_type":        &me.ExtractionType,
		"source_field_name":      &me.SourceFieldName,
	})
}
