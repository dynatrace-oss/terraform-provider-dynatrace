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

type MonthlyRecurrence struct {
	DayOfMonth      *int32           `json:"dayOfMonth"`      // The day of the month for monthly maintenance.  The value of `31` is treated as the last day of the month for months that don't have a 31st day. The value of `30` is also treated as the last day of the month for February.	TimeWindow      TimeWindow      `json:"timeWindow"`
	TimeWindow      *TimeWindow      `json:"timeWindow"`      // The time window of the maintenance window
	RecurrenceRange *RecurrenceRange `json:"recurrenceRange"` // The recurrence date range of the maintenance window
}

func (me *MonthlyRecurrence) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"day_of_month": {
			Type:        schema.TypeInt,
			Description: "The day of the month for monthly maintenance.  The value of `31` is treated as the last day of the month for months that don't have a 31st day. The value of `30` is also treated as the last day of the month for February",
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

func (me *MonthlyRecurrence) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"day_of_month":     me.DayOfMonth,
		"time_window":      me.TimeWindow,
		"recurrence_range": me.RecurrenceRange,
	})
}

func (me *MonthlyRecurrence) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"day_of_month":     &me.DayOfMonth,
		"time_window":      &me.TimeWindow,
		"recurrence_range": &me.RecurrenceRange,
	})
}
