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

package log

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CalculatedLogMetric Descriptor of a calculated log metric.
type CalculatedLogMetric struct {
	Description         *string              `json:"description,omitempty"`         // Descriptor of a calculated log metric.
	MetricKey           string               `json:"metricKey"`                     // The unique key of the calculated log metric.
	Active              *bool                `json:"active,omitempty"`              // The metric is enabled (`true`) or disabled (`false`).
	DisplayName         string               `json:"displayName"`                   // The name of the metric, displayed in the UI.
	Unit                Unit                 `json:"unit"`                          // The unit of the metric. Possible values: [ AMPERE, BILLION, BIT, BIT_PER_HOUR, BIT_PER_MINUTE, BIT_PER_SECOND, BYTE, BYTE_PER_HOUR, BYTE_PER_MINUTE, BYTE_PER_SECOND, CORES, COUNT, DAY, DECIBEL_MILLI_WATT, GIBI_BYTE, GIBI_BYTE_PER_HOUR, GIBI_BYTE_PER_MINUTE, GIBI_BYTE_PER_SECOND, GIGA, GIGA_BYTE, GIGA_BYTE_PER_HOUR, GIGA_BYTE_PER_MINUTE, GIGA_BYTE_PER_SECOND, HERTZ, HOUR, KIBI_BYTE, KIBI_BYTE_PER_HOUR, KIBI_BYTE_PER_MINUTE, KIBI_BYTE_PER_SECOND, KILO, KILO_BYTE, KILO_BYTE_PER_HOUR, KILO_BYTE_PER_MINUTE, KILO_BYTE_PER_SECOND, KILO_METRE_PER_HOUR, MEBI_BYTE, MEBI_BYTE_PER_HOUR, MEBI_BYTE_PER_MINUTE, MEBI_BYTE_PER_SECOND, MEGA, MEGA_BYTE, MEGA_BYTE_PER_HOUR, MEGA_BYTE_PER_MINUTE, MEGA_BYTE_PER_SECOND, METRE_PER_HOUR, METRE_PER_SECOND, MICRO_SECOND, MILLION, MILLI_CORES, MILLI_SECOND, MILLI_SECOND_PER_MINUTE, MINUTE, MONTH, MSU, NANO_SECOND, NANO_SECOND_PER_MINUTE, NOT_APPLICABLE, PERCENT, PER_HOUR, PER_MINUTE, PER_SECOND, PIXEL, PROMILLE, RATIO, SECOND, STATE, TRILLION, UNSPECIFIED, VOLT, WATT, WEEK, YEAR ]
	UnitDisplayName     *string              `json:"unitDisplayName,omitempty"`     // The display name of the unit. Only applicable if the unit is set to UNSPECIFIED.
	SearchString        string               `json:"searchString"`                  // The pattern to look for in logs. Use the Dynatrace search query language to specify it. Quotes must be escaped. To return all results, leave the field blank.
	MetricValueType     MetricValueType      `json:"metricValueType"`               // The type of the metric data points calculation. For now the only allowed value is OCCURRENCES. Possible values: [ FP_COLUMN_AVG, FP_COLUMN_COUNT, FP_COLUMN_MAX, FP_COLUMN_MIN, FP_COLUMN_SUM, FP_COLUMN_TOP_X_AVG, FP_COLUMN_TOP_X_COUNT, FP_COLUMN_TOP_X_MAX, FP_COLUMN_TOP_X_MIN, FP_COLUMN_TOP_X_SUM, INT_COLUMN_AVG, INT_COLUMN_COUNT, INT_COLUMN_MAX, INT_COLUMN_MIN, INT_COLUMN_SUM, INT_COLUMN_TOP_X_AVG, INT_COLUMN_TOP_X_COUNT, INT_COLUMN_TOP_X_MAX, INT_COLUMN_TOP_X_MIN, INT_COLUMN_TOP_X_SUM, OCCURRENCES ]
	ColumnDefiningValue *ColumnDefiningValue `json:"columnDefiningValue,omitempty"` // Definition of numeric column.
	LogSourceFilters    LogSourceFilters     `json:"logSourceFilters"`              // A list of filters to define the logs to look into. If several filters are specified, the OR logic applies.
}

