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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomFilterChartSeriesConfig is the configuration of a charted metric
type CustomFilterChartSeriesConfig struct {
	Metric          string                                    `json:"metric"`               // The name of the charted metric
	Aggregation     Aggregation                               `json:"aggregation"`          // The charted aggregation of the metric
	Percentile      *int64                                    `json:"percentile,omitempty"` // The charted percentile. Only applicable if the **aggregation** is set to `PERCENTILE`.
	Type            CustomFilterChartSeriesConfigType         `json:"type"`                 // The visualization of the timeseries chart
	EntityType      string                                    `json:"entityType"`           // The type of the Dynatrace entity that delivered the charted metric
	Dimensions      []*CustomFilterChartSeriesDimensionConfig `json:"dimensions"`           // Configuration of the charted metric splitting
	SortAscending   bool                                      `json:"sortAscending"`        // Sort ascending (`true`) or descending (`false`)
	SortColumn      bool                                      `json:"sortColumn"`
	AggregationRate *AggregationRate                          `json:"aggregationRate,omitempty"`
	Unknowns        map[string]json.RawMessage                `json:"-"`
}

func (me *CustomFilterChartSeriesConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric": {
			Type:        schema.TypeString,
			Description: "The name of the charted metric",
			Required:    true,
		},
		"aggregation": {
			Type:        schema.TypeString,
			Description: "The charted aggregation of the metric",
			Required:    true,
		},
		"percentile": {
			Type:        schema.TypeInt,
			Description: "The charted percentile. Only applicable if the **aggregation** is set to `PERCENTILE`",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The visualization of the timeseries chart. Possible values are `AREA`, `BAR` and `LINE`.",
			Required:    true,
		},
		"entity_type": {
			Type:        schema.TypeString,
			Description: "The visualization of the timeseries chart",
			Required:    true,
		},
		"dimension": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Configuration of the charted metric splitting",
			Elem: &schema.Resource{
				Schema: new(CustomFilterChartSeriesDimensionConfig).Schema(),
			},
		},
		"sort_ascending": {
			Type:        schema.TypeBool,
			Description: "Sort ascending (`true`) or descending (`false`)",
			Optional:    true,
		},
		"sort_column": {
			Type:        schema.TypeBool,
			Description: "Sort the column (`true`) or (`false`)",
			Optional:    true,
		},
		"aggregation_rate": {
			Type:        schema.TypeString,
			Description: "",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomFilterChartSeriesConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("metric", me.Metric); err != nil {
		return err
	}
	if err := properties.Encode("aggregation", string(me.Aggregation)); err != nil {
		return err
	}
	if err := properties.Encode("percentile", int(opt.Int64(me.Percentile))); err != nil {
		return err
	}
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("entity_type", me.EntityType); err != nil {
		return err
	}
	if err := properties.Encode("sort_ascending", me.SortAscending); err != nil {
		return err
	}
	if err := properties.Encode("sort_column", me.SortColumn); err != nil {
		return err
	}
	if err := properties.Encode("aggregation_rate", me.AggregationRate); err != nil {
		return err
	}
	if err := properties.Encode("dimension", me.Dimensions); err != nil {
		return err
	}
	return nil
}

func (me *CustomFilterChartSeriesConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "metric")
		delete(me.Unknowns, "aggregation")
		delete(me.Unknowns, "percentile")
		delete(me.Unknowns, "type")
		delete(me.Unknowns, "entity_type")
		delete(me.Unknowns, "dimensions")
		delete(me.Unknowns, "sort_ascending")
		delete(me.Unknowns, "sort_column")
		delete(me.Unknowns, "aggregation_rate")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("metric"); ok {
		me.Metric = value.(string)
	}
	if value, ok := decoder.GetOk("aggregation"); ok {
		me.Aggregation = Aggregation(value.(string))
	}
	if value, ok := decoder.GetOk("percentile"); ok {
		me.Percentile = opt.NewInt64(int64(value.(int)))
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = CustomFilterChartSeriesConfigType(value.(string))
	}
	if value, ok := decoder.GetOk("entity_type"); ok {
		me.EntityType = value.(string)
	}
	if result, ok := decoder.GetOk("dimension.#"); ok {
		me.Dimensions = []*CustomFilterChartSeriesDimensionConfig{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(CustomFilterChartSeriesDimensionConfig)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "dimension", idx)); err != nil {
				return err
			}
			me.Dimensions = append(me.Dimensions, entry)
		}
	}
	if value, ok := decoder.GetOk("sort_ascending"); ok {
		me.SortAscending = value.(bool)
	}
	if value, ok := decoder.GetOk("sort_column"); ok {
		me.SortColumn = value.(bool)
	}
	if value, ok := decoder.GetOk("aggregation_rate"); ok {
		me.AggregationRate = AggregationRate(value.(string)).Ref()
	}
	return nil
}

