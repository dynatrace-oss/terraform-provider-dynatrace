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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WeeklyRecurrence struct {
	DayOfWeek       DayOfWeek        `json:"dayOfWeek"`       // The day of the week for weekly maintenance.  The format is the full name of the day in upper case, for example `THURSDAY`.
	TimeWindow      *TimeWindow      `json:"timeWindow"`      // The time window of the maintenance window
	RecurrenceRange *RecurrenceRange `json:"recurrenceRange"` // The recurrence date range of the maintenance window
}

func (me *WeeklyRecurrence) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"day_of_week": {
			Type:        schema.TypeString,
			Description: "The day of the week for weekly maintenance.  The format is the full name of the day in upper case, for example `THURSDAY`",
			Required:    true,
		},
		"time_window": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "The time window of the maintenance window",
			Elem: &schema.Resource{
				Schema: new(TimeWindow).Schema(),
			},
		},
		"recurrence_range": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "The recurrence date range of the maintenance window",
			Elem: &schema.Resource{
				Schema: new(RecurrenceRange).Schema(),
			},
		},
	}
}

func (me *WeeklyRecurrence) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"day_of_week":      me.DayOfWeek,
		"time_window":      me.TimeWindow,
		"recurrence_range": me.RecurrenceRange,
	})
}

func (me *WeeklyRecurrence) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"day_of_week":      &me.DayOfWeek,
		"time_window":      &me.TimeWindow,
		"recurrence_range": &me.RecurrenceRange,
	})
}

// DayOfWeek The day of the week for weekly maintenance.
// The format is the full name of the day in upper case, for example `THURSDAY`.
type DayOfWeek string

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
