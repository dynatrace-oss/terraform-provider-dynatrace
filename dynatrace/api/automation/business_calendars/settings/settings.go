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

package business_calendars

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Title       string   `json:"title" maxlength:"200" minlength:"1"`         // The title / name of the Business Calendar
	Description string   `json:"description,omitempty"`                       // An optional description for the Business Calendar
	WeekStart   *int     `json:"weekstart,omitempty" minimum:"1" maximum:"7"` // Specifies the day of the week that's considered to be the first day in the week. `1` for Monday, `7` for Sunday
	WeekDays    []int    `json:"weekdays"`                                    // The days to be considered week days in this calendar. `1' = `Monday`, `2` = `Tuesday`, `3` = `Wednesday`, `4` = `Thursday`, `5` = `Friday`, `6` = `Saturday`, `7` = `Sunday`
	Holidays    Holidays `json:"holidays,omitempty"`                          // A list of holidays valid in this calendar
	ValidFrom   *string  `json:"validFrom,omitempty" format:"date"`           // The date from when on this calendar is valid from
	ValidTo     *string  `json:"validTo,omitempty" format:"date"`             // The date until when on this calendar is valid to
	// Labels      map[string]string `json:"labels,omitempty" pattern:"^[A-Za-z0-9\\.\\s_-]*$"`
}

func (me *Settings) Name() string {
	return me.Title
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:             schema.TypeString,
			Description:      "The title / name of the Business Calendar",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(200),
		},
		"description": {
			Type:        schema.TypeString,
			Description: "An optional description for the Business Calendar",
			Optional:    true,
		},
		"week_start": {
			Type:        schema.TypeInt,
			Description: "Specifies the day of the week that's considered to be the first day in the week. `1` for Monday, `7` for Sunday",
			Optional:    true,
			Default:     1,
		},
		"week_days": {
			Type:        schema.TypeSet,
			Description: "The days to be considered week days in this calendar. `1' = `Monday`, `2` = `Tuesday`, `3` = `Wednesday`, `4` = `Thursday`, `5` = `Friday`, `6` = `Saturday`, `7` = `Sunday`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"holidays": {
			Type:        schema.TypeList,
			Description: "A list of holidays valid in this calendar",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Holidays).Schema()},
		},
		"valid_from": {
			Type:        schema.TypeString,
			Description: "The date from when on this calendar is valid from. Example: `2023-07-04` for July 4th 2023",
			Optional:    true,
		},
		"valid_to": {
			Type:        schema.TypeString,
			Description: "The date until when on this calendar is valid to. Example: `2023-07-04` for July 4th 2023",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"title":       me.Title,
		"description": me.Description,

		"week_start": me.WeekStart,
		"week_days":  me.WeekDays,
		"holidays":   me.Holidays,
		"valid_from": me.ValidFrom,
		"valid_to":   me.ValidTo,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"title":       &me.Title,
		"description": &me.Description,

		"week_start": &me.WeekStart,
		"week_days":  &me.WeekDays,
		"holidays":   &me.Holidays,
		"valid_from": &me.ValidFrom,
		"valid_to":   &me.ValidTo,
	})
}
