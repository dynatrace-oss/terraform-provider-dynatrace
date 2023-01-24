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

type QueryDefinition struct {
	Type            Type             `json:"type"`                      // The type of query definition
	MetricSelector  *string          `json:"metricSelector,omitempty"`  /// To learn more, visit [Metric Selector](https://dt-url.net/metselad)
	MetricKey       string           `json:"metricKey"`                 // The metric key of the query definition
	Aggregation     *Aggregation     `json:"aggregation,omitempty"`     // The aggregation of the query definition
	QueryOffset     *int             `json:"queryOffset,omitempty"`     // Minute offset of sliding evaluation window for metrics with latency
	EntityFilter    *EntityFilter    `json:"entityFilter,omitempty"`    // Use rule-based filters to define the scope this event monitors.
	DimensionFilter DimensionFilters `json:"dimensionFilter,omitempty"` // The dimension filters of the query definition
}

func (me *QueryDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of query definition",
			Required:    true,
		},
		"metric_selector": {
			Type:             schema.TypeString,
			Description:      "To learn more, visit [Metric Selector](https://dt-url.net/metselad)",
			Optional:         true,
			DiffSuppressFunc: hcl.SuppressEOT,
		},
		"metric_key": {
			Type:        schema.TypeString,
			Description: "The metric key of the query definition",
			Required:    true,
		},
		"aggregation": {
			Type:        schema.TypeString,
			Description: "The aggregation of the query definition",
			Optional:    true,
		},
		"query_offset": {
			Type:        schema.TypeInt,
			Description: "Minute offset of sliding evaluation window for metrics with latency",
			Optional:    true,
		},
		"entity_filter": {
			Type:        schema.TypeList,
			Description: "Use rule-based filters to define the scope this event monitors.",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(EntityFilter).Schema()},
			Optional:    true,
		},
		"dimension_filter": {
			Type:        schema.TypeList,
			Description: "The dimension filters of the query definition",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DimensionFilters).Schema()},
			Optional:    true,
		},
	}
}

func (me *QueryDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":             me.Type,
		"metric_selector":  me.MetricSelector,
		"metric_key":       me.MetricKey,
		"aggregation":      me.Aggregation,
		"query_offset":     me.QueryOffset,
		"entity_filter":    me.EntityFilter,
		"dimension_filter": me.DimensionFilter,
	})
}

func (me *QueryDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":             &me.Type,
		"metric_selector":  &me.MetricSelector,
		"metric_key":       &me.MetricKey,
		"aggregation":      &me.Aggregation,
		"query_offset":     &me.QueryOffset,
		"entity_filter":    &me.EntityFilter,
		"dimension_filter": &me.DimensionFilter,
	})
}
