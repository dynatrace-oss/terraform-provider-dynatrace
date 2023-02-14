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

type MonthlyRecurrence struct {
	Every            int              `json:"every"`            // Every **X** months:\n* `1` = every month,\n* `2` = every two months,\n* `3` = every three months,\n* etc.
	RecurrenceRange  *RecurrenceRange `json:"recurrenceRange"`  // Recurrence range
	SelectedMonthDay int              `json:"selectedMonthDay"` // Day of the month
	UpdateTime       *UpdateTime      `json:"updateTime"`       // Update time
}

func (me *MonthlyRecurrence) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"every": {
			Type:        schema.TypeInt,
			Description: "Every **X** months:\n* `1` = every month,\n* `2` = every two months,\n* `3` = every three months,\n* etc.",
			Required:    true,
		},
		"recurrence_range": {
			Type:        schema.TypeList,
			Description: "Recurrence range",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(RecurrenceRange).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"selected_month_day": {
			Type:        schema.TypeInt,
			Description: "Day of the month",
			Required:    true,
		},
		"update_time": {
			Type:        schema.TypeList,
			Description: "Update time",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(UpdateTime).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *MonthlyRecurrence) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"every":              me.Every,
		"recurrence_range":   me.RecurrenceRange,
		"selected_month_day": me.SelectedMonthDay,
		"update_time":        me.UpdateTime,
	})
}

func (me *MonthlyRecurrence) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"every":              &me.Every,
		"recurrence_range":   &me.RecurrenceRange,
		"selected_month_day": &me.SelectedMonthDay,
		"update_time":        &me.UpdateTime,
	})
}
