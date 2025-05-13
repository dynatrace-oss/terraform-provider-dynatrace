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

package capturecustomproperties

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomProperties []*CustomProperty

func (me *CustomProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_session_properties_allow": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CustomProperty).Schema()},
		},
	}
}

func (me CustomProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("custom_session_properties_allow", me)
}

func (me *CustomProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("custom_session_properties_allow", me)
}

type CustomProperty struct {
	FieldDataType FieldDataType `json:"fieldDataType"` // Possible Values: `BOOLEAN`, `NUMBER`, `STRING`
	FieldName     string        `json:"fieldName"`     // Field name
}

func (me *CustomProperty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_data_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BOOLEAN`, `NUMBER`, `STRING`",
			Required:    true,
		},
		"field_name": {
			Type:        schema.TypeString,
			Description: "Field name",
			Required:    true,
		},
	}
}

func (me *CustomProperty) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"field_data_type": me.FieldDataType,
		"field_name":      me.FieldName,
	})
}

func (me *CustomProperty) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"field_data_type": &me.FieldDataType,
		"field_name":      &me.FieldName,
	})
}
