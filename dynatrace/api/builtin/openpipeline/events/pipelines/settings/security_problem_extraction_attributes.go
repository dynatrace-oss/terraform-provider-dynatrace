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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecurityProblemExtractionAttributes struct {
	Category               string                          `json:"category"`
	FieldsCopiedFromRecord []string                        `json:"fieldsCopiedFromRecord,omitempty"` // Record fields copied verbatim onto the extracted security problem.
	IdComponents           SmartscapeIdComponentsEntries   `json:"idComponents"`                     // ID components
	NodeName               *GenericValueAssignment         `json:"nodeName,omitempty"`               // Node name
	SecurityContext        *SecurityProblemSecurityContext `json:"securityContext,omitempty"`        // Security context written on the extracted security problem. It controls who is allowed to read the security problem. If not set, no security context is written.
}

func (me *SecurityProblemExtractionAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"category": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Required:    true,
		},
		"fields_copied_from_record": {
			Type:        schema.TypeList,
			Description: "Record fields copied verbatim onto the extracted security problem.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"id_components": {
			Type:        schema.TypeList,
			Description: "ID components",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SmartscapeIdComponentsEntries).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"node_name": {
			Type:        schema.TypeList,
			Description: "Node name",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(GenericValueAssignment).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"security_context": {
			Type:        schema.TypeList,
			Description: "Security context written on the extracted security problem. It controls who is allowed to read the security problem. If not set, no security context is written.",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(SecurityProblemSecurityContext).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *SecurityProblemExtractionAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"category":                  me.Category,
		"fields_copied_from_record": me.FieldsCopiedFromRecord,
		"id_components":             me.IdComponents,
		"node_name":                 me.NodeName,
		"security_context":          me.SecurityContext,
	})
}

func (me *SecurityProblemExtractionAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"category":                  &me.Category,
		"fields_copied_from_record": &me.FieldsCopiedFromRecord,
		"id_components":             &me.IdComponents,
		"node_name":                 &me.NodeName,
		"security_context":          &me.SecurityContext,
	})
}
