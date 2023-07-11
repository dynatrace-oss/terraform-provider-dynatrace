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

package workflows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DavisProblemConfig struct {
	EntityTagsMatch EntityTagsMatch         `json:"entityTagsMatch"`                // Possible values: `all` and `any`
	EntityTags      map[string]string       `json:"entityTags"`                     //
	OnProblemClose  bool                    `json:"onProblemClose" default:"false"` //
	Categories      *DavisProblemCategories `json:"categories"`                     //
	CustomFilter    string                  `json:"customFilter,omitempty"`         //
}

func (me *DavisProblemConfig) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_tags_match": {
			Type:         schema.TypeString,
			Description:  "Specifies whether all or just any of the configured entity tags need to match. Possible values: `all` and `any`. Omit this attribute if all entities should match",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
			RequiredWith: []string{prefix + ".0.entity_tags"},
		},
		"entity_tags": {
			Type:         schema.TypeMap,
			Description:  "key/value pairs for entity tags to match for. For tags that don't require a value, just specify an empty string as value. Omit this attribute if all entities should match",
			Optional:     true,
			Elem:         &schema.Schema{Type: schema.TypeString},
			RequiredWith: []string{prefix + ".0.entity_tags_match"},
		},
		"on_problem_close": {
			Type:        schema.TypeBool,
			Description: "If set to `true` closing a problem also is considered an event that triggers the execution",
			Optional:    true,
			Default:     false,
		},
		"categories": {
			Type:        schema.TypeList,
			Description: "",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DavisProblemCategories).Schema(prefix + ".0.categories")},
		},
		"custom_filter": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *DavisProblemConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"entity_tags_match": me.EntityTagsMatch,
		"entity_tags":       me.EntityTags,
		"on_problem_close":  me.OnProblemClose,
		"categories":        me.Categories,
		"custom_filter":     me.CustomFilter,
	})
}

func (me *DavisProblemConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"entity_tags_match": &me.EntityTagsMatch,
		"entity_tags":       &me.EntityTags,
		"on_problem_close":  &me.OnProblemClose,
		"categories":        &me.Categories,
		"custom_filter":     &me.CustomFilter,
	})
}
