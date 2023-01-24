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

type Schedule struct {
	Type              ScheduleType       `json:"scheduleType"`
	OnceRecurrence    *OnceRecurrence    `json:"onceRecurrence,omitempty"`
	DailyRecurrence   *DailyRecurrence   `json:"dailyRecurrence,omitempty"`
	WeeklyRecurrence  *WeeklyRecurrence  `json:"weeklyRecurrence,omitempty"`
	MonthlyRecurrence *MonthlyRecurrence `json:"monthlyRecurrence,omitempty"`
}

func (me *Schedule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The time window of the maintenance window",
		},
		"once_recurrence": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration for maintenance windows occuring once",
			Elem: &schema.Resource{
				Schema: new(OnceRecurrence).Schema(),
			},
		},
		"daily_recurrence": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration for maintenance windows occuring daily",
			Elem: &schema.Resource{
				Schema: new(DailyRecurrence).Schema(),
			},
		},
		"weekly_recurrence": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration for maintenance windows occuring weekly",
			Elem: &schema.Resource{
				Schema: new(WeeklyRecurrence).Schema(),
			},
		},
		"monthly_recurrence": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration for maintenance windows occuring monthly",
			Elem: &schema.Resource{
				Schema: new(MonthlyRecurrence).Schema(),
			},
		},
	}
}

func (me *Schedule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":               me.Type,
		"once_recurrence":    me.OnceRecurrence,
		"daily_recurrence":   me.DailyRecurrence,
		"weekly_recurrence":  me.WeeklyRecurrence,
		"monthly_recurrence": me.MonthlyRecurrence,
	})
}

func (me *Schedule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":               &me.Type,
		"once_recurrence":    &me.OnceRecurrence,
		"daily_recurrence":   &me.DailyRecurrence,
		"weekly_recurrence":  &me.WeeklyRecurrence,
		"monthly_recurrence": &me.MonthlyRecurrence,
	})
}

// ScheduleType The type of the maintenance: planned or unplanned.
type ScheduleType string

// ScheduleTypes offers the known enum values
var ScheduleTypes = struct {
	Once    ScheduleType
	Daily   ScheduleType
	Weekly  ScheduleType
	Monthly ScheduleType
}{
	"ONCE",
	"DAILY",
	"WEEKLY",
	"MONTHLY",
}
