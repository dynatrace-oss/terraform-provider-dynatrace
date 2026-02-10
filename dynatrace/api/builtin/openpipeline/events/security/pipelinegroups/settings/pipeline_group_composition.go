/**
* @license
* Copyright 2026 Dynatrace LLC
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

package pipelinegroups

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PipelineGroupCompositions []*PipelineGroupComposition

func (me *PipelineGroupCompositions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pipeline_group_composition": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(PipelineGroupComposition).Schema()},
		},
	}
}

func (me PipelineGroupCompositions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("pipeline_group_composition", me)
}

func (me *PipelineGroupCompositions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("pipeline_group_composition", me)
}

type PipelineGroupComposition struct {
	IsPipelinePlaceholder bool         `json:"isPipelinePlaceholder"` // Placeholder for the wrapped pipeline
	PipelineID            *string      `json:"pipelineId,omitempty"`  // Pipeline ID
	Stages                *StageConfig `json:"stages,omitempty"`      // stage configuration for this pipelines
}

func (me *PipelineGroupComposition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"is_pipeline_placeholder": {
			Type:        schema.TypeBool,
			Description: "Placeholder for the wrapped pipeline",
			Required:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Description: "Pipeline ID",
			Optional:    true, // precondition
		},
		"stages": {
			Type:        schema.TypeList,
			Description: "stage configuration for this pipelines",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(StageConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *PipelineGroupComposition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"is_pipeline_placeholder": me.IsPipelinePlaceholder,
		"pipeline_id":             me.PipelineID,
		"stages":                  me.Stages,
	})
}

func (me *PipelineGroupComposition) HandlePreconditions() error {
	if (me.PipelineID == nil) && (!me.IsPipelinePlaceholder) {
		return fmt.Errorf("'pipeline_id' must be specified if 'is_pipeline_placeholder' is set to '%v'", me.IsPipelinePlaceholder)
	}
	if (me.PipelineID != nil) && (me.IsPipelinePlaceholder) {
		return fmt.Errorf("'pipeline_id' must not be specified if 'is_pipeline_placeholder' is set to '%v'", me.IsPipelinePlaceholder)
	}
	if (me.Stages == nil) && (!me.IsPipelinePlaceholder) {
		return fmt.Errorf("'stages' must be specified if 'is_pipeline_placeholder' is set to '%v'", me.IsPipelinePlaceholder)
	}
	if (me.Stages != nil) && (me.IsPipelinePlaceholder) {
		return fmt.Errorf("'stages' must not be specified if 'is_pipeline_placeholder' is set to '%v'", me.IsPipelinePlaceholder)
	}
	return nil
}

func (me *PipelineGroupComposition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"is_pipeline_placeholder": &me.IsPipelinePlaceholder,
		"pipeline_id":             &me.PipelineID,
		"stages":                  &me.Stages,
	})
}
