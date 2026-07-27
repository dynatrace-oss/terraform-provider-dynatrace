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

package updatewindows

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DailyRecurrence   *DailyRecurrence   `json:"dailyRecurrence,omitempty"`
	Enabled           bool               `json:"enabled"` // This setting is enabled (`true`) or disabled (`false`)
	MonthlyRecurrence *MonthlyRecurrence `json:"monthlyRecurrence,omitempty"`
	Name              string             `json:"name"` // Name
	OnceRecurrence    *OnceRecurrence    `json:"onceRecurrence,omitempty"`
	Recurrence        RecurrenceEnum     `json:"recurrence"` // Recurrence. Possible values: `DAILY`, `MONTHLY`, `ONCE`, `WEEKLY`
	WeeklyRecurrence  *WeeklyRecurrence  `json:"weeklyRecurrence,omitempty"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"daily_recurrence": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(DailyRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"monthly_recurrence": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(MonthlyRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"once_recurrence": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(OnceRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"recurrence": {
			Type:        schema.TypeString,
			Description: "Recurrence. Possible values: `DAILY`, `MONTHLY`, `ONCE`, `WEEKLY`",
			Required:    true,
		},
		"weekly_recurrence": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(WeeklyRecurrence).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"daily_recurrence":   me.DailyRecurrence,
		"enabled":            me.Enabled,
		"monthly_recurrence": me.MonthlyRecurrence,
		"name":               me.Name,
		"once_recurrence":    me.OnceRecurrence,
		"recurrence":         me.Recurrence,
		"weekly_recurrence":  me.WeeklyRecurrence,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.DailyRecurrence != nil) && (string(me.Recurrence) != "DAILY") {
		return fmt.Errorf("'daily_recurrence' must not be specified unless 'recurrence' is set to 'DAILY'; got 'recurrence'='%v'", me.Recurrence)
	}
	if (me.DailyRecurrence == nil) && (string(me.Recurrence) == "DAILY") {
		return fmt.Errorf("'daily_recurrence' must be specified when 'recurrence' is set to 'DAILY'; got 'recurrence'='%v'", me.Recurrence)
	}
	if (me.MonthlyRecurrence != nil) && (string(me.Recurrence) != "MONTHLY") {
		return fmt.Errorf("'monthly_recurrence' must not be specified unless 'recurrence' is set to 'MONTHLY'; got 'recurrence'='%v'", me.Recurrence)
	}
	if (me.MonthlyRecurrence == nil) && (string(me.Recurrence) == "MONTHLY") {
		return fmt.Errorf("'monthly_recurrence' must be specified when 'recurrence' is set to 'MONTHLY'; got 'recurrence'='%v'", me.Recurrence)
	}
	if (me.OnceRecurrence != nil) && (string(me.Recurrence) != "ONCE") {
		return fmt.Errorf("'once_recurrence' must not be specified unless 'recurrence' is set to 'ONCE'; got 'recurrence'='%v'", me.Recurrence)
	}
	if (me.OnceRecurrence == nil) && (string(me.Recurrence) == "ONCE") {
		return fmt.Errorf("'once_recurrence' must be specified when 'recurrence' is set to 'ONCE'; got 'recurrence'='%v'", me.Recurrence)
	}
	if (me.WeeklyRecurrence != nil) && (string(me.Recurrence) != "WEEKLY") {
		return fmt.Errorf("'weekly_recurrence' must not be specified unless 'recurrence' is set to 'WEEKLY'; got 'recurrence'='%v'", me.Recurrence)
	}
	if (me.WeeklyRecurrence == nil) && (string(me.Recurrence) == "WEEKLY") {
		return fmt.Errorf("'weekly_recurrence' must be specified when 'recurrence' is set to 'WEEKLY'; got 'recurrence'='%v'", me.Recurrence)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"daily_recurrence":   &me.DailyRecurrence,
		"enabled":            &me.Enabled,
		"monthly_recurrence": &me.MonthlyRecurrence,
		"name":               &me.Name,
		"once_recurrence":    &me.OnceRecurrence,
		"recurrence":         &me.Recurrence,
		"weekly_recurrence":  &me.WeeklyRecurrence,
	})
}
