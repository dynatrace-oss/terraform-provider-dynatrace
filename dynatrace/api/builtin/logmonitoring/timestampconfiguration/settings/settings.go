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

package timestampconfiguration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Config_item_title string   `json:"config-item-title"` // Name
	Date_time_pattern string   `json:"date-time-pattern"` // Date-time pattern
	Enabled           bool     `json:"enabled"`           // This setting is enabled (`true`) or disabled (`false`)
	Matchers          Matchers `json:"matchers"`
	Scope             *string  `json:"-" scope:"scope"` // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	Timezone          string   `json:"timezone"`        // Timezone
}

func (me *Settings) Name() string {
	return *me.Scope + "_" + me.Config_item_title
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"config_item_title": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"date_time_pattern": {
			Type:        schema.TypeString,
			Description: "Date-time pattern",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"matchers": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(Matchers).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
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
		"config_item_title": me.Config_item_title,
		"date_time_pattern": me.Date_time_pattern,
		"enabled":           me.Enabled,
		"matchers":          me.Matchers,
		"scope":             me.Scope,
		"timezone":          me.Timezone,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"config_item_title": &me.Config_item_title,
		"date_time_pattern": &me.Date_time_pattern,
		"enabled":           &me.Enabled,
		"matchers":          &me.Matchers,
		"scope":             &me.Scope,
		"timezone":          &me.Timezone,
	})
}
