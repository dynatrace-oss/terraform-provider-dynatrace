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

type RecurrenceRange struct {
	StartDate string `json:"scheduleStartDate"` // The start date of the recurrence range in YYYY-MM-DD format
	EndDate   string `json:"scheduleEndDate"`   // The end date of the recurrence range in YYYY-MM-DD format
}

func (me *RecurrenceRange) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_date": {
			Type:        schema.TypeString,
			Description: "The start date of the recurrence range in YYYY-MM-DD format",
			Required:    true,
		},
		"end_date": {
			Type:        schema.TypeString,
			Description: "The end date of the recurrence range in YYYY-MM-DD format",
			Required:    true,
		},
	}
}

func (me *RecurrenceRange) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"start_date": me.StartDate,
		"end_date":   me.EndDate,
	})
}

func (me *RecurrenceRange) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"start_date": &me.StartDate,
		"end_date":   &me.EndDate,
	})
}
