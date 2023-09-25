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

package scheduling_rules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GroupingRule struct {
	Combine   []string `json:"combine"`             // The IDs of scheduling rules determining the days the schedule should apply to
	Intersect []string `json:"intersect,omitempty"` // The IDs of scheduling rules determining the days the schedule is allowed apply to. If specified, only days that are covered by `combine` and `intersect` are valid days for the schedule
	Subtract  []string `json:"subtract,omitempty"`  // The IDs of scheduling rules determing the days the schedule must not apply. If specified it reduces down the set of days covered by `combine` and `intersect`
}

func (me *GroupingRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"combine": {
			Type:        schema.TypeSet,
			Description: "The IDs of scheduling rules determining the days the schedule should apply to",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"intersect": {
			Type:        schema.TypeSet,
			Description: "The IDs of scheduling rules determining the days the schedule is allowed apply to. If specified, only days that are covered by `combine` and `intersect` are valid days for the schedule",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"subtract": {
			Type:        schema.TypeSet,
			Description: "The IDs of scheduling rules determing the days the schedule must not apply. If specified it reduces down the set of days covered by `combine` and `intersect`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *GroupingRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"combine":   me.Combine,
		"intersect": me.Intersect,
		"subtract":  me.Subtract,
	})
}

func (me *GroupingRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"combine":   &me.Combine,
		"intersect": &me.Intersect,
		"subtract":  &me.Subtract,
	})
}
