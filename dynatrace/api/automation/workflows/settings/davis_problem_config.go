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
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DavisProblemConfig struct {
	EntityTagsMatch *EntityTagsMatch        `json:"entityTagsMatch"`                // Possible values: `all` and `any`
	EntityTags      map[string]StringArray  `json:"entityTags"`                     // key/value pairs for entity tags to match for. For tags that don't require a value, just specify an empty string as value. Multiple values can be provided separated by whitespace (e.g. \"val1 val2\") and will be parsed as multiple tag values. Omit this attribute if all entities should match
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
			Description:  "key/value pairs for entity tags to match for. For tags that don't require a value, just specify an empty string as value. Multiple values can be provided separated by whitespace (e.g. \"val1 val2\") and will be parsed as multiple tag values. Omit this attribute if all entities should match",
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
			Description: "Additional DQL matcher expression to further filter events to match",
			Optional:    true,
		},
	}
}

func (me *DavisProblemConfig) MarshalHCL(properties hcl.Properties) error {
	if err := me.MarshalEntityTagsHCL(properties); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"entity_tags_match": me.EntityTagsMatch,
		"on_problem_close":  me.OnProblemClose,
		"categories":        me.Categories,
		"custom_filter":     me.CustomFilter,
	})
}

func (me *DavisProblemConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := me.UnmarshalEntityTagsHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"entity_tags_match": &me.EntityTagsMatch,
		"on_problem_close":  &me.OnProblemClose,
		"categories":        &me.Categories,
		"custom_filter":     &me.CustomFilter,
	})
}

func (me *DavisProblemConfig) MarshalEntityTagsHCL(properties hcl.Properties) error {
	entityTagsMap := map[string]string{}
	for k, v := range me.EntityTags {
		if len(k) == 0 {
			continue
		}
		if len(v) == 0 {
			continue
		}
		entityTagsMap[k] = strings.Join([]string(v), " ")
	}
	if len(entityTagsMap) > 0 {
		if err := properties.Encode("entity_tags", entityTagsMap); err != nil {
			return err
		}
	}
	return nil
}

func (me *DavisProblemConfig) UnmarshalEntityTagsHCL(decoder hcl.Decoder) error {
	entityTagsMap := map[string]string{}
	if err := decoder.Decode("entity_tags", &entityTagsMap); err != nil {
		return err
	}
	for k, v := range entityTagsMap {
		if len(k) == 0 {
			continue
		}
		if me.EntityTags == nil {
			me.EntityTags = map[string]StringArray{}
		}
		parts := strings.Split(v, " ")
		var sa StringArray
		for _, p := range parts {
			p = strings.TrimSpace(p)
			sa = append(sa, p)
		}
		me.EntityTags[k] = sa
	}
	return nil
}
