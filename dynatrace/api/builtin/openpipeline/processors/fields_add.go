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

type FieldsAddAttributes struct {
	Fields []*FieldsAddItem `json:"fields"`
}

func (faa *FieldsAddAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field": {
			Type:        schema.TypeList,
			Elem:        &schema.Resource{Schema: new(FieldsAddItem).Schema()},
			Description: "Field to add to the record",
			Required:    true,
			MinItems:    1,
			MaxItems:    50,
		},
	}
}

func (faa *FieldsAddAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("field", faa.Fields)
}

func (faa *FieldsAddAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("field", &faa.Fields)
}

type FieldsAddItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const FieldsAddMaxNameLength = 256
const FieldsAddMaxValueLength = 512

func (f *FieldsAddItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Description:  "Name of the field",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, FieldsAddMaxNameLength),
		},
		"value": {
			Type:         schema.TypeString,
			Description:  "Value to assign to the field",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(0, FieldsAddMaxValueLength),
		},
	}
}

func (f *FieldsAddItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  f.Name,
		"value": f.Value,
	})
}

func (f *FieldsAddItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &f.Name,
		"value": &f.Value,
	})
}
