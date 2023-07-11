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

package workflows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ScheduleFilterParameters struct {
	Count             *int     `json:"count,omitempty" default:"10" minimum:"1"`                     // If specified, the schedule will end triggering executions af the given amount of executions. Minimum: 1, Maximum: 10
	EarliestStart     *string  `json:"earliestStart,omitempty" format:"date"`                        // If specified, the schedule won't trigger executions before the given date
	EarliestStartTime *string  `json:"earliestStartTime,omitempty" default:"00:00:00" format:"time"` // If specified, the schedule won't trigger executions before the given time
	Until             *string  `json:"until,omitempty" format:"date"`                                // If specified, the schedule won't trigger executions after the given date
	IncludeDates      []string `json:"includeDates,omitempty" format:"date"`                         // If specified, the schedule will trigger executions on the given dates, even if the main configuration prohibits it
	ExcludeDates      []string `json:"excludeDates,omitempty" format:"date"`                         // If specified, the schedule won't trigger exeuctions on the given dates
}

func (me *ScheduleFilterParameters) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"count": {
			Type:        schema.TypeInt,
			Description: "If specified, the schedule will end triggering executions af the given amount of executions. Minimum: 1, Maximum: 10",
			Optional:    true,
		},
		"earliest_start": {
			Type:        schema.TypeString,
			Description: "If specified, the schedule won't trigger executions before the given date",
			Optional:    true,
			// TODO: ValidateDiag 2023-07-01
		},
		"earliest_start_time": {
			Type:        schema.TypeString,
			Description: "If specified, the schedule won't trigger executions before the given time",
			Optional:    true,
			// TODO: ValidateDiag 22:32:22
		},
		"until": {
			Type:        schema.TypeString,
			Description: "If specified, the schedule won't trigger executions after the given date",
			Optional:    true,
			// TODO: ValidateDiag 2023-07-01
		},
		"include_dates": {
			Type:        schema.TypeSet,
			Description: "If specified, the schedule will trigger executions on the given dates, even if the main configuration prohibits it",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
			// TODO: ValidateDiag 2023-07-01
		},
		"exclude_dates": {
			Type:        schema.TypeSet,
			Description: "If specified, the schedule won't trigger exeuctions on the given dates",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
			// TODO: ValidateDiag 2023-07-01
		},
	}
}

func (me *ScheduleFilterParameters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"count":               me.Count,
		"earliest_start":      me.EarliestStart,
		"earliest_start_time": me.EarliestStartTime,
		"until":               me.Until,
		"include_dates":       me.IncludeDates,
		"exclude_dates":       me.ExcludeDates,
	})
}

func (me *ScheduleFilterParameters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"count":               &me.Count,
		"earliest_start":      &me.EarliestStart,
		"earliest_start_time": &me.EarliestStartTime,
		"until":               &me.Until,
		"include_dates":       &me.IncludeDates,
		"exclude_dates":       &me.ExcludeDates,
	})
}
