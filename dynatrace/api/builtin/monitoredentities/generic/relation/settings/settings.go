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

package relation

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CreatedBy      string        `json:"createdBy"`          // The user or extension that created this relationship.
	Enabled        bool          `json:"enabled"`            // This setting is enabled (`true`) or disabled (`false`)
	FromRole       *string       `json:"fromRole,omitempty"` // Specify a role for the source entity. If both source and destination type are the same, referring different roles will allow identification of a relationships direction. If role is left blank, any role of the source type is considered for the relationship.
	FromType       string        `json:"fromType"`           // Define an entity type as the source of the relationship.
	Sources        SourceFilters `json:"sources"`            // Specify all sources which should be evaluated for this relationship rule. The relationship is only created when any of the filters match.
	ToRole         *string       `json:"toRole,omitempty"`   // Specify a role for the destination entity. If both source and destination type are the same, referring different roles will allow identification of a relationships direction. If role is left blank, any role of the destination type is considered for the relationship.
	ToType         string        `json:"toType"`             // Define an entity type as the destination of the relationship. You can choose the same type as the source type. In this case you also may assign different roles for source and destination for having directed relationships.
	TypeOfRelation RelationType  `json:"typeOfRelation"`     // Possible Values: `CALLS`, `CHILD_OF`, `INSTANCE_OF`, `PART_OF`, `RUNS_ON`, `SAME_AS`
}

func (me *Settings) Name() string {
	return fmt.Sprintf("%s_%s_%s", me.FromType, me.TypeOfRelation, me.ToType)
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Type:        schema.TypeString,
			Description: "The user or extension that created this relationship.",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"from_role": {
			Type:        schema.TypeString,
			Description: "Specify a role for the source entity. If both source and destination type are the same, referring different roles will allow identification of a relationships direction. If role is left blank, any role of the source type is considered for the relationship.",
			Optional:    true,
		},
		"from_type": {
			Type:        schema.TypeString,
			Description: "Define an entity type as the source of the relationship.",
			Required:    true,
		},
		"sources": {
			Type:        schema.TypeList,
			Description: "Specify all sources which should be evaluated for this relationship rule. The relationship is only created when any of the filters match.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SourceFilters).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"to_role": {
			Type:        schema.TypeString,
			Description: "Specify a role for the destination entity. If both source and destination type are the same, referring different roles will allow identification of a relationships direction. If role is left blank, any role of the destination type is considered for the relationship.",
			Optional:    true,
		},
		"to_type": {
			Type:        schema.TypeString,
			Description: "Define an entity type as the destination of the relationship. You can choose the same type as the source type. In this case you also may assign different roles for source and destination for having directed relationships.",
			Required:    true,
		},
		"type_of_relation": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CALLS`, `CHILD_OF`, `INSTANCE_OF`, `PART_OF`, `RUNS_ON`, `SAME_AS`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"created_by":       me.CreatedBy,
		"enabled":          me.Enabled,
		"from_role":        me.FromRole,
		"from_type":        me.FromType,
		"sources":          me.Sources,
		"to_role":          me.ToRole,
		"to_type":          me.ToType,
		"type_of_relation": me.TypeOfRelation,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"created_by":       &me.CreatedBy,
		"enabled":          &me.Enabled,
		"from_role":        &me.FromRole,
		"from_type":        &me.FromType,
		"sources":          &me.Sources,
		"to_role":          &me.ToRole,
		"to_type":          &me.ToType,
		"type_of_relation": &me.TypeOfRelation,
	})
}