func (me *CustomFilterChartSeriesConfig) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.Metric)
		if err != nil {
			return nil, err
		}
		m["metric"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Aggregation)
		if err != nil {
			return nil, err
		}
		m["aggregation"] = rawMessage
	}
	if me.Percentile != nil {
		rawMessage, err := json.Marshal(me.Percentile)
		if err != nil {
			return nil, err
		}
		m["percentile"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Type)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.EntityType)
		if err != nil {
			return nil, err
		}
		m["entityType"] = rawMessage
	}
	dimensions := me.Dimensions
	if dimensions == nil {
		dimensions = []*CustomFilterChartSeriesDimensionConfig{}
	}
	{
		rawMessage, err := json.Marshal(dimensions)
		if err != nil {
			return nil, err
		}
		m["dimensions"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.SortAscending)
		if err != nil {
			return nil, err
		}
		m["sortAscending"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.SortColumn)
		if err != nil {
			return nil, err
		}
		m["sortColumn"] = rawMessage
	}
	if me.AggregationRate != nil {
		rawMessage, err := json.Marshal(me.AggregationRate)
		if err != nil {
			return nil, err
		}
		m["aggregationRate"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *CustomFilterChartSeriesConfig) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["metric"]; found {
		if err := json.Unmarshal(v, &me.Metric); err != nil {
			return err
		}
	}
	if v, found := m["aggregation"]; found {
		if err := json.Unmarshal(v, &me.Aggregation); err != nil {
			return err
		}
	}
	if v, found := m["percentile"]; found {
		if err := json.Unmarshal(v, &me.Percentile); err != nil {
			return err
		}
	}
	if v, found := m["type"]; found {
		if err := json.Unmarshal(v, &me.Type); err != nil {
			return err
		}
	}
	if v, found := m["entityType"]; found {
		if err := json.Unmarshal(v, &me.EntityType); err != nil {
			return err
		}
	}
	if v, found := m["dimensions"]; found {
		if err := json.Unmarshal(v, &me.Dimensions); err != nil {
			return err
		}
	}
	if v, found := m["sortAscending"]; found {
		if err := json.Unmarshal(v, &me.SortAscending); err != nil {
			return err
		}
	}
	if v, found := m["sortColumn"]; found {
		if err := json.Unmarshal(v, &me.SortColumn); err != nil {
			return err
		}
	}
	if v, found := m["aggregationRate"]; found {
		if err := json.Unmarshal(v, &me.AggregationRate); err != nil {
			return err
		}
	}
	delete(m, "metric")
	delete(m, "aggregation")
	delete(m, "percentile")
	delete(m, "type")
	delete(m, "entityType")
	delete(m, "dimensions")
	delete(m, "sortAscending")
	delete(m, "sortColumn")
	delete(m, "aggregationRate")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// Aggregation has no documentation
type Aggregation string

// Aggregations offers the known enum values
var Aggregations = struct {
	Avg              Aggregation
	Count            Aggregation
	Distinct         Aggregation
	Fastest10percent Aggregation
	Max              Aggregation
	Median           Aggregation
	Min              Aggregation
	None             Aggregation
	OfInterestRatio  Aggregation
	OtherRatio       Aggregation
	Percentile       Aggregation
	PerMin           Aggregation
	Slowest10percent Aggregation
	Slowest5percent  Aggregation
	Sum              Aggregation
	SumDimensions    Aggregation
}{
	"AVG",
	"COUNT",
	"DISTINCT",
	"FASTEST10PERCENT",
	"MAX",
	"MEDIAN",
	"MIN",
	"NONE",
	"OF_INTEREST_RATIO",
	"OTHER_RATIO",
	"PERCENTILE",
	"PER_MIN",
	"SLOWEST10PERCENT",
	"SLOWEST5PERCENT",
	"SUM",
	"SUM_DIMENSIONS",
}

// CustomFilterChartSeriesConfigType has no documentation
type CustomFilterChartSeriesConfigType string

// CustomFilterChartSeriesConfigTypes offers the known enum values
var CustomFilterChartSeriesConfigTypes = struct {
	Area CustomFilterChartSeriesConfigType
	Bar  CustomFilterChartSeriesConfigType
	Line CustomFilterChartSeriesConfigType
}{
	"AREA",
	"BAR",
	"LINE",
}

// AggregationRate has no documentation
type AggregationRate string

func (me AggregationRate) Ref() *AggregationRate {
	return &me
}

// AggregationRates offers the known enum values
var AggregationRates = struct {
	Hour   AggregationRate
	Minute AggregationRate
	Second AggregationRate
	Total  AggregationRate
}{
	"HOUR",
	"MINUTE",
	"SECOND",
	"TOTAL",
}
