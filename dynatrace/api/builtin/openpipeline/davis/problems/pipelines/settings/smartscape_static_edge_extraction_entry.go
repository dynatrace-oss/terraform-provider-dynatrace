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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SmartscapeStaticEdgeExtractionEntries []*SmartscapeStaticEdgeExtractionEntry

func (me *SmartscapeStaticEdgeExtractionEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"smartscape_static_edge_extraction_entry": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SmartscapeStaticEdgeExtractionEntry).Schema()},
		},
	}
}

func (me SmartscapeStaticEdgeExtractionEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("smartscape_static_edge_extraction_entry", me)
}

func (me *SmartscapeStaticEdgeExtractionEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("smartscape_static_edge_extraction_entry", me)
}

type SmartscapeStaticEdgeExtractionEntry struct {
	EdgeType          string `json:"edgeType"`          // Edge type
	TargetIdFieldName string `json:"targetIdFieldName"` // Target ID field name
	TargetType        string `json:"targetType"`        // Target type
}

func (me *SmartscapeStaticEdgeExtractionEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"edge_type": {
			Type:        schema.TypeString,
			Description: "Edge type",
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

func (me *SmartscapeStaticEdgeExtractionEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"edge_type":            me.EdgeType,
		"target_id_field_name": me.TargetIdFieldName,
		"target_type":          me.TargetType,
	})
}

func (me *SmartscapeStaticEdgeExtractionEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"edge_type":            &me.EdgeType,
		"target_id_field_name": &me.TargetIdFieldName,
		"target_type":          &me.TargetType,
	})
}
