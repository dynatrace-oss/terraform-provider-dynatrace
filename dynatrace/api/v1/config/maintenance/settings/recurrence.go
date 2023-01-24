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

package maintenance

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Recurrence The recurrence of the maintenance window.
type Recurrence struct {
	DayOfMonth      *int32                     `json:"dayOfMonth,omitempty"` // The day of the month for monthly maintenance.  The value of `31` is treated as the last day of the month for months that don't have a 31st day. The value of `30` is also treated as the last day of the month for February.
	DayOfWeek       *DayOfWeek                 `json:"dayOfWeek,omitempty"`  // The day of the week for weekly maintenance.  The format is the full name of the day in upper case, for example `THURSDAY`.
	DurationMinutes int32                      `json:"durationMinutes"`      // The duration of the maintenance window in minutes.
	StartTime       string                     `json:"startTime"`            // The start time of the maintenance window in HH:mm format.
	Unknowns        map[string]json.RawMessage `json:"-"`
}

func (me *Recurrence) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"day_of_month": {
			Type:        schema.TypeInt,
			Description: "The day of the month for monthly maintenance.  The value of `31` is treated as the last day of the month for months that don't have a 31st day. The value of `30` is also treated as the last day of the month for February",
			Optional:    true,
		},
		"day_of_week": {
			Type:        schema.TypeString,
			Description: "The day of the week for weekly maintenance.  The format is the full name of the day in upper case, for example `THURSDAY`",
			Optional:    true,
		},
		"duration_minutes": {
			Type:        schema.TypeInt,
			Description: "The duration of the maintenance window in minutes",
			Required:    true,
		},
		"start_time": {
			Type:        schema.TypeString,
			Description: "The start time of the maintenance window in HH:mm format",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Recurrence) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("day_of_month", int(opt.Int32(me.DayOfMonth))); err != nil {
		return err
	}
	if err := properties.Encode("day_of_week", me.DayOfWeek); err != nil {
		return err
	}
	if err := properties.Encode("duration_minutes", int(me.DurationMinutes)); err != nil {
		return err
	}
	if err := properties.Encode("start_time", me.StartTime); err != nil {
		return err
	}
	return nil
}

func (me *Recurrence) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "day_of_month")
		delete(me.Unknowns, "day_of_week")
		delete(me.Unknowns, "duration_minutes")
		delete(me.Unknowns, "start_time")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("day_of_month"); ok {
		me.DayOfMonth = opt.NewInt32(int32(value.(int)))
	}
	if value, ok := decoder.GetOk("day_of_week"); ok {
		me.DayOfWeek = DayOfWeek(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("duration_minutes"); ok {
		me.DurationMinutes = int32(value.(int))
	}
	if value, ok := decoder.GetOk("start_time"); ok {
		me.StartTime = value.(string)
	}
	return nil
}

func (me *Recurrence) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("dayOfMonth", me.DayOfMonth); err != nil {
		return nil, err
	}
	if err := m.Marshal("dayOfWeek", me.DayOfWeek); err != nil {
		return nil, err
	}
	if err := m.Marshal("durationMinutes", me.DurationMinutes); err != nil {
		return nil, err
	}
	if err := m.Marshal("startTime", me.StartTime); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *Recurrence) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("dayOfMonth", &me.DayOfMonth); err != nil {
		return err
	}
	if err := m.Unmarshal("dayOfWeek", &me.DayOfWeek); err != nil {
		return err
	}
	if err := m.Unmarshal("durationMinutes", &me.DurationMinutes); err != nil {
		return err
	}
	if err := m.Unmarshal("startTime", &me.StartTime); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// DayOfWeek The day of the week for weekly maintenance.
// The format is the full name of the day in upper case, for example `THURSDAY`.
type DayOfWeek string

func (me DayOfWeek) Ref() *DayOfWeek {
	return &me
}

// DayOfWeeks offers the known enum values
var DayOfWeeks = struct {
	Friday    DayOfWeek
	Monday    DayOfWeek
	Saturday  DayOfWeek
	Sunday    DayOfWeek
	Thursday  DayOfWeek
	Tuesday   DayOfWeek
	Wednesday DayOfWeek
}{
	"FRIDAY",
	"MONDAY",
	"SATURDAY",
	"SUNDAY",
	"THURSDAY",
	"TUESDAY",
	"WEDNESDAY",
}
