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
	*me = hcl.FilterEmpty(*me, FieldExtractionEntry{})
	return nil
}

type FieldExtractionEntry struct {
	ConstantFieldName    *string                   `json:"constantFieldName,omitempty"`    // Destination field name
	ConstantValue        *string                   `json:"constantValue,omitempty"`        // Constant value to be assigned to field
	DefaultValue         *string                   `json:"defaultValue,omitempty"`         // Default value
	DestinationFieldName *string                   `json:"destinationFieldName,omitempty"` // Destination field name
	ExtractionType       *FieldValueExtractionType `json:"extractionType,omitempty"`       // Field value extraction type. Possible values: `constant`, `field`
	SourceFieldName      *string                   `json:"sourceFieldName,omitempty"`      // Source field name
	Strategy             *FieldExtractionStrategy  `json:"strategy,omitempty"`             // Strategy for field extraction. Possible values: `equals`, `startsWith`
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
			Optional:    true, // precondition
		},
		"source_field_name": {
			Type:        schema.TypeString,
			Description: "Source field name",
			Optional:    true, // precondition
		},
		"strategy": {
			Type:        schema.TypeString,
			Description: "Strategy for field extraction. Possible values: `equals`, `startsWith`",
			Optional:    true, // nullable
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
		"strategy":               me.Strategy,
	})
}

func (me *FieldExtractionEntry) HandlePreconditions() error {
	empty := FieldExtractionEntry{}
	if *me == empty {
		// ignore empty item, else it will error on missing conditional required fields
		return nil
	}
	if (me.ConstantFieldName == nil) && (((me.Strategy != nil && string(*me.Strategy) == "equals") || (me.Strategy == nil)) && (me.ExtractionType != nil && string(*me.ExtractionType) == "constant")) {
		me.ConstantFieldName = new("")
	}
	if (me.ConstantValue == nil) && (((me.Strategy != nil && string(*me.Strategy) == "equals") || (me.Strategy == nil)) && (me.ExtractionType != nil && string(*me.ExtractionType) == "constant")) {
		me.ConstantValue = new("")
	}
	if (me.SourceFieldName == nil) && ((me.ExtractionType != nil && string(*me.ExtractionType) == "field") || (me.Strategy != nil && string(*me.Strategy) == "startsWith")) {
		me.SourceFieldName = new("")
	}
	if (me.ConstantFieldName != nil) && (((me.Strategy == nil || string(*me.Strategy) != "equals") && (me.Strategy != nil)) || (me.ExtractionType == nil || string(*me.ExtractionType) != "constant")) {
		return fmt.Errorf("'constant_field_name' must not be specified unless (('strategy' is set to 'equals' or 'strategy' is not set) and 'extraction_type' is set to 'constant'); got 'strategy'='%v', 'extraction_type'='%v'", opt.ValOrNil(me.Strategy), opt.ValOrNil(me.ExtractionType))
	}
	if (me.ConstantValue != nil) && (((me.Strategy == nil || string(*me.Strategy) != "equals") && (me.Strategy != nil)) || (me.ExtractionType == nil || string(*me.ExtractionType) != "constant")) {
		return fmt.Errorf("'constant_value' must not be specified unless (('strategy' is set to 'equals' or 'strategy' is not set) and 'extraction_type' is set to 'constant'); got 'strategy'='%v', 'extraction_type'='%v'", opt.ValOrNil(me.Strategy), opt.ValOrNil(me.ExtractionType))
	}
	if (me.DefaultValue != nil) && (((me.Strategy == nil || string(*me.Strategy) != "equals") && (me.Strategy != nil)) || (me.ExtractionType == nil || string(*me.ExtractionType) != "field")) {
		return fmt.Errorf("'default_value' must not be specified unless (('strategy' is set to 'equals' or 'strategy' is not set) and 'extraction_type' is set to 'field'); got 'strategy'='%v', 'extraction_type'='%v'", opt.ValOrNil(me.Strategy), opt.ValOrNil(me.ExtractionType))
	}
	if (me.DestinationFieldName != nil) && (((me.Strategy == nil || string(*me.Strategy) != "equals") && (me.Strategy != nil)) || (me.ExtractionType == nil || string(*me.ExtractionType) != "field")) {
		return fmt.Errorf("'destination_field_name' must not be specified unless (('strategy' is set to 'equals' or 'strategy' is not set) and 'extraction_type' is set to 'field'); got 'strategy'='%v', 'extraction_type'='%v'", opt.ValOrNil(me.Strategy), opt.ValOrNil(me.ExtractionType))
	}
	if (me.ExtractionType != nil) && ((me.Strategy == nil || string(*me.Strategy) != "equals") && (me.Strategy != nil)) {
		return fmt.Errorf("'extraction_type' must not be specified unless ('strategy' is set to 'equals' or 'strategy' is not set); got 'strategy'='%v'", opt.ValOrNil(me.Strategy))
	}
	if (me.ExtractionType == nil) && ((me.Strategy != nil && string(*me.Strategy) == "equals") || (me.Strategy == nil)) {
		return fmt.Errorf("'extraction_type' must be specified when ('strategy' is set to 'equals' or 'strategy' is not set); got 'strategy'='%v'", opt.ValOrNil(me.Strategy))
	}
	if (me.SourceFieldName != nil) && ((me.ExtractionType == nil || string(*me.ExtractionType) != "field") && (me.Strategy == nil || string(*me.Strategy) != "startsWith")) {
		return fmt.Errorf("'source_field_name' must not be specified unless ('extraction_type' is set to 'field' or 'strategy' is set to 'startsWith'); got 'extraction_type'='%v', 'strategy'='%v'", opt.ValOrNil(me.ExtractionType), opt.ValOrNil(me.Strategy))
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
		"strategy":               &me.Strategy,
	})
}
