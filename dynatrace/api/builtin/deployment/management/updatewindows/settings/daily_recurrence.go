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

type DailyRecurrence struct {
	Every           int              `json:"every"`           // Every **X** days:\n* `1` = every day,\n* `2` = every two days,\n* `3` = every three days,\n* etc.
	RecurrenceRange *RecurrenceRange `json:"recurrenceRange"` // Recurrence range
	UpdateTime      *UpdateTime      `json:"updateTime"`      // Update time
}

func (me *DailyRecurrence) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"every": {
			Type:        schema.TypeInt,
			Description: "Every **X** days:\n* `1` = every day,\n* `2` = every two days,\n* `3` = every three days,\n* etc.",
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

func (me *DailyRecurrence) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"every":            me.Every,
		"recurrence_range": me.RecurrenceRange,
		"update_time":      me.UpdateTime,
	})
}

func (me *DailyRecurrence) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"every":            &me.Every,
		"recurrence_range": &me.RecurrenceRange,
		"update_time":      &me.UpdateTime,
	})
}
