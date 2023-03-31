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

package dashboards

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GenericTagFilters []*GenericTagFilter

func (me *GenericTagFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(GenericTagFilter).Schema()},
		},
	}
}

func (me GenericTagFilters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *GenericTagFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter", me)
}

type GenericTagFilter struct {
	DisplayName               *string  `json:"displayName,omitempty"`               // The display name used to identify this generic filter
	TagKey                    *string  `json:"tagKey,omitempty"`                    // The tag key for this filter
	EntityTypes               []string `json:"entityTypes"`                         // Entity types affected by tag
	SuggestionsFromEntityType *string  `json:"suggestionsFromEntityType,omitempty"` // The entity type for which the suggestions should be provided.
}

func (me *GenericTagFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The display name used to identify this generic filter",
			Optional:    true,
		},
		"tag_key": {
			Type:        schema.TypeString,
			Description: "The tag key for this filter",
			Optional:    true,
		},
		"suggestions_from_entity_type": {
			Type:        schema.TypeString,
			Description: "The entity type for which the suggestions should be provided.",
			Optional:    true,
		},
		"entity_types": {
			Type:        schema.TypeSet,
			Description: "Entity types affected by tag",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
		},
	}
}

func (me *GenericTagFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":         me.DisplayName,
		"tag_key":      me.TagKey,
		"type":         me.SuggestionsFromEntityType,
		"entity_types": me.EntityTypes,
	})
}

func (me *GenericTagFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.DisplayName,
		"tag_key":      &me.TagKey,
		"type":         &me.SuggestionsFromEntityType,
		"entity_types": &me.EntityTypes,
	})
}
