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

package reports

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Subscriptions struct {
	Month []string `json:"MONTH,omitempty"` // A list of monthly subscribers.\n Monthly subscribers receive the report on the first Monday of the month at midnight.\n You can specify email addresses or Dynatrace user IDs here.
	Week  []string `json:"WEEK,omitempty"`  // A list of weekly subscribers.\n Weekly subscribers receive the report every Monday at midnight.\n You can specify email addresses or Dynatrace user IDs here.
}

func (me *Subscriptions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"month": {
			Type:        schema.TypeSet,
			Description: "A list of monthly subscribers.\n Monthly subscribers receive the report on the first Monday of the month at midnight.\n You can specify email addresses or Dynatrace user IDs here.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"week": {
			Type:        schema.TypeSet,
			Description: "A list of weekly subscribers.\n Weekly subscribers receive the report every Monday at midnight.\n You can specify email addresses or Dynatrace user IDs here.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Subscriptions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"month": me.Month,
		"week":  me.Week,
	})
}

func (me *Subscriptions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"month": &me.Month,
		"week":  &me.Week,
	})
}
