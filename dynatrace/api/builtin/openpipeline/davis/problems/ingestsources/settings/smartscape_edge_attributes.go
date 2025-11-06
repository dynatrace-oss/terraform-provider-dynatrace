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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SmartscapeEdgeAttributes struct {
	EdgeType          string `json:"edgeType"`          // Edge type
	SourceIdFieldName string `json:"sourceIdFieldName"` // Source ID field name
	SourceType        string `json:"sourceType"`        // Source type
	TargetIdFieldName string `json:"targetIdFieldName"` // Target ID field name
	TargetType        string `json:"targetType"`        // Target type
}

func (me *SmartscapeEdgeAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"edge_type": {
			Type:        schema.TypeString,
			Description: "Edge type",
			Required:    true,
		},
		"source_id_field_name": {
			Type:        schema.TypeString,
			Description: "Source ID field name",
			Required:    true,
		},
		"source_type": {
			Type:        schema.TypeString,
			Description: "Source type",
			Required:    true,
		},
		"target_id_field_name": {
			Type:        schema.TypeString,
			Description: "Target ID field name",
			Required:    true,
		},
		"target_type": {
			Type:        schema.TypeString,
			Description: "Target type",
			Required:    true,
		},
	}
}

func (me *SmartscapeEdgeAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"edge_type":            me.EdgeType,
		"source_id_field_name": me.SourceIdFieldName,
		"source_type":          me.SourceType,
		"target_id_field_name": me.TargetIdFieldName,
		"target_type":          me.TargetType,
	})
}

func (me *SmartscapeEdgeAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"edge_type":            &me.EdgeType,
		"source_id_field_name": &me.SourceIdFieldName,
		"source_type":          &me.SourceType,
		"target_id_field_name": &me.TargetIdFieldName,
		"target_type":          &me.TargetType,
	})
}
