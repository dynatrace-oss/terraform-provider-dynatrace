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

package maintenance

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Filters []*Filter

type Filter struct {
	EntityType      *string  `json:"entityType,omitempty"`      // Type of entities this maintenance window should match
	EntityId        *string  `json:"entityId,omitempty"`        // A specific entity that should match this maintenance window
	EntityTags      []string `json:"entityTags,omitempty"`      // The tags you want to use for matching in the format key or key:value
	ManagementZones []string `json:"managementZones,omitempty"` // The IDs of management zones to which the matched entities must belong
}

func (me *Filter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_type": {
			Type:        schema.TypeString,
			Description: "Type of entities this maintenance window should match",
			Optional:    true,
		},
		"entity_id": {
			Type:        schema.TypeString,
			Description: "A specific entity that should match this maintenance window",
			Optional:    true,
		},
		"entity_tags": {
			Type:        schema.TypeSet,
			Description: "The tags you want to use for matching in the format key or key:value",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"management_zones": {
			Type:        schema.TypeSet,
			Description: "The IDs of management zones to which the matched entities must belong",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Filter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"entity_type":      me.EntityType,
		"entity_id":        me.EntityId,
		"entity_tags":      me.EntityTags,
		"management_zones": me.ManagementZones,
	})
}

func (me *Filter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"entity_type":      &me.EntityType,
		"entity_id":        &me.EntityId,
		"entity_tags":      &me.EntityTags,
		"management_zones": &me.ManagementZones,
	})
}

func (me *Filters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A list of matching rules for dynamic filter formation.  If several rules are set, the OR logic applies",
			Elem:        &schema.Resource{Schema: new(Filter).Schema()},
		},
	}
}

func (me Filters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *Filters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter", me)
}