func (me *CalculatedLogMetric) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Descriptor of a calculated log metric.",
		},
		"metric_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The unique key of the calculated log metric.",
		},
		"active": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The metric is enabled (`true`) or disabled (`false`).",
		},
		"display_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the metric, displayed in the UI.",
		},
		"unit": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The unit of the metric. Possible values: [ AMPERE, BILLION, BIT, BIT_PER_HOUR, BIT_PER_MINUTE, BIT_PER_SECOND, BYTE, BYTE_PER_HOUR, BYTE_PER_MINUTE, BYTE_PER_SECOND, CORES, COUNT, DAY, DECIBEL_MILLI_WATT, GIBI_BYTE, GIBI_BYTE_PER_HOUR, GIBI_BYTE_PER_MINUTE, GIBI_BYTE_PER_SECOND, GIGA, GIGA_BYTE, GIGA_BYTE_PER_HOUR, GIGA_BYTE_PER_MINUTE, GIGA_BYTE_PER_SECOND, HERTZ, HOUR, KIBI_BYTE, KIBI_BYTE_PER_HOUR, KIBI_BYTE_PER_MINUTE, KIBI_BYTE_PER_SECOND, KILO, KILO_BYTE, KILO_BYTE_PER_HOUR, KILO_BYTE_PER_MINUTE, KILO_BYTE_PER_SECOND, KILO_METRE_PER_HOUR, MEBI_BYTE, MEBI_BYTE_PER_HOUR, MEBI_BYTE_PER_MINUTE, MEBI_BYTE_PER_SECOND, MEGA, MEGA_BYTE, MEGA_BYTE_PER_HOUR, MEGA_BYTE_PER_MINUTE, MEGA_BYTE_PER_SECOND, METRE_PER_HOUR, METRE_PER_SECOND, MICRO_SECOND, MILLION, MILLI_CORES, MILLI_SECOND, MILLI_SECOND_PER_MINUTE, MINUTE, MONTH, MSU, NANO_SECOND, NANO_SECOND_PER_MINUTE, NOT_APPLICABLE, PERCENT, PER_HOUR, PER_MINUTE, PER_SECOND, PIXEL, PROMILLE, RATIO, SECOND, STATE, TRILLION, UNSPECIFIED, VOLT, WATT, WEEK, YEAR ]",
		},
		"unit_display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The display name of the unit. Only applicable if the unit is set to UNSPECIFIED.",
		},
		"search_string": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The pattern to look for in logs. Use the Dynatrace search query language to specify it. Quotes must be escaped. To return all results, leave the field blank.",
		},
		"metric_value_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the metric data points calculation. For now the only allowed value is OCCURRENCES. Possible values: [ FP_COLUMN_AVG, FP_COLUMN_COUNT, FP_COLUMN_MAX, FP_COLUMN_MIN, FP_COLUMN_SUM, FP_COLUMN_TOP_X_AVG, FP_COLUMN_TOP_X_COUNT, FP_COLUMN_TOP_X_MAX, FP_COLUMN_TOP_X_MIN, FP_COLUMN_TOP_X_SUM, INT_COLUMN_AVG, INT_COLUMN_COUNT, INT_COLUMN_MAX, INT_COLUMN_MIN, INT_COLUMN_SUM, INT_COLUMN_TOP_X_AVG, INT_COLUMN_TOP_X_COUNT, INT_COLUMN_TOP_X_MAX, INT_COLUMN_TOP_X_MIN, INT_COLUMN_TOP_X_SUM, OCCURRENCES ]",
		},
		"column_defining_value": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Definition of numeric column.",
			Elem:        &schema.Resource{Schema: new(ColumnDefiningValue).Schema()},
		},
		"log_source_filters": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "A list of filters to define the logs to look into. If several filters are specified, the OR logic applies.",
			Elem:        &schema.Resource{Schema: new(LogSourceFilters).Schema()},
		},
	}
}

