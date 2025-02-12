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

package sitereliabilityguardian

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DisplayUnit struct {
	Base     string `json:"base"`     // Base Unit
	Decimals int    `json:"decimals"` // Decimals
	Display  string `json:"display"`  // display as unit
}

func (me *DisplayUnit) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"base": {
			Type:        schema.TypeString,
			Description: "Base Unit",
			Required:    true,
		},
		"decimals": {
			Type:        schema.TypeInt,
			Description: "Decimals",
			Required:    true,
		},
		"display": {
			Type:        schema.TypeString,
			Description: "display as unit",
			Required:    true,
		},
	}
}

func (me *DisplayUnit) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"base":     me.Base,
		"decimals": me.Decimals,
		"display":  me.Display,
	})
}

func (me *DisplayUnit) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"base":     &me.Base,
		"decimals": &me.Decimals,
		"display":  &me.Display,
	})
}
