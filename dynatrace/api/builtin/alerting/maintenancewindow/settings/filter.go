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

package maintenancewindow

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Filters []*Filter

func (me *Filters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Filter).Schema()},
		},
	}
}

func (me Filters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *Filters) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("filter", me); err != nil {
		return err
	}
	// slice may contain empty values because of SDK bug
	filters := Filters{}
	for _, filter := range *me {
		// empty value
		if filter.EntityID == nil && len(filter.EntityTags) == 0 && filter.EntityType == nil && len(filter.ManagementZones) == 0 {
			continue
		}
		filters = append(filters, filter)
	}
	*me = filters
	return nil
}

// Filter. Configured values of one filter are evaluated together (**AND**).. The maintenance window is applied onto an entity if it matches all of the values of at least one filter.
type Filter struct {
	EntityID        *string  `json:"entityId,omitempty"`        // A specific entity that should match this maintenance window.. **Note**: If an entity type filter value is set, it must be equal to the type of the selected entity. Otherwise this maintenance window will not match.
	EntityTags      []string `json:"entityTags,omitempty"`      // Entities which contain all of the configured tags will match this maintenance window.
	EntityType      *string  `json:"entityType,omitempty"`      // Type of entities this maintenance window should match.. If no entity type is selected all entities regardless of the type will match.
	ManagementZones []string `json:"managementZones,omitempty"` // Entities which are part of all the configured management zones will match this maintenance window.
}

func (me *Filter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_id": {
			Type:        schema.TypeString,
			Description: "A specific entity that should match this maintenance window.. **Note**: If an entity type filter value is set, it must be equal to the type of the selected entity. Otherwise this maintenance window will not match.",
			Optional:    true, // nullable
		},
		"entity_tags": {
			Type:        schema.TypeSet,
			Description: "Entities which contain all of the configured tags will match this maintenance window.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"entity_type": {
			Type:        schema.TypeString,
			Description: "Type of entities this maintenance window should match.. If no entity type is selected all entities regardless of the type will match.",
			Optional:    true, // nullable
		},
		"management_zones": {
			Type:        schema.TypeSet,
			Description: "Entities which are part of all the configured management zones will match this maintenance window.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Filter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"entity_id":        me.EntityID,
		"entity_tags":      me.EntityTags,
		"entity_type":      me.EntityType,
		"management_zones": me.ManagementZones,
	})
}

func (me *Filter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"entity_id":        &me.EntityID,
		"entity_tags":      &me.EntityTags,
		"entity_type":      &me.EntityType,
		"management_zones": &me.ManagementZones,
	})
}
