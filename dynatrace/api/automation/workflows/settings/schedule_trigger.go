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
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ScheduleTrigger struct {
	Type ScheduleTriggerType `json:"type"`

	// TimeTrigger - type <= time
	Time string `json:"time" pattern:"^([0-1]{1}[0-9]|2[0-3]):[0-5][0-9]$"` // Specifies a fixed time the schedule will trigger at in 24h format (e.g. `14:23:59`). Conflicts with `cron`, `interval_minutes`, `between_start` and `between_end`

	// CronTrigger - type <= cron
	Cron string `json:"cron" example:"0 0 * * *"` // Configures using cron syntax. Conflicts with `time`, `interval_minutes`, `between_start` and `between_end`.

	// IntervalTrigger - type <= interval
	IntervalMinutes int    `json:"intervalMinutes" minimum:"1" maximum:"720"`                            // Triggers the schedule every n minutes within a given time frame. Conflicts with `cron` and `time`. Required with `between_start` and `between_end`.
	BetweenStart    string `json:"betweenStart,omitempty" pattern:"^([0-1]{1}[0-9]|2[0-3]):[0-5][0-9]$"` // Triggers the schedule every n minutes within a given time frame. Conflicts with `cron` and `time`. Required with `interval_minutes` and `between_end`.
	BetweenEnd      string `json:"betweenEnd,omitempty" pattern:"^([0-1]{1}[0-9]|2[0-3]):[0-5][0-9]$"`   // Triggers the schedule every n minutes within a given time frame. Conflicts with `cron` and `time`. Required with `between_start` and `interval_minutes`.
}

func (me *ScheduleTrigger) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"time": {
			Type:             schema.TypeString,
			Description:      "Specifies a fixed time the schedule will trigger at in 24h format (e.g. `14:23:59`). Conflicts with `cron`, `interval_minutes`, `between_start` and `between_end`",
			Optional:         true,
			ValidateDiagFunc: ValidateRegex(regexp.MustCompile("^([0-1]{1}[0-9]|2[0-3]):[0-5][0-9]$"), "Example: 14:23:59"),
			ConflictsWith:    []string{prefix + ".0.cron", prefix + ".0.interval_minutes", prefix + ".0.between_start", prefix + ".0.between_end"},
		},
		"cron": {
			Type:          schema.TypeString,
			Description:   "Configures using cron syntax. Conflicts with `time`, `interval_minutes`, `between_start` and `between_end`",
			Optional:      true,
			ConflictsWith: []string{prefix + ".0.time", prefix + ".0.interval_minutes", prefix + ".0.between_start", prefix + ".0.between_end"},
		},
		"interval_minutes": {
			Type:          schema.TypeString,
			Description:   "Triggers the schedule every n minutes within a given time frame. Minimum: 1, Maximum: 720. Required with `between_start` and `between_end`. Conflicts with `cron` and `time`",
			Optional:      true,
			ConflictsWith: []string{prefix + ".0.time", prefix + ".0.cron"},
			RequiredWith:  []string{prefix + ".0.between_start", prefix + ".0.between_end"},
		},
		"between_start": {
			Type:          schema.TypeString,
			Description:   "Triggers the schedule every n minutes within a given time frame - specifying the start time on any valid day in 24h format (e.g. 13:22:44). Conflicts with `cron` and `time`. Required with `interval_minutes` and `between_end`",
			Optional:      true,
			ConflictsWith: []string{prefix + ".0.time", prefix + ".0.cron"},
			RequiredWith:  []string{prefix + ".0.interval_minutes", prefix + ".0.between_end"},
		},
		"between_end": {
			Type:          schema.TypeString,
			Description:   "Triggers the schedule every n minutes within a given time frame - specifying the end time on any valid day in 24h format (e.g. 14:22:44). Conflicts with `cron` and `time`. Required with `interval_minutes` and `between_start`",
			Optional:      true,
			ConflictsWith: []string{prefix + ".0.time", prefix + ".0.cron"},
			RequiredWith:  []string{prefix + ".0.interval_minutes", prefix + ".0.between_start"},
		},
	}
}

func (me *ScheduleTrigger) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"time":             me.Time,
		"cron":             me.Cron,
		"interval_minutes": me.IntervalMinutes,
		"between_start":    me.BetweenStart,
		"between_end":      me.BetweenEnd,
	})
}

func (me *ScheduleTrigger) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"time":             &me.Time,
		"cron":             &me.Cron,
		"interval_minutes": &me.IntervalMinutes,
		"between_start":    &me.BetweenStart,
		"between_end":      &me.BetweenEnd,
	})
}
