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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type UserFriendlySettings struct {
	DailyRecurrence   *DailyRecurrence   `json:"dailyRecurrence,omitempty"`
	Enabled           bool               `json:"enabled"` // This setting is enabled (`true`) or disabled (`false`)
	MonthlyRecurrence *MonthlyRecurrence `json:"monthlyRecurrence,omitempty"`
	Name              string             `json:"name"` // Name
	OnceRecurrence    *OnceRecurrence    `json:"onceRecurrence,omitempty"`
	Recurrence        RecurrenceEnum     `json:"recurrence"` // Possible Values: `MONTHLY`, `ONCE`, `DAILY`, `WEEKLY`
	WeeklyRecurrence  *WeeklyRecurrence  `json:"weeklyRecurrence,omitempty"`
}

func (me *UserFriendlySettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"daily": {
			Type:         schema.TypeList,
			Description:  "no documentation available",
			Optional:     true,
			Elem:         &schema.Resource{Schema: new(DailyRecurrence).Schema()},
			MinItems:     1,
			MaxItems:     1,
			ExactlyOneOf: []string{"daily", "monthly", "weekly", "once"},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"monthly": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(MonthlyRecurrence).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"once": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(OnceRecurrence).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		// "recurrence": {
		// 	Type:        schema.TypeString,
		// 	Description: "Possible Values: `MONTHLY`, `ONCE`, `DAILY`, `WEEKLY`",
		// 	Required:    true,
		// },
		"weekly": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(WeeklyRecurrence).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *UserFriendlySettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"daily":   me.DailyRecurrence,
		"enabled": me.Enabled,
		"monthly": me.MonthlyRecurrence,
		"name":    me.Name,
		"once":    me.OnceRecurrence,
		// "recurrence": me.Recurrence,
		"weekly": me.WeeklyRecurrence,
	})
}

func (me *UserFriendlySettings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"daily":   &me.DailyRecurrence,
		"enabled": &me.Enabled,
		"monthly": &me.MonthlyRecurrence,
		"name":    &me.Name,
		"once":    &me.OnceRecurrence,
		// "recurrence": &me.Recurrence,
		"weekly": &me.WeeklyRecurrence,
	}); err != nil {
		return err
	}
	if me.DailyRecurrence != nil {
		me.Recurrence = RecurrenceEnums.Daily
	} else if me.MonthlyRecurrence != nil {
		me.Recurrence = RecurrenceEnums.Monthly
	} else if me.OnceRecurrence != nil {
		me.Recurrence = RecurrenceEnums.Once
	} else if me.WeeklyRecurrence != nil {
		me.Recurrence = RecurrenceEnums.Weekly
	}
	return nil
}
