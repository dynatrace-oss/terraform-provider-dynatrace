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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Composition     PipelineGroupCompositions `json:"composition,omitempty"`     // Composition
	DisplayName     string                    `json:"displayName"`               // Display name
	MemberPipelines []string                  `json:"memberPipelines,omitempty"` // Pipelines wrapped by this group
	MemberStages    *StageConfig              `json:"memberStages"`              // stage configuration of the member pipelines
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"composition": {
			Type:        schema.TypeList,
			Description: "Composition",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(PipelineGroupCompositions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name",
			Required:    true,
		},
		"member_pipelines": {
			Type:        schema.TypeSet,
			Description: "Pipelines wrapped by this group",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"member_stages": {
			Type:        schema.TypeList,
			Description: "stage configuration of the member pipelines",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(StageConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"composition":      me.Composition,
		"display_name":     me.DisplayName,
		"member_pipelines": me.MemberPipelines,
		"member_stages":    me.MemberStages,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"composition":      &me.Composition,
		"display_name":     &me.DisplayName,
		"member_pipelines": &me.MemberPipelines,
		"member_stages":    &me.MemberStages,
	})
}
