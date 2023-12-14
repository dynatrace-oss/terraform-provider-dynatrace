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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SchemaVersion = 3

type Workflow struct {
	Title       string `json:"title" maxlength:"200"` // The title / name of the workflow
	Description string `json:"description,omitempty"` // An optional description for the workflow

	Actor         string   `json:"actor,omitempty" maxlength:"36" format:"uuid"` // The user context the executions of the workflow will happen with
	Owner         string   `json:"owner,omitempty" format:"uuid"`                // The ID of the owner of this workflow
	Private       bool     `json:"isPrivate" default:"true"`                     // Defines whether this workflow is private to the owner or not. Default is `true`
	SchemaVersion int      `json:"schemaVersion,omitempty"`                      //
	Trigger       *Trigger `json:"trigger,omitempty"`                            // Configures how executions of the workflows are getting triggered. If no trigger is specified it means the workflow is getting manually triggered
	Tasks         Tasks    `json:"tasks,omitempty"`                              // The tasks to run for every execution of this workflow
}

func (me *Workflow) Name() string {
	return me.Title
}

func (me *Workflow) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:             schema.TypeString,
			Description:      "The title / name of the workflow",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(200),
		},
		"description": {
			Type:        schema.TypeString,
			Description: "An optional description for the workflow",
			Optional:    true,
		},

		"actor": {
			Type:             schema.TypeString,
			Description:      "The user context the executions of the workflow will happen with",
			Optional:         true,
			ValidateDiagFunc: Validate(ValidateUUID, ValidateMaxLength(36)),
		},
		"owner": {
			Type:             schema.TypeString,
			Description:      "The ID of the owner of this workflow",
			Optional:         true,
			ValidateDiagFunc: ValidateUUID,
		},
		"private": {
			Type:        schema.TypeBool,
			Description: "Defines whether this workflow is private to the owner or not. Default is `true`",
			Optional:    true,
			Default:     true,
		},
		"trigger": {
			Type:        schema.TypeList,
			Description: "Configures how executions of the workflows are getting triggered. If no trigger is specified it means the workflow is getting manually triggered",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Trigger).Schema("trigger")},
		},
		"tasks": {
			Type:        schema.TypeList,
			Description: "The tasks to run for every execution of this workflow",
			MinItems:    1,
			MaxItems:    1,
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Tasks).Schema("tasks")},
		},
	}
}

func (me *Workflow) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"title":       me.Title,
		"description": me.Description,

		"actor":   me.Actor,
		"owner":   me.Owner,
		"private": me.Private,
		"trigger": me.Trigger,
		"tasks":   me.Tasks,
	})
}

func (me *Workflow) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"title":       &me.Title,
		"description": &me.Description,

		"actor":   &me.Actor,
		"owner":   &me.Owner,
		"private": &me.Private,
		"trigger": &me.Trigger,
		"tasks":   &me.Tasks,
	})
}

func (me *Workflow) MarshalJSON() ([]byte, error) {
	wf := struct {
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`

		Actor         string   `json:"actor,omitempty"`
		Owner         string   `json:"owner,omitempty"`
		Private       bool     `json:"isPrivate"`
		SchemaVersion int      `json:"schemaVersion,omitempty"`
		Trigger       *Trigger `json:"trigger,omitempty"`
		Tasks         Tasks    `json:"tasks,omitempty"`
	}{
		SchemaVersion: SchemaVersion, // adding the Schema Version is the purpose of this custome `MarshalJSON` function
		Title:         me.Title,
		Description:   me.Description,
		Actor:         me.Actor,
		Owner:         me.Owner,
		Private:       me.Private,
		Trigger:       me.Trigger,
		Tasks:         me.Tasks,
	}
	return json.Marshal(wf)
}
