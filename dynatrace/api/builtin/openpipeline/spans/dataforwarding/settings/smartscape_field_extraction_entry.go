/**
* @license
* Copyright 2026 Dynatrace LLC
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

package dataforwarding

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SmartscapeFieldExtractionEntries []*SmartscapeFieldExtractionEntry

func (me *SmartscapeFieldExtractionEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"smartscape_field_extraction_entry": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SmartscapeFieldExtractionEntry).Schema()},
		},
	}
}

func (me SmartscapeFieldExtractionEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("smartscape_field_extraction_entry", me)
}

func (me *SmartscapeFieldExtractionEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("smartscape_field_extraction_entry", me)
}

type SmartscapeFieldExtractionEntry struct {
	FieldName           *string                  `json:"fieldName,omitempty"` // Field name
	ReferencedFieldName string                   `json:"referencedFieldName"` // Referenced field name
	Strategy            *FieldExtractionStrategy `json:"strategy,omitempty"`  // Strategy for field extraction. Possible values: `equals`, `startsWith`
}

func (me *SmartscapeFieldExtractionEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_name": {
			Type:        schema.TypeString,
			Description: "Field name",
			Optional:    true, // precondition
		},
		"referenced_field_name": {
			Type:        schema.TypeString,
			Description: "Referenced field name",
			Required:    true,
		},
		"strategy": {
			Type:        schema.TypeString,
			Description: "Strategy for field extraction. Possible values: `equals`, `startsWith`",
			Optional:    true, // nullable
		},
	}
}

func (me *SmartscapeFieldExtractionEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"field_name":            me.FieldName,
		"referenced_field_name": me.ReferencedFieldName,
		"strategy":              me.Strategy,
	})
}

func (me *SmartscapeFieldExtractionEntry) HandlePreconditions() error {
	if (me.FieldName != nil) && ((me.Strategy == nil || string(*me.Strategy) != "equals") && (me.Strategy != nil)) {
		return fmt.Errorf("'field_name' must not be specified unless ('strategy' is set to 'equals' or 'strategy' is not set); got 'strategy'='%v'", opt.ValOrNil(me.Strategy))
	}
	if (me.FieldName == nil) && ((me.Strategy != nil && string(*me.Strategy) == "equals") || (me.Strategy == nil)) {
		return fmt.Errorf("'field_name' must be specified when ('strategy' is set to 'equals' or 'strategy' is not set); got 'strategy'='%v'", opt.ValOrNil(me.Strategy))
	}
	return nil
}

func (me *SmartscapeFieldExtractionEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"field_name":            &me.FieldName,
		"referenced_field_name": &me.ReferencedFieldName,
		"strategy":              &me.Strategy,
	})
}
