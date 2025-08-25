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

package settings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type PipelineReference struct {
	PipelineID        *string `json:"pipelineId,omitempty"`
	BuiltinPipelineID *string `json:"builtinPipelineId,omitempty"`
	PipelineType      string  `json:"pipelineType"`
}

const BuiltinPipelineIDMaxLength = 500

func (pr *PipelineReference) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pipeline_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "ID of the pipeline. Only used if the pipeline type is \"custom\"",
			ExactlyOneOf: []string{prefix + "pipeline_id", prefix + "builtin_pipeline_id"},
		},
		"builtin_pipeline_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "ID of the pipeline. Only used if the pipeline type is \"builtin\"",
			ValidateFunc: validation.StringLenBetween(1, BuiltinPipelineIDMaxLength),
			ExactlyOneOf: []string{prefix + "pipeline_id", prefix + "builtin_pipeline_id"},
		},
		"pipeline_type": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Type of the pipeline. Must be \"custom\" or \"builtin\"",
			ValidateFunc: validation.StringInSlice([]string{"custom", "builtin"}, true),
		},
	}
}

func (pr *PipelineReference) MarshalHCL(properties hcl.Properties) error {
	err := properties.Encode("pipeline_type", pr.PipelineType)
	if err != nil {
		return err
	}

	if pr.PipelineID != nil {
		err = properties.Encode("pipeline_id", pr.PipelineID)
		if err != nil {
			return err
		}
	}

	if pr.BuiltinPipelineID != nil {
		err = properties.Encode("builtin_pipeline_id", pr.BuiltinPipelineID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pr *PipelineReference) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"pipeline_id":         &pr.PipelineID,
		"builtin_pipeline_id": &pr.BuiltinPipelineID,
		"pipeline_type":       &pr.PipelineType,
	})
}
