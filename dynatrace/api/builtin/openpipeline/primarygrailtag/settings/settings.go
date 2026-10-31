/**
* @license
* Copyright 2026 Dynatrace LLC
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

package primarygrailtag

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled             bool         `json:"enabled"`             // Whether this rule is active
	KeepFields          bool         `json:"keepFields"`          // Whether to keep (true) or drop (false) the source field
	MatchingCondition   string       `json:"matchingCondition"`   // A DQL expression that determines which records this rule applies to
	PrimaryTagFieldName string       `json:"primaryTagFieldName"` // The unique target identifier within the primary_tags namespace. Must be unique across all rules.
	RuleName            string       `json:"ruleName"`            // A descriptive name for this rule
	SourceFields        SourceFields `json:"sourceFields"`        // Ordered list of source fields to consider. The first field that exist should be considered a source value for the primary tag
	InsertAfter         string       `json:"-"`
}

func (me *Settings) Name() string {
	return me.RuleName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Whether this rule is active",
			Required:    true,
		},
		"keep_fields": {
			Type:        schema.TypeBool,
			Description: "Whether to keep (true) or drop (false) the source field",
			Required:    true,
		},
		"matching_condition": {
			Type:        schema.TypeString,
			Description: "A DQL expression that determines which records this rule applies to",
			Required:    true,
		},
		"primary_tag_field_name": {
			Type:        schema.TypeString,
			Description: "The unique target identifier within the primary_tags namespace. Must be unique across all rules.",
			Required:    true,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "A descriptive name for this rule",
			Required:    true,
		},
		"source_fields": {
			Type:        schema.TypeList,
			Description: "Ordered list of source fields to consider. The first field that exist should be considered a source value for the primary tag",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SourceFields).Schema()},
			MinItems:    1,
			MaxItems:    1,
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
		"enabled":                me.Enabled,
		"keep_fields":            me.KeepFields,
		"matching_condition":     me.MatchingCondition,
		"primary_tag_field_name": me.PrimaryTagFieldName,
		"rule_name":              me.RuleName,
		"source_fields":          me.SourceFields,
		"insert_after":           me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":                &me.Enabled,
		"keep_fields":            &me.KeepFields,
		"matching_condition":     &me.MatchingCondition,
		"primary_tag_field_name": &me.PrimaryTagFieldName,
		"rule_name":              &me.RuleName,
		"source_fields":          &me.SourceFields,
		"insert_after":           &me.InsertAfter,
	})
}
