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

type TransformationFields []*TransformationField

func (me *TransformationFields) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transformation_field": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(TransformationField).Schema()},
		},
	}
}

func (me TransformationFields) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("transformation_field", me)
}

func (me *TransformationFields) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("transformation_field", me)
}

type TransformationField struct {
	Array    bool                    `json:"array"` // Is Array
	Name     string                  `json:"name"`
	Optional bool                    `json:"optional"`
	Readonly bool                    `json:"readonly"` // Read-only
	Type     TransformationFieldType `json:"type"`     // Possible Values: `BOOLEAN`, `DOUBLE`, `DURATION`, `INT`, `IPADDR`, `LONG`, `STRING`, `TIMESTAMP`
}

func (me *TransformationField) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"array": {
			Type:        schema.TypeBool,
			Description: "Is Array",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"optional": {
			Type:        schema.TypeBool,
			Description: "no documentation available",
			Required:    true,
		},
		"readonly": {
			Type:        schema.TypeBool,
			Description: "Read-only",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BOOLEAN`, `DOUBLE`, `DURATION`, `INT`, `IPADDR`, `LONG`, `STRING`, `TIMESTAMP`",
			Required:    true,
		},
	}
}

func (me *TransformationField) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"array":    me.Array,
		"name":     me.Name,
		"optional": me.Optional,
		"readonly": me.Readonly,
		"type":     me.Type,
	})
}

func (me *TransformationField) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"array":    &me.Array,
		"name":     &me.Name,
		"optional": &me.Optional,
		"readonly": &me.Readonly,
		"type":     &me.Type,
	})
}
