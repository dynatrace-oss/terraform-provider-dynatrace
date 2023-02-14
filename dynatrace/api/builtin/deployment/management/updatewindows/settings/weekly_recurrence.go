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

type WeeklyRecurrence struct {
	Every            int               `json:"every"`            // Every **X** weeks:\n* `1` = every week,\n* `2` = every two weeks,\n* `3` = every three weeks,\n* etc.
	RecurrenceRange  *RecurrenceRange  `json:"recurrenceRange"`  // Recurrence range
	SelectedWeekDays *SelectedWeekDays `json:"selectedWeekDays"` // Day of the week
	UpdateTime       *UpdateTime       `json:"updateTime"`       // Update time
}

func (me *WeeklyRecurrence) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"every": {
			Type:        schema.TypeInt,
			Description: "Every **X** weeks:\n* `1` = every week,\n* `2` = every two weeks,\n* `3` = every three weeks,\n* etc.",
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
		"selected_week_days": {
			Type:        schema.TypeList,
			Description: "Day of the week",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SelectedWeekDays).Schema()},
			MinItems: 1,
			MaxItems: 1,
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

func (me *WeeklyRecurrence) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"every":              me.Every,
		"recurrence_range":   me.RecurrenceRange,
		"selected_week_days": me.SelectedWeekDays,
		"update_time":        me.UpdateTime,
	})
}

func (me *WeeklyRecurrence) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"every":              &me.Every,
		"recurrence_range":   &me.RecurrenceRange,
		"selected_week_days": &me.SelectedWeekDays,
		"update_time":        &me.UpdateTime,
	})
}
