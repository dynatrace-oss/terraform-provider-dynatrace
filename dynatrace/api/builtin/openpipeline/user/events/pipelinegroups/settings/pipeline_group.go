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

package pipelinegroups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PipelineGroups []*PipelineGroup

func (me *PipelineGroups) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"included_pipeline": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(PipelineGroup).Schema()},
		},
	}
}

func (me PipelineGroups) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("included_pipeline", me)
}

func (me *PipelineGroups) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("included_pipeline", me)
}

type PipelineGroup struct {
	EnabledStages               []StageType `json:"enabledStages,omitempty"`     // Enabled Stages. Possible Values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`.
	IsTargetPipelinePlaceholder bool        `json:"isTargetPipelinePlaceholder"` // Placeholder for the wrapped Pipeline
	PipelineID                  *string     `json:"pipelineId,omitempty"`        // Pipeline ID
}

func (me *PipelineGroup) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled_stages": {
			Type:        schema.TypeSet,
			Description: "Enabled Stages. Possible Values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"is_target_pipeline_placeholder": {
			Type:        schema.TypeBool,
			Description: "Placeholder for the wrapped Pipeline",
			Required:    true,
		},
		"pipeline_id": {
			Type:        schema.TypeString,
			Description: "Pipeline ID",
			Optional:    true, // nullable
		},
	}
}

func (me *PipelineGroup) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled_stages":                 me.EnabledStages,
		"is_target_pipeline_placeholder": me.IsTargetPipelinePlaceholder,
		"pipeline_id":                    me.PipelineID,
	})
}

func (me *PipelineGroup) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled_stages":                 &me.EnabledStages,
		"is_target_pipeline_placeholder": &me.IsTargetPipelinePlaceholder,
		"pipeline_id":                    &me.PipelineID,
	})
}
