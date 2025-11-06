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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SmartscapeNodeAttributes struct {
	ExtractNode          bool                                  `json:"extractNode"`                    // Extract node
	FieldsToExtract      SmartscapeFieldExtractionEntries      `json:"fieldsToExtract,omitempty"`      // Fields to extract
	IdComponents         SmartscapeIdComponentsEntries         `json:"idComponents"`                   // ID components
	NodeIdFieldName      string                                `json:"nodeIdFieldName"`                // Node ID field name
	NodeName             *GenericValueAssignment               `json:"nodeName,omitempty"`             // Node name
	NodeType             string                                `json:"nodeType"`                       // Node type
	StaticEdgesToExtract SmartscapeStaticEdgeExtractionEntries `json:"staticEdgesToExtract,omitempty"` // Static edges to extract
}

func (me *SmartscapeNodeAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"extract_node": {
			Type:        schema.TypeBool,
			Description: "Extract node",
			Required:    true,
		},
		"fields_to_extract": {
			Type:        schema.TypeList,
			Description: "Fields to extract",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(SmartscapeFieldExtractionEntries).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"id_components": {
			Type:        schema.TypeList,
			Description: "ID components",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SmartscapeIdComponentsEntries).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"node_id_field_name": {
			Type:        schema.TypeString,
			Description: "Node ID field name",
			Required:    true,
		},
		"node_name": {
			Type:        schema.TypeList,
			Description: "Node name",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(GenericValueAssignment).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"node_type": {
			Type:        schema.TypeString,
			Description: "Node type",
			Required:    true,
		},
		"static_edges_to_extract": {
			Type:        schema.TypeList,
			Description: "Static edges to extract",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(SmartscapeStaticEdgeExtractionEntries).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *SmartscapeNodeAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"extract_node":            me.ExtractNode,
		"fields_to_extract":       me.FieldsToExtract,
		"id_components":           me.IdComponents,
		"node_id_field_name":      me.NodeIdFieldName,
		"node_name":               me.NodeName,
		"node_type":               me.NodeType,
		"static_edges_to_extract": me.StaticEdgesToExtract,
	})
}

func (me *SmartscapeNodeAttributes) HandlePreconditions() error {
	if (me.NodeName == nil) && (me.ExtractNode) {
		return fmt.Errorf("'node_name' must be specified if 'extract_node' is set to '%v'", me.ExtractNode)
	}
	if (me.NodeName != nil) && (!me.ExtractNode) {
		return fmt.Errorf("'node_name' must not be specified if 'extract_node' is set to '%v'", me.ExtractNode)
	}
	// ---- FieldsToExtract SmartscapeFieldExtractionEntries -> {"expectedValue":true,"property":"extractNode","type":"EQUALS"}
	// ---- StaticEdgesToExtract SmartscapeStaticEdgeExtractionEntries -> {"expectedValue":true,"property":"extractNode","type":"EQUALS"}
	return nil
}

func (me *SmartscapeNodeAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"extract_node":            &me.ExtractNode,
		"fields_to_extract":       &me.FieldsToExtract,
		"id_components":           &me.IdComponents,
		"node_id_field_name":      &me.NodeIdFieldName,
		"node_name":               &me.NodeName,
		"node_type":               &me.NodeType,
		"static_edges_to_extract": &me.StaticEdgesToExtract,
	})
}
