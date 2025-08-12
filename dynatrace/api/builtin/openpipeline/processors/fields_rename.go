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

package processors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type FieldsRenameAttributes struct {
	Fields []*FieldsRenameItem `json:"fields"`
}

func (fra *FieldsRenameAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field": {
			Type:        schema.TypeList,
			Elem:        &schema.Resource{Schema: new(FieldsRenameItem).Schema()},
			Description: "Field to rename on the record",
			Required:    true,
			MinItems:    1,
			MaxItems:    50,
		},
	}
}

func (fra *FieldsRenameAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("field", fra.Fields)
}

func (fra *FieldsRenameAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("field", &fra.Fields)
}

type FieldsRenameItem struct {
	FromName string `json:"fromName"`
	ToName   string `json:"toName"`
}

const FieldsRenameMaxNameLength = 256

func (f *FieldsRenameItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from_name": {
			Type:         schema.TypeString,
			Description:  "The field to rename",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, FieldsRenameMaxNameLength),
		},
		"to_name": {
			Type:         schema.TypeString,
			Description:  "The new field name",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, FieldsRenameMaxNameLength),
		},
	}
}

func (f *FieldsRenameItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"from_name": f.FromName,
		"to_name":   f.ToName,
	})
}

func (f *FieldsRenameItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"from_name": &f.FromName,
		"to_name":   &f.ToName,
	})
}
