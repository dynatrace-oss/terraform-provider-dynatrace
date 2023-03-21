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

package usersettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Language Language `json:"language"`        // Possible Values: `Auto`, `En`, `Ja`
	Region   string   `json:"region"`          // Region
	Scope    string   `json:"-" scope:"scope"` // The scope of this setting (user, userdefaults)
	Theme    Theme    `json:"theme"`           // Possible Values: `Auto`, `Dark`, `Light`
	Timezone string   `json:"timezone"`        // Timezone
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"language": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Auto`, `En`, `Ja`",
			Required:    true,
		},
		"region": {
			Type:        schema.TypeString,
			Description: "Region",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (user, userdefaults)",
			Required:    true,
		},
		"theme": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Auto`, `Dark`, `Light`",
			Required:    true,
		},
		"timezone": {
			Type:        schema.TypeString,
			Description: "Timezone",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"language": me.Language,
		"region":   me.Region,
		"scope":    me.Scope,
		"theme":    me.Theme,
		"timezone": me.Timezone,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"language": &me.Language,
		"region":   &me.Region,
		"scope":    &me.Scope,
		"theme":    &me.Theme,
		"timezone": &me.Timezone,
	})
}
