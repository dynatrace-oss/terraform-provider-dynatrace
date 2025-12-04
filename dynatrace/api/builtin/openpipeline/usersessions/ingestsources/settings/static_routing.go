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

package ingestsources

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type StaticRouting struct {
	BuiltinPipelineID *string      `json:"builtinPipelineId,omitempty"` // Builtin Pipeline ID
	PipelineID        *string      `json:"pipelineId,omitempty"`        // Pipeline ID
	PipelineType      PipelineType `json:"pipelineType"`                // Pipeline Type. Possible Values: `builtin`, `custom`
}

func (me *StaticRouting) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"builtin_pipeline_id": {
			Type:        schema.TypeString,
			Description: "Builtin Pipeline ID",
			Optional:    true, // precondition
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

func (me *StaticRouting) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"builtin_pipeline_id": me.BuiltinPipelineID,
		"pipeline_id":         me.PipelineID,
		"pipeline_type":       me.PipelineType,
	})
}

func (me *StaticRouting) HandlePreconditions() error {
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

func (me *StaticRouting) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"builtin_pipeline_id": &me.BuiltinPipelineID,
		"pipeline_id":         &me.PipelineID,
		"pipeline_type":       &me.PipelineType,
	})
}
