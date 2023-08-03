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

package scheduling_rules

import (
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type RecurrenceRule struct {
	Frequency  Frequency `json:"freq"`                                       // Possible values are `YEARLY`, `MONTHLY`, `WEEKLY`, `DAILY`, `HOURLY`, `MINUTELY` and `SECONDLY`
	DateStart  string    `json:"datestart" format:"date"`                    // The recurrence start
	Interval   int       `json:"interval,omitempty" default:"1" minimum:"0"` // The interval between each freq iteration
	ByMonth    []int     `json:"bymonth,omitempty"`                          // The months to apply the recurrence to. `1` for `January`, `2` for `February`, ..., `12` for `December`
	ByMonthDay []int     `json:"bymonthday,omitempty"`                       // The days within a month to apply the recurrence to. `1`, `2`, `3`, ... refers to the first, second, third day in the month. You can also specify negative values to refer to values relative to the last day. `-1` refers to the last day, `-2` refers to the second to the last day, ...
	ByYearDay  []int     `json:"byyearday,omitempty"`                        // The days within a year to apply the recurrence to. `1`, `2`, `3`, ... refers to the first, second, third day in the year. You can also specify negative values to refer to values relative to the last day. `-1` refers to the last day, `-2` refers to the second to the last day, ...
	ByEaster   []int     `json:"byeaster,omitempty"`                         // Each value will define an offset from the Easter Sunday. The offset `0` will yield the Easter Sunday itself
	ByWeekNo   []int     `json:"byweekno,omitempty"`                         // The calendar week within the year to apply the recurrence to. `1`, `2`, `3`, ... refers to the first, second, third week of the year. You can also specify negative values to refer to values relative to the last week. `-1` refers to the last week, `-2` refers to the second to the last week, ...
	ByWeekDay  []string  `json:"byday,omitempty"`                            // Define the weekdays where the recurrence will be applied. Possible values are `MO`, `TU`, `WE`, `TH`, `FR`, `SA` and `SU`
	ByWorkDay  string    `json:"automation_server_byworkday,omitempty"`      // Possible values are `WORKING` (Work days), `HOLIDAYS` (Holidays) and `OFF` (Weekends + Holidays)

	// WeekStartDay string    `json:"wkst,omitempty"`                             // The week start day. Possible values are `MO`, `TU`, `WE`, `TH`, `FR`, `SA` and `SU`
	// Count        int       `json:"count,omitempty"`                            // How many occurrences will be generated
	// Until        string    `json:"until,omitempty" format:"date-time"`         // Specifying the upper-bound limit of the recurrence
	// BySetPos     []int     `json:"bysetpos,omitempty"`                         // Each given integer will specify an occurrence number, corresponding to the nth occurrence of the rule inside the frequency period
	// ByHour     []int    `json:"byhour,omitempty"`     // The hours to apply the recurrence to
	// ByMinute   []int    `json:"byminute,omitempty"`   // The minutes to apply the recurrence to
	// BySecond   []int    `json:"bysecond,omitempty"`   // The seconds to apply the recurrence to
}

type jsonrule RecurrenceRule

func (me *RecurrenceRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"frequency": {
			Type:         schema.TypeString,
			Description:  "Possible values are `YEARLY`, `MONTHLY`, `WEEKLY`, `DAILY`, `HOURLY`, `MINUTELY` and `SECONDLY`. Example: `frequency` = `DAILY` and `interval` = `2` schedules for every other day",
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"YEARLY", "MONTHLY", "WEEKLY", "DAILY", "HOURLY", "MINUTELY", "SECONDLY"}, false),
		},
		"datestart": {
			Type:        schema.TypeString,
			Description: "The recurrence start. Example: `2017-07-04` represents July 4th 2017",
			Required:    true,
		},
		"interval": {
			Type:        schema.TypeInt,
			Description: "The interval between each iteration. Default: 1. Example: `frequency` = `DAILY` and `interval` = `2` schedules for every other day",
			Default:     1,
			Optional:    true,
		},
		"months": {
			Type:        schema.TypeSet,
			Description: "Restricts the recurrence to specific months. `1` for `January`, `2` for `February`, ..., `12` for `December`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"days_in_month": {
			Type:        schema.TypeSet,
			Description: "Restricts the recurrence to specific days within a month. `1`, `2`, `3`, ... refers to the first, second, third day in the month. You can also specify negative values to refer to values relative to the last day. `-1` refers to the last day, `-2` refers to the second to the last day, ...",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"days_in_year": {
			Type:        schema.TypeSet,
			Description: "Restricts the recurrence to specific days within a year. `1`, `2`, `3`, ... refers to the first, second, third day of the year. You can also specify negative values to refer to values relative to the last day. `-1` refers to the last day, `-2` refers to the second to the last day, ...",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"easter": {
			Type:        schema.TypeSet,
			Description: "Restricts the recurrence to specific days relative to Easter Sunday. `0` will yield the Easter Sunday itself",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"weeks": {
			Type:        schema.TypeSet,
			Description: "Restricts the recurrence to specific weeks within a year. `1`, `2`, `3`, ... refers to the first, second, third week of the year. You can also specify negative values to refer to values relative to the last week. `-1` refers to the last week, `-2` refers to the second to the last week, ...",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"weekdays": {
			Type:        schema.TypeSet,
			Description: "Restricts the recurrence to specific week days. Possible values are `MO`, `TU`, `WE`, `TH`, `FR`, `SA` and `SU`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"workdays": {
			Type:         schema.TypeString,
			Description:  "Possible values are `WORKING` (Work days), `HOLIDAYS` (Holidays) and `OFF` (Weekends + Holidays)",
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"WORKING", "HOLIDAYS", "OFF"}, false),
		},
	}
}

