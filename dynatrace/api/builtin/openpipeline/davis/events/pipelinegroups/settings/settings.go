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

type Settings struct {
	DisplayName       string         `json:"displayName"`                 // Display name
	IncludedPipelines PipelineGroups `json:"includedPipelines,omitempty"` // Included Pipelines
	TargetedPipelines []string       `json:"targetedPipelines,omitempty"` // Pipelines wrapped by this group
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name",
			Required:    true,
		},
		"included_pipelines": {
			Type:        schema.TypeList,
			Description: "Included Pipelines",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(PipelineGroups).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"targeted_pipelines": {
			Type:        schema.TypeSet,
			Description: "Pipelines wrapped by this group",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"display_name":       me.DisplayName,
		"included_pipelines": me.IncludedPipelines,
		"targeted_pipelines": me.TargetedPipelines,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"display_name":       &me.DisplayName,
		"included_pipelines": &me.IncludedPipelines,
		"targeted_pipelines": &me.TargetedPipelines,
	})
}
