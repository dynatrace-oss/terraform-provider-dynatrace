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

package log

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ColumnDefiningValue of a calculated log metric.
type ColumnDefiningValue struct {
	Name string `json:"name"` // The name of the column definiton
	Type string `json:"type"` // Defines the actual set of fields depending on the value. See one of the following objects: `CUSTOM` or `JSON`
}

func (me *ColumnDefiningValue) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the column definiton",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the actual set of fields depending on the value. See one of the following objects: `CUSTOM` or `JSON`",
		},
	}
}

func (me *ColumnDefiningValue) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name": me.Name,
		"type": me.Type,
	})
}

func (me *ColumnDefiningValue) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name": &me.Name,
		"type": &me.Type,
	})
}
