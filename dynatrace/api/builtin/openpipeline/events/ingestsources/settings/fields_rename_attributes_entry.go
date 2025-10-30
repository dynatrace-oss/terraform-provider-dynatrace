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

type FieldsRenameAttributesEntries []*FieldsRenameAttributesEntry

func (me *FieldsRenameAttributesEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(FieldsRenameAttributesEntry).Schema()},
		},
	}
}

func (me FieldsRenameAttributesEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("field", me)
}

func (me *FieldsRenameAttributesEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("field", me)
}

type FieldsRenameAttributesEntry struct {
	FromName string `json:"fromName"` // Fields's name
	ToName   string `json:"toName"`   // New field's name
}

func (me *FieldsRenameAttributesEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from_name": {
			Type:        schema.TypeString,
			Description: "Fields's name",
			Required:    true,
		},
		"to_name": {
			Type:        schema.TypeString,
			Description: "New field's name",
			Required:    true,
		},
	}
}

func (me *FieldsRenameAttributesEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"from_name": me.FromName,
		"to_name":   me.ToName,
	})
}

func (me *FieldsRenameAttributesEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"from_name": &me.FromName,
		"to_name":   &me.ToName,
	})
}
