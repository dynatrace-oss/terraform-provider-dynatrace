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

// CustomFilterChartConfig Configuration of a custom chart
type CustomFilterChartConfig struct {
	LegendShown         *bool                                        `json:"legendShown,omitempty"` // Defines if a legend should be shown
	Type                CustomFilterChartConfigType                  `json:"type"`                  // The type of the chart
	Series              []*CustomFilterChartSeriesConfig             `json:"series"`                // A list of charted metrics
	ResultMetadata      map[string]*CustomChartingItemMetadataConfig `json:"resultMetadata"`        // Additional information about charted metric
	AxisLimits          map[string]float64                           `json:"axisLimits,omitempty"`  // The optional custom y-axis limits
	LeftAxisCustomUnit  *LeftAxisCustomUnit                          `json:"leftAxisCustomUnit,omitempty"`
	RightAxisCustomUnit *RightAxisCustomUnit                         `json:"rightAxisCustomUnit,omitempty"`
	Unknowns            map[string]json.RawMessage                   `json:"-"`
}

func (me *CustomFilterChartConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"legend": {
			Type:        schema.TypeBool,
			Description: "Defines if a legend should be shown",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the chart",
			Required:    true,
		},
		"series": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A list of charted metrics",
			Elem: &schema.Resource{
				Schema: new(CustomFilterChartSeriesConfig).Schema(),
			},
		},
		"result_metadata": {
			Type:        schema.TypeList,
			Description: "Additional information about charted metric",
			Optional:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: new(ResultMetadata).Schema(),
			},
		},
		"axis_limits": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The optional custom y-axis limits",
			Elem:        &schema.Schema{Type: schema.TypeFloat},
		},
		"left_axis_custom_unit": {
			Type:        schema.TypeString,
			Description: "Either one of `Bit`, `BitPerHour`, `BitPerMinute`, `BitPerSecond`, `Byte`, `BytePerHour`, `BytePerMinute`, `BytePerSecond`, `Cores`, `Count`, `Day`, `DecibelMilliWatt`, `GibiByte`, `Giga`, `GigaByte`, `Hour`, `KibiByte`, `KibiBytePerHour`, `KibiBytePerMinute`, `KibiBytePerSecond`, `Kilo`, `KiloByte`, `KiloBytePerHour`, `KiloBytePerMinute`, `KiloBytePerSecond`, `MebiByte`, `MebiBytePerHour`, `MebiBytePerMinute`, `MebiBytePerSecond`, `Mega`, `MegaByte`, `MegaBytePerHour`, `MegaBytePerMinute`, `MegaBytePerSecond`, `MicroSecond`, `MilliCores`, `MilliSecond`, `MilliSecondPerMinute`, `Minute`, `Month`, `NanoSecond`, `NanoSecondPerMinute`, `NotApplicable`, `PerHour`, `PerMinute`, `PerSecond`, `Percent`, `Pixel`, `Promille`, `Ratio`, `Second`, `State`, `Unspecified`, `Week`, `Year`",
			Optional:    true,
		},
		"right_axis_custom_unit": {
			Type:        schema.TypeString,
			Description: "Either one of `Bit`, `BitPerHour`, `BitPerMinute`, `BitPerSecond`, `Byte`, `BytePerHour`, `BytePerMinute`, `BytePerSecond`, `Cores`, `Count`, `Day`, `DecibelMilliWatt`, `GibiByte`, `Giga`, `GigaByte`, `Hour`, `KibiByte`, `KibiBytePerHour`, `KibiBytePerMinute`, `KibiBytePerSecond`, `Kilo`, `KiloByte`, `KiloBytePerHour`, `KiloBytePerMinute`, `KiloBytePerSecond`, `MebiByte`, `MebiBytePerHour`, `MebiBytePerMinute`, `MebiBytePerSecond`, `Mega`, `MegaByte`, `MegaBytePerHour`, `MegaBytePerMinute`, `MegaBytePerSecond`, `MicroSecond`, `MilliCores`, `MilliSecond`, `MilliSecondPerMinute`, `Minute`, `Month`, `NanoSecond`, `NanoSecondPerMinute`, `NotApplicable`, `PerHour`, `PerMinute`, `PerSecond`, `Percent`, `Pixel`, `Promille`, `Ratio`, `Second`, `State`, `Unspecified`, `Week`, `Year`",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomFilterChartConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("legend", opt.Bool(me.LegendShown)); err != nil {
		return err
	}
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("series", me.Series); err != nil {
		return err
	}
	if len(me.ResultMetadata) > 0 {
		resultMetadata := &ResultMetadata{Entries: []*ResultMetadataEntry{}}
		for k, v := range me.ResultMetadata {
			entry := &ResultMetadataEntry{
				Key:    k,
				Config: v,
			}
			resultMetadata.Entries = append(resultMetadata.Entries, entry)
		}
		if err := properties.Encode("result_metadata", resultMetadata); err != nil {
			return err
		}
	}
	if len(me.AxisLimits) > 0 {
		data, err := json.Marshal(me.AxisLimits)
		if err != nil {
			return err
		}
		if err := properties.Encode("axis_limits", string(data)); err != nil {
			return err
		}
	}
	if err := properties.Encode("left_axis_custom_unit", me.LeftAxisCustomUnit); err != nil {
		return err
	}
	if err := properties.Encode("right_axis_custom_unit", me.RightAxisCustomUnit); err != nil {
		return err
	}
	return nil
}

