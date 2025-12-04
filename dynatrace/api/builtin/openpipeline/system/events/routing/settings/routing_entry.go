/**
* @license
* Copyright 2025 Dynatrace LLC
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

package routing

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RoutingEntries []*RoutingEntry

func (me *RoutingEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"routing_entry": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(RoutingEntry).Schema()},
		},
	}
}

func (me RoutingEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("routing_entry", me)
}

func (me *RoutingEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("routing_entry", me)
}

type RoutingEntry struct {
	BuiltinPipelineID *string      `json:"builtinPipelineId,omitempty"` // Builtin Pipeline ID
	Description       string       `json:"description"`
	Enabled           bool         `json:"enabled"`              // This setting is enabled (`true`) or disabled (`false`)
	Matcher           string       `json:"matcher"`              // Query which determines whether the record should be routed to the target pipeline of this rule.
	PipelineID        *string      `json:"pipelineId,omitempty"` // Pipeline ID
	PipelineType      PipelineType `json:"pipelineType"`         // Pipeline Type. Possible Values: `builtin`, `custom`
}

func (me *RoutingEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"builtin_pipeline_id": {
			Type:        schema.TypeString,
			Description: "Builtin Pipeline ID",
			Optional:    true, // precondition
		},
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Query which determines whether the record should be routed to the target pipeline of this rule.",
			Required:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Description: "Pipeline ID",
			Optional:    true, // precondition
		},
		"pipeline_type": {
			Type:        schema.TypeString,
			Description: "Pipeline Type. Possible Values: `builtin`, `custom`",
			Required:    true,
		},
	}
}

func (me *RoutingEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"builtin_pipeline_id": me.BuiltinPipelineID,
		"description":         me.Description,
		"enabled":             me.Enabled,
		"matcher":             me.Matcher,
		"pipeline_id":         me.PipelineID,
		"pipeline_type":       me.PipelineType,
	})
}

func (me *RoutingEntry) HandlePreconditions() error {
	if (me.PipelineID == nil) && (string(me.PipelineType) == "custom") {
		me.PipelineID = opt.NewString("")
	}
	if (me.BuiltinPipelineID == nil) && (string(me.PipelineType) == "builtin") {
		return fmt.Errorf("'builtin_pipeline_id' must be specified if 'pipeline_type' is set to '%v'", me.PipelineType)
	}
	if (me.BuiltinPipelineID != nil) && (string(me.PipelineType) != "builtin") {
		return fmt.Errorf("'builtin_pipeline_id' must not be specified if 'pipeline_type' is set to '%v'", me.PipelineType)
	}
	return nil
}

func (me *RoutingEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"builtin_pipeline_id": &me.BuiltinPipelineID,
		"description":         &me.Description,
		"enabled":             &me.Enabled,
		"matcher":             &me.Matcher,
		"pipeline_id":         &me.PipelineID,
		"pipeline_type":       &me.PipelineType,
	})
}