func (me *CalculatedLogMetric) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":           me.Description,
		"metric_key":            me.MetricKey,
		"active":                me.Active,
		"display_name":          me.DisplayName,
		"unit":                  me.Unit,
		"unit_display_name":     me.UnitDisplayName,
		"search_string":         me.SearchString,
		"metric_value_type":     me.MetricValueType,
		"column_defining_value": me.ColumnDefiningValue,
		"log_source_filters":    me.LogSourceFilters,
	})
}

func (me *CalculatedLogMetric) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description":           &me.Description,
		"metric_key":            &me.MetricKey,
		"active":                &me.Active,
		"display_name":          &me.DisplayName,
		"unit":                  &me.Unit,
		"unit_display_name":     &me.UnitDisplayName,
		"search_string":         &me.SearchString,
		"metric_value_type":     &me.MetricValueType,
		"column_defining_value": &me.ColumnDefiningValue,
		"log_source_filters":    &me.LogSourceFilters,
	})
}

type Unit string

var Units = struct {
	Ampere               Unit
	Billion              Unit
	Bit                  Unit
	BitPerHour           Unit
	BitPerMinute         Unit
	BitPerSecond         Unit
	Byte                 Unit
	BytePerHour          Unit
	BytePerMinute        Unit
	BytePerSecond        Unit
	Cores                Unit
	Count                Unit
	Day                  Unit
	DecibelMilliWatt     Unit
	GibiByte             Unit
	GibiBytePerHour      Unit
	GibiBytePerMinute    Unit
	GibiBytePerSecond    Unit
	Giga                 Unit
	GigaByte             Unit
	GigaBytePerHour      Unit
	GigaBytePerMinute    Unit
	GigaBytePerSecond    Unit
	Hertz                Unit
	Hour                 Unit
	KibiByte             Unit
	KibiBytePerHour      Unit
	KibiBytePerMinute    Unit
	KibiBytePerSecond    Unit
	Kilo                 Unit
	KiloByte             Unit
	KiloBytePerHour      Unit
	KiloBytePerMinute    Unit
	KiloBytePerSecond    Unit
	KiloMetrePerHour     Unit
	MebiByte             Unit
	MebiBytePerHour      Unit
	MebiBytePerMinute    Unit
	MebiBytePerSecond    Unit
	Mega                 Unit
	MegaByte             Unit
	MegaBytePerHour      Unit
	MegaBytePerMinute    Unit
	MegaBytePerSecond    Unit
	MetrePerHour         Unit
	MetrePerSecond       Unit
	MicroSecond          Unit
	Million              Unit
	MilliCores           Unit
	MilliSecond          Unit
	MilliSecondPerMinute Unit
	Minute               Unit
	Month                Unit
	MSU                  Unit
	NanoSecond           Unit
	NanoSecondPerMinute  Unit
	NotApplicable        Unit
	Percent              Unit
	PerHour              Unit
	PerMinute            Unit
	PerSecond            Unit
	Pixel                Unit
	Promille             Unit
	Ratio                Unit
	Second               Unit
	State                Unit
	Trillion             Unit
	Unspecified          Unit
	Volt                 Unit
	Watt                 Unit
	Week                 Unit
	Year                 Unit
}{
	"AMPERE",
	"BILLION",
	"BIT",
	"BIT_PER_HOUR",
	"BIT_PER_MINUTE",
	"BIT_PER_SECOND",
	"BYTE",
	"BYTE_PER_HOUR",
	"BYTE_PER_MINUTE",
	"BYTE_PER_SECOND",
	"CORES",
	"COUNT",
	"DAY",
	"DECIBEL_MILLI_WATT",
	"GIBI_BYTE",
	"GIBI_BYTE_PER_HOUR",
	"GIBI_BYTE_PER_MINUTE",
	"GIBI_BYTE_PER_SECOND",
	"GIGA",
	"GIGA_BYTE",
	"GIGA_BYTE_PER_HOUR",
	"GIGA_BYTE_PER_MINUTE",
	"GIGA_BYTE_PER_SECOND",
	"HERTZ",
	"HOUR",
	"KIBI_BYTE",
	"KIBI_BYTE_PER_HOUR",
	"KIBI_BYTE_PER_MINUTE",
	"KIBI_BYTE_PER_SECOND",
	"KILO",
	"KILO_BYTE",
	"KILO_BYTE_PER_HOUR",
	"KILO_BYTE_PER_MINUTE",
	"KILO_BYTE_PER_SECOND",
	"KILO_METRE_PER_HOUR",
	"MEBI_BYTE",
	"MEBI_BYTE_PER_HOUR",
	"MEBI_BYTE_PER_MINUTE",
	"MEBI_BYTE_PER_SECOND",
	"MEGA",
	"MEGA_BYTE",
	"MEGA_BYTE_PER_HOUR",
	"MEGA_BYTE_PER_MINUTE",
	"MEGA_BYTE_PER_SECOND",
	"METRE_PER_HOUR",
	"METRE_PER_SECOND",
	"MICRO_SECOND",
	"MILLION",
	"MILLI_CORES",
	"MILLI_SECOND",
	"MILLI_SECOND_PER_MINUTE",
	"MINUTE",
	"MONTH",
	"MSU",
	"NANO_SECOND",
	"NANO_SECOND_PER_MINUTE",
	"NOT_APPLICABLE",
	"PERCENT",
	"PER_HOUR",
	"PER_MINUTE",
	"PER_SECOND",
	"PIXEL",
	"PROMILLE",
	"RATIO",
	"SECOND",
	"STATE",
	"TRILLION",
	"UNSPECIFIED",
	"VOLT",
	"WATT",
	"WEEK",
	"YEAR",
}