func (me *CustomFilterChartConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "legend")
		delete(me.Unknowns, "type")
		delete(me.Unknowns, "series")
		delete(me.Unknowns, "result_metadata")
		delete(me.Unknowns, "axis_limits")
		delete(me.Unknowns, "left_axis_custom_unit")
		delete(me.Unknowns, "right_axis_custom_unit")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("legend"); ok {
		me.LegendShown = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = CustomFilterChartConfigType(value.(string))
	}
	if result, ok := decoder.GetOk("series.#"); ok {
		me.Series = []*CustomFilterChartSeriesConfig{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(CustomFilterChartSeriesConfig)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "series", idx)); err != nil {
				return err
			}
			me.Series = append(me.Series, entry)
		}
	}
	if _, ok := decoder.GetOk("result_metadata.#"); ok {
		resultMetadata := new(ResultMetadata)
		if err := resultMetadata.UnmarshalHCL(hcl.NewDecoder(decoder, "result_metadata", 0)); err != nil {
			return err
		}
		if len(resultMetadata.Entries) > 0 {
			me.ResultMetadata = map[string]*CustomChartingItemMetadataConfig{}
			for _, entry := range resultMetadata.Entries {
				me.ResultMetadata[entry.Key] = entry.Config
			}
		}
	}
	if value, ok := decoder.GetOk("axis_limits"); ok {
		if err := json.Unmarshal([]byte(value.(string)), &me.AxisLimits); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("left_axis_custom_unit"); ok {
		me.LeftAxisCustomUnit = LeftAxisCustomUnit(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("right_axis_custom_unit"); ok {
		me.RightAxisCustomUnit = RightAxisCustomUnit(value.(string)).Ref()
	}
	return nil
}

func (me *CustomFilterChartConfig) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.LegendShown != nil {
		rawMessage, err := json.Marshal(me.LegendShown)
		if err != nil {
			return nil, err
		}
		m["legendShown"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Type)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}

	series := me.Series
	if series == nil {
		series = []*CustomFilterChartSeriesConfig{}
	}
	{
		rawMessage, err := json.Marshal(series)
		if err != nil {
			return nil, err
		}
		m["series"] = rawMessage
	}
	resultMetadata := me.ResultMetadata
	if resultMetadata == nil {
		resultMetadata = map[string]*CustomChartingItemMetadataConfig{}
	}
	{
		rawMessage, err := json.Marshal(resultMetadata)
		if err != nil {
			return nil, err
		}
		m["resultMetadata"] = rawMessage
	}
	if len(me.AxisLimits) > 0 {
		rawMessage, err := json.Marshal(me.AxisLimits)
		if err != nil {
			return nil, err
		}
		m["axisLimits"] = rawMessage
	}
	if me.LeftAxisCustomUnit != nil {
		rawMessage, err := json.Marshal(me.LeftAxisCustomUnit)
		if err != nil {
			return nil, err
		}
		m["leftAxisCustomUnit"] = rawMessage
	}
	if me.RightAxisCustomUnit != nil {
		rawMessage, err := json.Marshal(me.RightAxisCustomUnit)
		if err != nil {
			return nil, err
		}
		m["rightAxisCustomUnit"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *CustomFilterChartConfig) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["legendShown"]; found {
		if err := json.Unmarshal(v, &me.LegendShown); err != nil {
			return err
		}
	}
	if v, found := m["type"]; found {
		if err := json.Unmarshal(v, &me.Type); err != nil {
			return err
		}
	}
	if v, found := m["series"]; found {
		if err := json.Unmarshal(v, &me.Series); err != nil {
			return err
		}
	}
	if v, found := m["resultMetadata"]; found {
		if err := json.Unmarshal(v, &me.ResultMetadata); err != nil {
			return err
		}
	}
	if v, found := m["axisLimits"]; found {
		if err := json.Unmarshal(v, &me.AxisLimits); err != nil {
			return err
		}
	}
	if v, found := m["leftAxisCustomUnit"]; found {
		if err := json.Unmarshal(v, &me.LeftAxisCustomUnit); err != nil {
			return err
		}
	}
	if v, found := m["rightAxisCustomUnit"]; found {
		if err := json.Unmarshal(v, &me.RightAxisCustomUnit); err != nil {
			return err
		}
	}
	delete(m, "legendShown")
	delete(m, "type")
	delete(m, "series")
	delete(m, "resultMetadata")
	delete(m, "axisLimits")
	delete(m, "leftAxisCustomUnit")
	delete(m, "rightAxisCustomUnit")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// CustomFilterChartConfigType has no documentation
type CustomFilterChartConfigType string

// CustomFilterChartConfigTypes offers the known enum values
var CustomFilterChartConfigTypes = struct {
	Pie         CustomFilterChartConfigType
	SingleValue CustomFilterChartConfigType
	TimeSeries  CustomFilterChartConfigType
	TopList     CustomFilterChartConfigType
}{
	"PIE",
	"SINGLE_VALUE",
	"TIMESERIES",
	"TOP_LIST",
}

// LeftAxisCustomUnit has no documentation
type LeftAxisCustomUnit string

func (me LeftAxisCustomUnit) Ref() *LeftAxisCustomUnit {
	return &me
}

// LeftAxisCustomUnits has no documentation
var LeftAxisCustomUnits = struct {
	Bit                  LeftAxisCustomUnit
	Bitperhour           LeftAxisCustomUnit
	Bitperminute         LeftAxisCustomUnit
	Bitpersecond         LeftAxisCustomUnit
	Byte                 LeftAxisCustomUnit
	Byteperhour          LeftAxisCustomUnit
	Byteperminute        LeftAxisCustomUnit
	Bytepersecond        LeftAxisCustomUnit
	Cores                LeftAxisCustomUnit
	Count                LeftAxisCustomUnit
	Day                  LeftAxisCustomUnit
	Decibelmilliwatt     LeftAxisCustomUnit
	Gibibyte             LeftAxisCustomUnit
	Giga                 LeftAxisCustomUnit
	Gigabyte             LeftAxisCustomUnit
	Hour                 LeftAxisCustomUnit
	Kibibyte             LeftAxisCustomUnit
	Kibibyteperhour      LeftAxisCustomUnit
	Kibibyteperminute    LeftAxisCustomUnit
	Kibibytepersecond    LeftAxisCustomUnit
	Kilo                 LeftAxisCustomUnit
	Kilobyte             LeftAxisCustomUnit
	Kilobyteperhour      LeftAxisCustomUnit
	Kilobyteperminute    LeftAxisCustomUnit
	Kilobytepersecond    LeftAxisCustomUnit
	Mebibyte             LeftAxisCustomUnit
	Mebibyteperhour      LeftAxisCustomUnit
	Mebibyteperminute    LeftAxisCustomUnit
	Mebibytepersecond    LeftAxisCustomUnit
	Mega                 LeftAxisCustomUnit
	Megabyte             LeftAxisCustomUnit
	Megabyteperhour      LeftAxisCustomUnit
	Megabyteperminute    LeftAxisCustomUnit
	Megabytepersecond    LeftAxisCustomUnit
	Microsecond          LeftAxisCustomUnit
	Millicores           LeftAxisCustomUnit
	Millisecond          LeftAxisCustomUnit
	Millisecondperminute LeftAxisCustomUnit
	Minute               LeftAxisCustomUnit
	Month                LeftAxisCustomUnit
	Nanosecond           LeftAxisCustomUnit
	Nanosecondperminute  LeftAxisCustomUnit
	Notapplicable        LeftAxisCustomUnit
	Perhour              LeftAxisCustomUnit
	Perminute            LeftAxisCustomUnit
	Persecond            LeftAxisCustomUnit
	Percent              LeftAxisCustomUnit
	Pixel                LeftAxisCustomUnit
	Promille             LeftAxisCustomUnit
	Ratio                LeftAxisCustomUnit
	Second               LeftAxisCustomUnit
	State                LeftAxisCustomUnit
	Unspecified          LeftAxisCustomUnit
	Week                 LeftAxisCustomUnit
	Year                 LeftAxisCustomUnit
}{
	"Bit",
	"BitPerHour",
	"BitPerMinute",
	"BitPerSecond",
	"Byte",
	"BytePerHour",
	"BytePerMinute",
	"BytePerSecond",
	"Cores",
	"Count",
	"Day",
	"DecibelMilliWatt",
	"GibiByte",
	"Giga",
	"GigaByte",
	"Hour",
	"KibiByte",
	"KibiBytePerHour",
	"KibiBytePerMinute",
	"KibiBytePerSecond",
	"Kilo",
	"KiloByte",
	"KiloBytePerHour",
	"KiloBytePerMinute",
	"KiloBytePerSecond",
	"MebiByte",
	"MebiBytePerHour",
	"MebiBytePerMinute",
	"MebiBytePerSecond",
	"Mega",
	"MegaByte",
	"MegaBytePerHour",
	"MegaBytePerMinute",
	"MegaBytePerSecond",
	"MicroSecond",
	"MilliCores",
	"MilliSecond",
	"MilliSecondPerMinute",
	"Minute",
	"Month",
	"NanoSecond",
	"NanoSecondPerMinute",
	"NotApplicable",
	"PerHour",
	"PerMinute",
	"PerSecond",
	"Percent",
	"Pixel",
	"Promille",
	"Ratio",
	"Second",
	"State",
	"Unspecified",
	"Week",
	"Year",
}

// RightAxisCustomUnit has no documentation
type RightAxisCustomUnit string

func (me RightAxisCustomUnit) Ref() *RightAxisCustomUnit {
	return &me
}

// RightAxisCustomUnits has no documentation
var RightAxisCustomUnits = struct {
	Bit                  RightAxisCustomUnit
	Bitperhour           RightAxisCustomUnit
	Bitperminute         RightAxisCustomUnit
	Bitpersecond         RightAxisCustomUnit
	Byte                 RightAxisCustomUnit
	Byteperhour          RightAxisCustomUnit
	Byteperminute        RightAxisCustomUnit
	Bytepersecond        RightAxisCustomUnit
	Cores                RightAxisCustomUnit
	Count                RightAxisCustomUnit
	Day                  RightAxisCustomUnit
	Decibelmilliwatt     RightAxisCustomUnit
	Gibibyte             RightAxisCustomUnit
	Giga                 RightAxisCustomUnit
	Gigabyte             RightAxisCustomUnit
	Hour                 RightAxisCustomUnit
	Kibibyte             RightAxisCustomUnit
	Kibibyteperhour      RightAxisCustomUnit
	Kibibyteperminute    RightAxisCustomUnit
	Kibibytepersecond    RightAxisCustomUnit
	Kilo                 RightAxisCustomUnit
	Kilobyte             RightAxisCustomUnit
	Kilobyteperhour      RightAxisCustomUnit
	Kilobyteperminute    RightAxisCustomUnit
	Kilobytepersecond    RightAxisCustomUnit
	Mebibyte             RightAxisCustomUnit
	Mebibyteperhour      RightAxisCustomUnit
	Mebibyteperminute    RightAxisCustomUnit
	Mebibytepersecond    RightAxisCustomUnit
	Mega                 RightAxisCustomUnit
	Megabyte             RightAxisCustomUnit
	Megabyteperhour      RightAxisCustomUnit
	Megabyteperminute    RightAxisCustomUnit
	Megabytepersecond    RightAxisCustomUnit
	Microsecond          RightAxisCustomUnit
	Millicores           RightAxisCustomUnit
	Millisecond          RightAxisCustomUnit
	Millisecondperminute RightAxisCustomUnit
	Minute               RightAxisCustomUnit
	Month                RightAxisCustomUnit
	Nanosecond           RightAxisCustomUnit
	Nanosecondperminute  RightAxisCustomUnit
	Notapplicable        RightAxisCustomUnit
	Perhour              RightAxisCustomUnit
	Perminute            RightAxisCustomUnit
	Persecond            RightAxisCustomUnit
	Percent              RightAxisCustomUnit
	Pixel                RightAxisCustomUnit
	Promille             RightAxisCustomUnit
	Ratio                RightAxisCustomUnit
	Second               RightAxisCustomUnit
	State                RightAxisCustomUnit
	Unspecified          RightAxisCustomUnit
	Week                 RightAxisCustomUnit
	Year                 RightAxisCustomUnit
}{
	"Bit",
	"BitPerHour",
	"BitPerMinute",
	"BitPerSecond",
	"Byte",
	"BytePerHour",
	"BytePerMinute",
	"BytePerSecond",
	"Cores",
	"Count",
	"Day",
	"DecibelMilliWatt",
	"GibiByte",
	"Giga",
	"GigaByte",
	"Hour",
	"KibiByte",
	"KibiBytePerHour",
	"KibiBytePerMinute",
	"KibiBytePerSecond",
	"Kilo",
	"KiloByte",
	"KiloBytePerHour",
	"KiloBytePerMinute",
	"KiloBytePerSecond",
	"MebiByte",
	"MebiBytePerHour",
	"MebiBytePerMinute",
	"MebiBytePerSecond",
	"Mega",
	"MegaByte",
	"MegaBytePerHour",
	"MegaBytePerMinute",
	"MegaBytePerSecond",
	"MicroSecond",
	"MilliCores",
	"MilliSecond",
	"MilliSecondPerMinute",
	"Minute",
	"Month",
	"NanoSecond",
	"NanoSecondPerMinute",
	"NotApplicable",
	"PerHour",
	"PerMinute",
	"PerSecond",
	"Percent",
	"Pixel",
	"Promille",
	"Ratio",
	"Second",
	"State",
	"Unspecified",
	"Week",
	"Year",
}
