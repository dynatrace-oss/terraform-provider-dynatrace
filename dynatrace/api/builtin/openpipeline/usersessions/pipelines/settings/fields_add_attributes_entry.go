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
	Name  string `json:"name"`  // Fields's name
	Value string `json:"value"` // Field's value
}

func (me *FieldsAddAttributesEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Fields's name",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Field's value",
			Required:    true,
		},
	}
}

func (me *FieldsAddAttributesEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"value": me.Value,
	})
}

func (me *FieldsAddAttributesEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"value": &me.Value,
	})
}
