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

type TaskConditionOption struct {
	States map[string]Status `json:"states,omitempty"` // key/value pairs where the `key` is the name of another task and the value the status it needs to be for the current task to get executed. Possible values are `SUCCESS`, `ERROR`, `ANY`, `OK` and `NOK`
	Custom *string           `json:"custom,omitempty"` // A custom condition that needs to be met for the current task to get executed
	Else   *Else             `json:"else,omitempty"`
}

func (me *TaskConditionOption) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"states": {
			Type:        schema.TypeMap,
			Description: "key/value pairs where the `key` is the name of another task and the value the status it needs to be for the current task to get executed. Possible values are `SUCCESS`, `ERROR`, `ANY`, `OK` (Success or Skipped) and `NOK` (Error or Cancelled)",
			Required:    true,
		},
		"custom": {
			Type:        schema.TypeString,
			Description: "A custom condition that needs to be met for the current task to get executed",
			Optional:    true,
		},
		"else": {
			Type:        schema.TypeString,
			Description: "Possible values are `SKIP` and `STOP`",
			Optional:    true,
			Default:     "STOP",
		},
	}
}

func (me *TaskConditionOption) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"states": me.States,
		"custom": me.Custom,
		"else":   me.Else,
	})
}

func (me *TaskConditionOption) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"states": &me.States,
		"custom": &me.Custom,
		"else":   &me.Else,
	})
}