type MetricValueType string

var MetricValueTypes = struct {
	FpColumnAvg        MetricValueType
	FpColumnCount      MetricValueType
	FpColumnMax        MetricValueType
	FpColumnMin        MetricValueType
	FpColumnSum        MetricValueType
	FpColumnTopXAvg    MetricValueType
	FpColumnTopXCount  MetricValueType
	FpColumnTopXMax    MetricValueType
	FpColumnTopXMin    MetricValueType
	FpColumnTopXSum    MetricValueType
	IntColumnAvg       MetricValueType
	IntColumnCount     MetricValueType
	IntColumnMax       MetricValueType
	IntColumnMin       MetricValueType
	IntColumnSum       MetricValueType
	IntColumnTopXAvg   MetricValueType
	IntColumnTopXCount MetricValueType
	IntColumnTopXMax   MetricValueType
	IntColumnTopXMin   MetricValueType
	IntColumnTopXSum   MetricValueType
	Occurrences        MetricValueType
}{
	"FP_COLUMN_AVG",
	"FP_COLUMN_COUNT",
	"FP_COLUMN_MAX",
	"FP_COLUMN_MIN",
	"FP_COLUMN_SUM",
	"FP_COLUMN_TOP_X_AVG",
	"FP_COLUMN_TOP_X_COUNT",
	"FP_COLUMN_TOP_X_MAX",
	"FP_COLUMN_TOP_X_MIN",
	"FP_COLUMN_TOP_X_SUM",
	"INT_COLUMN_AVG",
	"INT_COLUMN_COUNT",
	"INT_COLUMN_MAX",
	"INT_COLUMN_MIN",
	"INT_COLUMN_SUM",
	"INT_COLUMN_TOP_X_AVG",
	"INT_COLUMN_TOP_X_COUNT",
	"INT_COLUMN_TOP_X_MAX",
	"INT_COLUMN_TOP_X_MIN",
	"INT_COLUMN_TOP_X_SUM",
	"OCCURRENCES",
}
