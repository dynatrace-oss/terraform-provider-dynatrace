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

type UpdateTime struct {
	Duration  int          `json:"duration"`  // Duration (minutes)
	StartTime string       `json:"startTime"` // Start time (24-hour clock)
	TimeZone  TimezoneEnum `json:"timeZone"`  // Possible Values: `GMT_06_00`, `GMT_12_00`, `GMT_10_00`, `GMT_07_00`, `GMT_00_00`, `GMT_11_00`, `GMT_03_00`, `GMT_01_00`, `GMT_05_00`, `GMT_09_00`, `GMT_02_00`, `GMT_04_00`, `GMT_08_00`
}

func (me *UpdateTime) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"duration": {
			Type:        schema.TypeInt,
			Description: "Duration (minutes)",
			Required:    true,
		},
		"start_time": {
			Type:        schema.TypeString,
			Description: "Start time (24-hour clock)",
			Required:    true,
		},
		"time_zone": {
			Type:        schema.TypeString,
			Description: "Possible Values: `GMT_06_00`, `GMT_12_00`, `GMT_10_00`, `GMT_07_00`, `GMT_00_00`, `GMT_11_00`, `GMT_03_00`, `GMT_01_00`, `GMT_05_00`, `GMT_09_00`, `GMT_02_00`, `GMT_04_00`, `GMT_08_00`",
			Required:    true,
		},
	}
}

func (me *UpdateTime) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"duration":   me.Duration,
		"start_time": me.StartTime,
		"time_zone":  me.TimeZone,
	})
}

func (me *UpdateTime) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"duration":   &me.Duration,
		"start_time": &me.StartTime,
		"time_zone":  &me.TimeZone,
	})
}
