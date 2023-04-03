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

package customunit

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Description string `json:"newUnitDescription"` // Unit description should provide additional information about the new unit
	Name        string `json:"newUnitName"`        // Unit name has to be unique and is used as identifier.
	PluralName  string `json:"newUnitPluralName"`  // Unit plural name represent the plural form of the unit name.
	Symbol      string `json:"newUnitSymbol"`      // Unit symbol has to be unique.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Unit description should provide additional information about the new unit",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Unit name has to be unique and is used as identifier.",
			Required:    true,
		},
		"plural_name": {
			Type:        schema.TypeString,
			Description: "Unit plural name represent the plural form of the unit name.",
			Required:    true,
		},
		"symbol": {
			Type:        schema.TypeString,
			Description: "Unit symbol has to be unique.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description": me.Description,
		"name":        me.Name,
		"plural_name": me.PluralName,
		"symbol":      me.Symbol,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description": &me.Description,
		"name":        &me.Name,
		"plural_name": &me.PluralName,
		"symbol":      &me.Symbol,
	})
}
