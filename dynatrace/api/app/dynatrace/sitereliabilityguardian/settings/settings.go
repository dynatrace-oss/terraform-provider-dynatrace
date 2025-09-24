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

package sitereliabilityguardian

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Description *string    `json:"description,omitempty"` // Description
	EventKind   *EventKind `json:"eventKind,omitempty"`   // If set to null/'BIZ_EVENT' validation events stored as bizevents in Grail. If set to 'SDLC_EVENT' validation events stored as SDLC events
	Name        string     `json:"name"`                  // Name
	Objectives  Objectives `json:"objectives"`            // Objectives
	Tags        []string   `json:"tags,omitempty"`        // Define key/value pairs that further describe this guardian.
	Variables   Variables  `json:"variables,omitempty"`   // Define variables for dynamically defining DQL queries
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"event_kind": {
			Type:        schema.TypeString,
			Description: "If set to null/'BIZ_EVENT' validation events stored as bizevents in Grail. If set to 'SDLC_EVENT' validation events stored as SDLC events",
			Optional:    true, // nullable
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"objectives": {
			Type:        schema.TypeList,
			Description: "Objectives",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Objectives).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "Define key/value pairs that further describe this guardian.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"variables": {
			Type:        schema.TypeList,
			Description: "Define variables for dynamically defining DQL queries",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Variables).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description": me.Description,
		"event_kind":  me.EventKind,
		"name":        me.Name,
		"objectives":  me.Objectives,
		"tags":        me.Tags,
		"variables":   me.Variables,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description": &me.Description,
		"event_kind":  &me.EventKind,
		"name":        &me.Name,
		"objectives":  &me.Objectives,
		"tags":        &me.Tags,
		"variables":   &me.Variables,
	})
}
