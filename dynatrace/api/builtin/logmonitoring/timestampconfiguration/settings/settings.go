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

package timestampconfiguration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Config_item_title   string             `json:"config-item-title"`            // Name
	Date_search_limit   *int               `json:"date-search-limit,omitempty"`  // Defines the number of characters in every log line (starting from the first character in the line) where the timestamp is searched.
	Date_time_pattern   string             `json:"date-time-pattern"`            // Date-time pattern
	Enabled             bool               `json:"enabled"`                      // This setting is enabled (`true`) or disabled (`false`)
	Entry_boundary      *EntryBoundary     `json:"entry-boundary,omitempty"`     // Optional field. Enter a fragment of the line text that starts the entry. No support for wildcards - the text is treated literally.
	Json_configuration  *JSONConfiguration `json:"json-configuration,omitempty"` // Detect JSON format
	Matchers            Matchers           `json:"matchers,omitempty"`
	Scope               *string            `json:"-" scope:"scope"`               // The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.
	Skip_indented_lines *bool              `json:"skip-indented-lines,omitempty"` // Don't parse timestamps in lines starting with white character
	Timezone            string             `json:"timezone"`                      // Timezone
	InsertAfter         string             `json:"-"`
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
		"date_search_limit": {
			Type:        schema.TypeInt,
			Description: "Defines the number of characters in every log line (starting from the first character in the line) where the timestamp is searched.",
			Optional:    true, // nullable
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
		"entry_boundary": {
			Type:        schema.TypeList,
			Description: "Optional field. Enter a fragment of the line text that starts the entry. No support for wildcards - the text is treated literally.",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(EntryBoundary).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"json_configuration": {
			Type:        schema.TypeList,
			Description: "Detect JSON format",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(JSONConfiguration).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"matchers": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Matchers).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"skip_indented_lines": {
			Type:        schema.TypeBool,
			Description: "Don't parse timestamps in lines starting with white character",
			Optional:    true, // nullable
		},
		"timezone": {
			Type:        schema.TypeString,
			Description: "Timezone",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"config_item_title":   me.Config_item_title,
		"date_search_limit":   me.Date_search_limit,
		"date_time_pattern":   me.Date_time_pattern,
		"enabled":             me.Enabled,
		"entry_boundary":      me.Entry_boundary,
		"json_configuration":  me.Json_configuration,
		"matchers":            me.Matchers,
		"scope":               me.Scope,
		"skip_indented_lines": me.Skip_indented_lines,
		"timezone":            me.Timezone,
		"insert_after":        me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"config_item_title":   &me.Config_item_title,
		"date_search_limit":   &me.Date_search_limit,
		"date_time_pattern":   &me.Date_time_pattern,
		"enabled":             &me.Enabled,
		"entry_boundary":      &me.Entry_boundary,
		"json_configuration":  &me.Json_configuration,
		"matchers":            &me.Matchers,
		"scope":               &me.Scope,
		"skip_indented_lines": &me.Skip_indented_lines,
		"timezone":            &me.Timezone,
		"insert_after":        &me.InsertAfter,
	})
}
