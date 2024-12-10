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

package fields

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled      bool   `json:"enabled"`      // This setting is enabled (`true`) or disabled (`false`)
	EventField   string `json:"eventField"`   // Field from the event that will be extracted.
	ProblemField string `json:"problemField"` // Field under which the extracted event data will be stored on the problem.
}

func (me *Settings) Name() string {
	return "environment"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"event_field": {
			Type:        schema.TypeString,
			Description: "Field from the event that will be extracted.",
			Required:    true,
		},
		"problem_field": {
			Type:        schema.TypeString,
			Description: "Field under which the extracted event data will be stored on the problem.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":       me.Enabled,
		"event_field":   me.EventField,
		"problem_field": me.ProblemField,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":       &me.Enabled,
		"event_field":   &me.EventField,
		"problem_field": &me.ProblemField,
	})
}
