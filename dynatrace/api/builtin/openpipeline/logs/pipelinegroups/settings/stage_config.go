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

type StageConfig struct {
	Exclude []StageType     `json:"exclude,omitempty"` // exclude stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
	Include []StageType     `json:"include,omitempty"` // include stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`
	Type    StageConfigType `json:"type"`              // Stage configuration type. Possible values: `exclude`, `include`, `includeAll`
}

func (me *StageConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"exclude": {
			Type:        schema.TypeSet,
			Description: "exclude stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"include": {
			Type:        schema.TypeSet,
			Description: "include stages. Possible values: `costAllocation`, `dataExtraction`, `davis`, `metricExtraction`, `processing`, `productAllocation`, `securityContext`, `smartscapeEdgeExtraction`, `smartscapeNodeExtraction`, `storage`",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Stage configuration type. Possible values: `exclude`, `include`, `includeAll`",
			Required:    true,
		},
	}
}

func (me *StageConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"exclude": me.Exclude,
		"include": me.Include,
		"type":    me.Type,
	})
}

func (me *StageConfig) HandlePreconditions() error {
	// ---- Exclude []StageType -> {"expectedValue":"exclude","property":"type","type":"EQUALS"}
	// ---- Include []StageType -> {"expectedValue":"include","property":"type","type":"EQUALS"}
	return nil
}

func (me *StageConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"exclude": &me.Exclude,
		"include": &me.Include,
		"type":    &me.Type,
	})
}
