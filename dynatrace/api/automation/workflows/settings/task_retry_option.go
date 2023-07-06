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

type TaskRetryOption struct {
	Count                    VarInt  `json:"count" minimum:"0" maximum:"99"`                    // Specifies a maximum number of times that a task can be repeated in case it fails on execution
	Delay                    *VarInt `json:"delay,omitempty" minimum:"0" maximum:"3600"`        // Specifies a delay in seconds between subsequent task retries
	FailedLoopIterationsOnly *bool   `json:"failedLoopIterationsOnly,omitempty" default:"true"` // Specifies whether retrying the failed iterations or the whole loop. Default: True
}

func (me *TaskRetryOption) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"count": { // TODO: minimum:"0" maximum:"99"
			Type:        schema.TypeString,
			Description: "Specifies a maximum number of times that a task can be repeated in case it fails on execution. You can specify either a number between 1 and 99 here or use an expression (`{{}}`). Default: 1",
			Default:     1,
			Optional:    true,
		},
		"delay": { // TODO: minimum:"0" maximum:"3600"
			Type:        schema.TypeString,
			Description: "Specifies a delay in seconds between subsequent task retries. You can specify either a number between 1 and 3600 here or an expression (`{{...}}`). Default: 1",
			Optional:    true,
		},
		"failed_loop_iterations_only": {
			Type:        schema.TypeBool,
			Description: "Specifies whether retrying the failed iterations or the whole loop. Default: true",
			Optional:    true,
			Default:     true,
		},
	}
}

func (me *TaskRetryOption) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"count":                       me.Count,
		"delay":                       me.Delay,
		"failed_loop_iterations_only": me.FailedLoopIterationsOnly,
	})
}

func (me *TaskRetryOption) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"count":                       &me.Count,
		"delay":                       &me.Delay,
		"failed_loop_iterations_only": &me.FailedLoopIterationsOnly,
	})
}
