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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type QueryDefinition struct {
	Aggregation     *Aggregation     `json:"aggregation,omitempty"`     // Possible Values: `AVG`, `COUNT`, `MAX`, `MEDIAN`, `MIN`, `PERCENTILE90`, `SUM`, `VALUE`
	DimensionFilter DimensionFilters `json:"dimensionFilter,omitempty"` // Dimension filter
	EntityFilter    *EntityFilter    `json:"entityFilter,omitempty"`    // Use rule-based filters to define the scope this event monitors.
	ManagementZone  *string          `json:"managementZone,omitempty"`  // Management zone
	MetricKey       *string          `json:"metricKey,omitempty"`       // Metric key
	MetricSelector  *string          `json:"metricSelector,omitempty"`  // To learn more, visit [Metric Selector](https://dt-url.net/metselad)
	QueryOffset     *int             `json:"queryOffset,omitempty"`     // Minute offset of sliding evaluation window for metrics with latency
	Type            Type             `json:"type"`                      // Possible Values: `METRIC_KEY`, `METRIC_SELECTOR`
}

func (me *QueryDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aggregation": {
			Type:        schema.TypeString,
			Description: "Possible Values: `AVG`, `COUNT`, `MAX`, `MEDIAN`, `MIN`, `PERCENTILE90`, `SUM`, `VALUE`",
			Optional:    true, // precondition
		},
		"dimension_filter": {
			Type:        schema.TypeList,
			Description: "Dimension filter",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(DimensionFilters).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"entity_filter": {
			Type:        schema.TypeList,
			Description: "Use rule-based filters to define the scope this event monitors.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(EntityFilter).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"management_zone": {
			Type:        schema.TypeString,
			Description: "Management zone",
			Optional:    true, // nullable & precondition
		},
		"metric_key": {
			Type:        schema.TypeString,
			Description: "Metric key",
			Optional:    true, // precondition
		},
		"metric_selector": {
			Type:        schema.TypeString,
			Description: "To learn more, visit [Metric Selector](https://dt-url.net/metselad)",
			Optional:    true, // precondition
		},
		"query_offset": {
			Type:        schema.TypeInt,
			Description: "Minute offset of sliding evaluation window for metrics with latency",
			Optional:    true, // nullable
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `METRIC_KEY`, `METRIC_SELECTOR`",
			Required:    true,
		},
	}
}

func (me *QueryDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"aggregation":      me.Aggregation,
		"dimension_filter": me.DimensionFilter,
		"entity_filter":    me.EntityFilter,
		"management_zone":  me.ManagementZone,
		"metric_key":       me.MetricKey,
		"metric_selector":  me.MetricSelector,
		"query_offset":     me.QueryOffset,
		"type":             me.Type,
	})
}

func (me *QueryDefinition) HandlePreconditions() error {
	if me.Aggregation == nil && (string(me.Type) == "METRIC_KEY") {
		return fmt.Errorf("'aggregation' must be specified if 'type' is set to '%v'", me.Type)
	}
	if me.Aggregation != nil && (string(me.Type) != "METRIC_KEY") {
		return fmt.Errorf("'aggregation' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if me.MetricKey == nil && (string(me.Type) == "METRIC_KEY") {
		return fmt.Errorf("'metric_key' must be specified if 'type' is set to '%v'", me.Type)
	}
	if me.MetricSelector == nil && (string(me.Type) == "METRIC_SELECTOR") {
		return fmt.Errorf("'metric_selector' must be specified if 'type' is set to '%v'", me.Type)
	}
	// ---- DimensionFilter DimensionFilters -> {"preconditions":[{"expectedValue":"METRIC_KEY","property":"type","type":"EQUALS"},{"precondition":{"property":"metricKey","type":"NULL"},"type":"NOT"}],"type":"AND"}
	// ---- EntityFilter *EntityFilter -> {"preconditions":[{"expectedValue":"METRIC_KEY","property":"type","type":"EQUALS"},{"precondition":{"property":"metricKey","type":"NULL"},"type":"NOT"}],"type":"AND"}
	// ---- ManagementZone *string -> {"preconditions":[{"expectedValue":"METRIC_KEY","property":"type","type":"EQUALS"},{"expectedValue":"METRIC_SELECTOR","property":"type","type":"EQUALS"}],"type":"OR"}
	return nil
}

func (me *QueryDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"aggregation":      &me.Aggregation,
		"dimension_filter": &me.DimensionFilter,
		"entity_filter":    &me.EntityFilter,
		"management_zone":  &me.ManagementZone,
		"metric_key":       &me.MetricKey,
		"metric_selector":  &me.MetricSelector,
		"query_offset":     &me.QueryOffset,
		"type":             &me.Type,
	})
}
