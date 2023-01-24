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

type MetricEvents struct {
	Enabled                 bool             `json:"enabled"`                           // Enabled toggle of metric event entry
	Summary                 string           `json:"summary"`                           // The textual summary of the metric event entry
	QueryDefinition         *QueryDefinition `json:"queryDefinition"`                   // The query definition of the metric event entry
	ModelProperties         *ModelProperties `json:"modelProperties"`                   // The model properties of the metric event entry
	EventTemplate           *EventTemplate   `json:"eventTemplate"`                     // The event template of the metric event entry
	EventEntityDimensionKey *string          `json:"eventEntityDimensionKey,omitempty"` // Controls the preferred entity type used for triggered events.
	LegacyId                *string          `json:"legacyId,omitempty"`                // The legacy id of the metric event entry
}

func (me *MetricEvents) Name() string {
	return me.Summary
}

func (me *MetricEvents) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enabled toggle of metric event entry",
			Optional:    true,
		},
		"summary": {
			Type:        schema.TypeString,
			Description: "The textual summary of the metric event entry",
			Required:    true,
		},
		"query_definition": {
			Type:        schema.TypeList,
			Description: "The query definition of the metric event entry",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(QueryDefinition).Schema()},
			Required:    true,
		},
		"model_properties": {
			Type:        schema.TypeList,
			Description: "The model properties of the metric event entry",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ModelProperties).Schema()},
			Required:    true,
		},
		"event_template": {
			Type:        schema.TypeList,
			Description: "The event template of the metric event entry",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(EventTemplate).Schema()},
			Required:    true,
		},
		"event_entity_dimension_key": {
			Type:        schema.TypeString,
			Description: "Controls the preferred entity type used for triggered events.",
			Optional:    true,
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The legacy id of the metric event entry",
			Optional:    true,
		},
	}
}

func (me *MetricEvents) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":                    me.Enabled,
		"summary":                    me.Summary,
		"query_definition":           me.QueryDefinition,
		"model_properties":           me.ModelProperties,
		"event_template":             me.EventTemplate,
		"event_entity_dimension_key": me.EventEntityDimensionKey,
		"legacy_id":                  me.LegacyId,
	})
}

func (me *MetricEvents) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":                    &me.Enabled,
		"summary":                    &me.Summary,
		"query_definition":           &me.QueryDefinition,
		"model_properties":           &me.ModelProperties,
		"event_template":             &me.EventTemplate,
		"event_entity_dimension_key": &me.EventEntityDimensionKey,
		"legacy_id":                  &me.LegacyId,
	})
}
