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

package metricevents

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled                 bool             `json:"enabled"`                           // This setting is enabled (`true`) or disabled (`false`)
	EventEntityDimensionKey *string          `json:"eventEntityDimensionKey,omitempty"` // Controls the preferred entity type used for triggered events.
	EventTemplate           *EventTemplate   `json:"eventTemplate"`                     // Event template
	LegacyID                *string          `json:"legacyId,omitempty"`                // Config id
	ModelProperties         *ModelProperties `json:"modelProperties"`                   // Monitoring strategy
	QueryDefinition         *QueryDefinition `json:"queryDefinition"`                   // Query definition
	Summary                 string           `json:"summary"`                           // The textual summary of the metric event entry
}

func (me *Settings) Name() string {
	return me.Summary
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"event_entity_dimension_key": {
			Type:        schema.TypeString,
			Description: "Controls the preferred entity type used for triggered events.",
			Optional:    true, // nullable
		},
		"event_template": {
			Type:        schema.TypeList,
			Description: "Event template",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventTemplate).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "Config id",
			Computed:    true,
			Optional:    true, // nullable
		},
		"model_properties": {
			Type:        schema.TypeList,
			Description: "Monitoring strategy",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ModelProperties).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"query_definition": {
			Type:        schema.TypeList,
			Description: "Query definition",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(QueryDefinition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"summary": {
			Type:        schema.TypeString,
			Description: "The textual summary of the metric event entry",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":                    me.Enabled,
		"event_entity_dimension_key": me.EventEntityDimensionKey,
		"event_template":             me.EventTemplate,
		"legacy_id":                  me.LegacyID,
		"model_properties":           me.ModelProperties,
		"query_definition":           me.QueryDefinition,
		"summary":                    me.Summary,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":                    &me.Enabled,
		"event_entity_dimension_key": &me.EventEntityDimensionKey,
		"event_template":             &me.EventTemplate,
		"legacy_id":                  &me.LegacyID,
		"model_properties":           &me.ModelProperties,
		"query_definition":           &me.QueryDefinition,
		"summary":                    &me.Summary,
	})
}
