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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FieldExtractionEntries []*FieldExtractionEntry

func (me *FieldExtractionEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dimension": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(FieldExtractionEntry).Schema()},
		},
	}
}

func (me FieldExtractionEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("dimension", me)
}

func (me *FieldExtractionEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("dimension", me)
}

type FieldExtractionEntry struct {
	DefaultValue         *string `json:"defaultValue,omitempty"`         // Default value
	DestinationFieldName *string `json:"destinationFieldName,omitempty"` // Destination field name
	SourceFieldName      string  `json:"sourceFieldName"`                // Source field name
}

func (me *FieldExtractionEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_value": {
			Type:        schema.TypeString,
			Description: "Default value",
			Optional:    true, // nullable
		},
		"destination_field_name": {
			Type:        schema.TypeString,
			Description: "Destination field name",
			Optional:    true, // nullable
		},
		"source_field_name": {
			Type:        schema.TypeString,
			Description: "Source field name",
			Required:    true,
		},
	}
}

func (me *FieldExtractionEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"default_value":          me.DefaultValue,
		"destination_field_name": me.DestinationFieldName,
		"source_field_name":      me.SourceFieldName,
	})
}

func (me *FieldExtractionEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"default_value":          &me.DefaultValue,
		"destination_field_name": &me.DestinationFieldName,
		"source_field_name":      &me.SourceFieldName,
	})
}
