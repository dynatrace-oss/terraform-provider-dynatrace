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

type SmartscapeIdComponentsEntries []*SmartscapeIdComponentsEntry

func (me *SmartscapeIdComponentsEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id_component": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SmartscapeIdComponentsEntry).Schema()},
		},
	}
}

func (me SmartscapeIdComponentsEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("id_component", me)
}

func (me *SmartscapeIdComponentsEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("id_component", me)
}

type SmartscapeIdComponentsEntry struct {
	IdComponent         string `json:"idComponent"`         // ID component
	ReferencedFieldName string `json:"referencedFieldName"` // Referenced field name
}

func (me *SmartscapeIdComponentsEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id_component": {
			Type:        schema.TypeString,
			Description: "ID component",
			Required:    true,
		},
		"referenced_field_name": {
			Type:        schema.TypeString,
			Description: "Referenced field name",
			Required:    true,
		},
	}
}

func (me *SmartscapeIdComponentsEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id_component":          me.IdComponent,
		"referenced_field_name": me.ReferencedFieldName,
	})
}

func (me *SmartscapeIdComponentsEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id_component":          &me.IdComponent,
		"referenced_field_name": &me.ReferencedFieldName,
	})
}
