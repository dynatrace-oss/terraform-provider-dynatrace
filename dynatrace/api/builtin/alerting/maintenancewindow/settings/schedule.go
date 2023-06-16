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

package maintenancewindow

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Schedule struct {
	DailyRecurrence   *DailyRecurrence   `json:"dailyRecurrence,omitempty"`   // The configuration for maintenance windows occuring daily
	MonthlyRecurrence *MonthlyRecurrence `json:"monthlyRecurrence,omitempty"` // The configuration for maintenance windows occuring monthly
	OnceRecurrence    *OnceRecurrence    `json:"onceRecurrence,omitempty"`    // The configuration for maintenance windows occuring daily once
	ScheduleType      ScheduleType       `json:"scheduleType"`                // The type maintenance window, possible values: `DAILY`, `MONTHLY`, `ONCE`, `WEEKLY`
	WeeklyRecurrence  *WeeklyRecurrence  `json:"weeklyRecurrence,omitempty"`  // The configuration for maintenance windows occuring weekly
}

func (me *Schedule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"daily_recurrence": {
			Type:        schema.TypeList,
			Description: "The configuration for maintenance windows occuring daily",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(DailyRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"monthly_recurrence": {
			Type:        schema.TypeList,
			Description: "The configuration for maintenance windows occuring monthly",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(MonthlyRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"once_recurrence": {
			Type:        schema.TypeList,
			Description: "The configuration for maintenance windows occuring once",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(OnceRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type maintenance window, possible values: `DAILY`, `MONTHLY`, `ONCE`, `WEEKLY`",
			Required:    true,
		},
		"weekly_recurrence": {
			Type:        schema.TypeList,
			Description: "The configuration for maintenance windows occuring weekly",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(WeeklyRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Schedule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"daily_recurrence":   me.DailyRecurrence,
		"monthly_recurrence": me.MonthlyRecurrence,
		"once_recurrence":    me.OnceRecurrence,
		"type":               me.ScheduleType,
		"weekly_recurrence":  me.WeeklyRecurrence,
	})
}

func (me *Schedule) HandlePreconditions() error {
	if me.DailyRecurrence == nil && (string(me.ScheduleType) == "DAILY") {
		return fmt.Errorf("'daily_recurrence' must be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	if me.DailyRecurrence != nil && (string(me.ScheduleType) != "DAILY") {
		return fmt.Errorf("'daily_recurrence' must not be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	if me.MonthlyRecurrence == nil && (string(me.ScheduleType) == "MONTHLY") {
		return fmt.Errorf("'monthly_recurrence' must be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	if me.MonthlyRecurrence != nil && (string(me.ScheduleType) != "MONTHLY") {
		return fmt.Errorf("'monthly_recurrence' must not be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	if me.OnceRecurrence == nil && (string(me.ScheduleType) == "ONCE") {
		return fmt.Errorf("'once_recurrence' must be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	if me.OnceRecurrence != nil && (string(me.ScheduleType) != "ONCE") {
		return fmt.Errorf("'once_recurrence' must not be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	if me.WeeklyRecurrence == nil && (string(me.ScheduleType) == "WEEKLY") {
		return fmt.Errorf("'weekly_recurrence' must be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	if me.WeeklyRecurrence != nil && (string(me.ScheduleType) != "WEEKLY") {
		return fmt.Errorf("'weekly_recurrence' must not be specified if 'type' is set to '%v'", me.ScheduleType)
	}
	return nil
}

func (me *Schedule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"daily_recurrence":   &me.DailyRecurrence,
		"monthly_recurrence": &me.MonthlyRecurrence,
		"once_recurrence":    &me.OnceRecurrence,
		"type":               &me.ScheduleType,
		"weekly_recurrence":  &me.WeeklyRecurrence,
	})
}
