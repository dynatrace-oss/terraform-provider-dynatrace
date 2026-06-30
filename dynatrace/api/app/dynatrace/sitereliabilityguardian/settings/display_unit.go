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

// Display Unit. Defines unit conversion and decimal formatting applied when showing a DQL result in the UI.
type DisplayUnit struct {
	Base     string `json:"base"`     // Unit the DQL query returns its result in. Source unit for conversion.
	Decimals int    `json:"decimals"` // Number of decimal places (0-4) used when formatting the displayed value.
	Display  string `json:"display"`  // Unit to display the value in after conversion. Use Default to show the base unit as-is.
}

func (me *DisplayUnit) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"base": {
			Type:        schema.TypeString,
			Description: "Unit the DQL query returns its result in. Source unit for conversion.",
			Required:    true,
		},
		"decimals": {
			Type:        schema.TypeInt,
			Description: "Number of decimal places (0-4) used when formatting the displayed value.",
			Required:    true,
		},
		"display": {
			Type:        schema.TypeString,
			Description: "Unit to display the value in after conversion. Use Default to show the base unit as-is.",
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
