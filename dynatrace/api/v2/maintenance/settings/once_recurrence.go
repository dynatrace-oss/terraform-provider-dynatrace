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

package maintenance

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type OnceRecurrence struct {
	StartTime string `json:"startTime"` // The start time of the maintenance window validity period in YYYY-MM-DDThh:mm:ss format (for example, `2022-01-01T08:00:00`)
	EndTime   string `json:"endTime"`   // The end time of the maintenance window validity period in YYYY-MM-DDThh:mm:ss format (for example, `2022-01-01T08:00:00`)
	TimeZone  string `json:"timeZone"`  // The time zone of the start and end time. Default time zone is UTC. You can use either UTC offset `UTC+01:00` format or the IANA Time Zone Database format (for example, `Europe/Vienna`)
}

func (me *OnceRecurrence) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_time": {
			Type:        schema.TypeString,
			Description: "The start time of the maintenance window validity period in YYYY-MM-DDThh:mm:ss format (for example, `2022-01-01T08:00:00`)",
			Required:    true,
		},
		"end_time": {
			Type:        schema.TypeString,
			Description: "The end time of the maintenance window validity period in YYYY-MM-DDThh:mm:ss format (for example, `2022-01-01T08:00:00`)",
			Required:    true,
		},
		"time_zone": {
			Type:        schema.TypeString,
			Description: "The time zone of the start and end time. Default time zone is UTC. You can use either UTC offset `UTC+01:00` format or the IANA Time Zone Database format (for example, `Europe/Vienna`)",
			Required:    true,
		},
	}
}

func (me *OnceRecurrence) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"start_time": me.StartTime,
		"end_time":   me.EndTime,
		"time_zone":  me.TimeZone,
	})
}

func (me *OnceRecurrence) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"start_time": &me.StartTime,
		"end_time":   &me.EndTime,
		"time_zone":  &me.TimeZone,
	})
}
