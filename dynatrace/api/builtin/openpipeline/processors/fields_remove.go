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
)

type FieldsRemoveAttributes struct {
	Fields []string `json:"fields"`
}

func (fra *FieldsRemoveAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fields": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Fields to be removed from the record",
			Required:    true,
			MinItems:    1,
			MaxItems:    50,
		},
	}
}

func (fra *FieldsRemoveAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("fields", fra.Fields)
}

func (fra *FieldsRemoveAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("fields", &fra.Fields)
}
