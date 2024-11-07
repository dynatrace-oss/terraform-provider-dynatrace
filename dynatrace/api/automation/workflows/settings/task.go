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
	"encoding/json"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Tasks []*Task

func (me Tasks) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	for _, task := range me {
		// Because of a bug in Terraform regarding Sets of Resources
		// These extra entries may occur
		if len(task.Name) == 0 {
			continue
		}
		data, err := json.Marshal(task)
		if err != nil {
			return nil, err
		}
		m[task.Name] = data
	}
	return json.Marshal(m)
}

func (me *Tasks) UnmarshalJSON(data []byte) error {
	m := map[string]*Task{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for task_name, task := range m {
		if len(task.Name) == 0 {
			task.Name = task_name
		}
		*me = append(*me, task)
	}
	return nil
}

func (me *Tasks) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"task": {
			Type:        schema.TypeSet,
			Description: "TODO: No documentation available",
			MinItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Task).Schema(prefix + ".0.task")},
		},
	}
}

func (me Tasks) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("task", me)
}

func (me *Tasks) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("task", me)
}

type Task struct {
	Name         string               `json:"name" pattern:"^(?!.*null$)([A-Za-z_][A-Za-z0-9-_]*)$"` // The name of the task
	Action       string               `json:"action" pattern:"^.+:.+$"`
	Description  *string              `json:"description,omitempty"` // A description for this task
	Input        map[string]any       `json:"input"`
	Active       bool                 `json:"active" default:"true"` // Specifies whether a task should be skipped as a no operation or not
	Position     *TaskPosition        `json:"position"`
	Predecessors []string             `json:"predecessors,omitempty"`
	Conditions   *TaskConditionOption `json:"conditions,omitempty"`
	WithItems    *string              `json:"withItems,omitempty"`
	Concurrency  *VarInt              `json:"concurrency" minimum:"1" maximum:"99"`
	Retry        *TaskRetryOption     `json:"retry,omitempty"`
	Timeout      *VarInt              `json:"timeout" default:"900" minimum:"1" maximum:"604800"` // Specifies a default task timeout. 15 * 60 (15min) is used when not set
	WaitBefore   *VarInt              `json:"waitBefore" default:"0" minimum:"0" maximum:"86400"` // Specifies a default task wait before in seconds. 0 is used when not set
}

func (me *Task) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Description:      "The name of the task",
			Required:         true,
			ValidateDiagFunc: ValidateRegex(regexp.MustCompile("^([A-Za-z_][A-Za-z0-9-_]*)$"), "Only letters, digits, _ and - are allowed. The name must start with a letter."),
		},
		"action": {
			Type:        schema.TypeString,
			Description: "Currently known and supported values are `dynatrace.automations:http-function`, `dynatrace.automations:run-javascript` and `dynatrace.automations:execute-dql-query`",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "A description for this task",
			Optional:    true,
		},
		"input": {
			Type:             schema.TypeString,
			Description:      "Parameters and values for this task as JSON code. Contents depend on the kind of task - determined by the attribute `action`",
			Optional:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
			Elem:             &schema.Schema{Type: schema.TypeString},
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "Specifies whether a task should be skipped as a no operation or not",
			Optional:    true,
			Default:     true,
		},
		"position": {
			Type:        schema.TypeList,
			Description: "Layouting information about the task tile when visualized. If not specified Dynatrace will position the task tiles automatically",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TaskPosition).Schema(prefix + ".0.position")},
		},
		"conditions": {
			Type:        schema.TypeList,
			Description: "Conditions that have to be met in order to execute that task",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TaskConditionOption).Schema(prefix + ".0.conditions")},
		},
		"with_items": {
			Type:        schema.TypeString,
			Description: "Iterates over items in a list, allowing actions to be executed repeatedly. Example: Specifying `item in [1, 2, 3]` here will execute the task three times for the numbers 1, 2 and 3 - with the current number available for scripting using the expression `{{ _.item }}`",
			Optional:    true,
		},
		"concurrency": {
			Type:             schema.TypeString,
			Description:      "Required if `with_items` is specified. By default loops execute sequentially with concurrency set to 1. You can increase how often it runs in parallel",
			Optional:         true,
			ValidateDiagFunc: ValidateRange(1, 99),
		},
		"retry": {
			Type:        schema.TypeList,
			Description: "Configure whether to automatically rerun the task on failure. If not specified no retries will be attempted",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TaskRetryOption).Schema(prefix + ".0.retry")},
		},
		"timeout": {
			Type:             schema.TypeString,
			Description:      "Specifies a default task timeout in seconds. 15 * 60 (15min) is used when not set. Minimum 1. Maximum 604800",
			Optional:         true,
			ValidateDiagFunc: ValidateRange(1, 604800),
			Default:          "900",
		},
		"wait_before": {
			Type:             schema.TypeString,
			Description:      "Specifies a default task wait before in seconds. 0 is used when not set",
			Optional:         true,
			ValidateDiagFunc: ValidateRange(0, 86400),
			Default:          "0",
		},
	}
}

func (me *Task) MarshalHCL(properties hcl.Properties) error {
	var inputJSON *string
	var err error
	if len(me.Input) > 0 {
		var data []byte
		data, err = json.Marshal(me.Input)
		if err != nil {
			return err
		}
		inputJSON = opt.NewString(string(data))
	}
	// Fix for #579. The REST API apparently now produces an empty `conditions` property
	// That leads to non-empty plans
	if me.Conditions != nil {
		if me.Conditions.Custom == nil && me.Conditions.Else == nil && len(me.Conditions.States) == 0 {
			me.Conditions = nil
		}
	}
	return properties.EncodeAll(map[string]any{
		"name":        me.Name,
		"action":      me.Action,
		"description": me.Description,
		"input":       inputJSON,
		"active":      me.Active,
		"position":    me.Position,
		"conditions":  me.Conditions,
		"with_items":  me.WithItems,
		"retry":       me.Retry,
		"timeout":     me.Timeout,
		"wait_before": me.WaitBefore,
	})
}

func (me *Task) UnmarshalHCL(decoder hcl.Decoder) error {
	var inputStr string
	if err := decoder.DecodeAll(map[string]any{
		"name":        &me.Name,
		"action":      &me.Action,
		"description": &me.Description,
		"input":       &inputStr,
		"active":      &me.Active,
		"position":    &me.Position,
		"conditions":  &me.Conditions,
		"with_items":  &me.WithItems,
		"retry":       &me.Retry,
		"timeout":     &me.Timeout,
		"wait_before": &me.WaitBefore,
	}); err != nil {
		return err
	}
	if len(inputStr) > 0 {
		if err := json.Unmarshal([]byte(inputStr), &me.Input); err != nil {
			return err
		}
	}

	// The REST API requires `predecessors` getting populated with the names of the tasks
	// mentioned within `conditions.states` (where the task names are the keys).
	//
	if me.Conditions != nil {
		if len(me.Conditions.States) > 0 {
			me.Predecessors = []string{}
			for k := range me.Conditions.States {
				me.Predecessors = append(me.Predecessors, k)
			}
		}
	}
	return nil
}
