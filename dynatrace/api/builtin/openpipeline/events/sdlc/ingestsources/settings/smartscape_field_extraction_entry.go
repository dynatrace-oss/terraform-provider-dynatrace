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
	FieldName           string `json:"fieldName"`           // Field name
	ReferencedFieldName string `json:"referencedFieldName"` // Referenced field name
}

func (me *SmartscapeFieldExtractionEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_name": {
			Type:        schema.TypeString,
			Description: "Field name",
			Required:    true,
		},
		"referenced_field_name": {
			Type:        schema.TypeString,
			Description: "Referenced field name",
			Required:    true,
		},
	}
}

func (me *SmartscapeFieldExtractionEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"field_name":            me.FieldName,
		"referenced_field_name": me.ReferencedFieldName,
	})
}

func (me *SmartscapeFieldExtractionEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"field_name":            &me.FieldName,
		"referenced_field_name": &me.ReferencedFieldName,
	})
}