func (me *RecurrenceRule) MarshalHCL(properties hcl.Properties) error {
	if len(me.ByMonth) > 0 {
		sort.Ints(me.ByMonth)
	}
	if len(me.ByMonthDay) > 0 {
		sort.Ints(me.ByMonthDay)
	}
	if len(me.ByYearDay) > 0 {
		sort.Ints(me.ByYearDay)
	}
	if len(me.ByEaster) > 0 {
		sort.Ints(me.ByEaster)
	}
	if len(me.ByWeekNo) > 0 {
		sort.Ints(me.ByWeekNo)
	}
	if len(me.ByWeekNo) > 0 {
		sort.Ints(me.ByWeekNo)
	}
	return properties.EncodeAll(map[string]any{
		"frequency":     me.Frequency,
		"datestart":     me.DateStart,
		"interval":      me.Interval,
		"months":        me.ByMonth,
		"days_in_month": me.ByMonthDay,
		"days_in_year":  me.ByYearDay,
		"easter":        me.ByEaster,
		"weeks":         me.ByWeekNo,
		"weekdays":      me.ByWeekDay,
		"workdays":      me.ByWorkDay,
	})
}

func (me *RecurrenceRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"frequency":     &me.Frequency,
		"datestart":     &me.DateStart,
		"interval":      &me.Interval,
		"months":        &me.ByMonth,
		"days_in_month": &me.ByMonthDay,
		"days_in_year":  &me.ByYearDay,
		"easter":        &me.ByEaster,
		"weeks":         &me.ByWeekNo,
		"weekdays":      &me.ByWeekDay,
		"workdays":      &me.ByWorkDay,
	})
}

// func (me *RecurrenceRule) MarshalJSON() ([]byte, error) {
// 	jr := jsonrule(*me)
// 	data, err := json.Marshal(jr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m := map[string]json.RawMessage{}
// 	if err := json.Unmarshal(data, &m); err != nil {
// 		return nil, err
// 	}
// 	if v, found := m["automation_server_byworkday"]; found {
// 		m["byworkday"] = v
// 	} else if v, found := m["byworkday"]; found {
// 		m["automation_server_byworkday"] = v
// 	}
// 	return json.Marshal(m)
// }

func (me *RecurrenceRule) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["automation_server_byworkday"]; found {
		m["byworkday"] = v
	} else if v, found := m["byworkday"]; found {
		m["automation_server_byworkday"] = v
	}
	var jr jsonrule
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &jr); err != nil {
		return err
	}
	*me = RecurrenceRule(jr)
	return nil
}
