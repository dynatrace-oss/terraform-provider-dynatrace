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

type SecurityProblemSecurityContext struct {
	Type   SecurityContextType `json:"type"`             // Whether the Dynatrace-defined default security context or a custom one is applied. Possible values: `custom`, `default`
	Values []string            `json:"values,omitempty"` // Security context
}

func (me *SecurityProblemSecurityContext) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Whether the Dynatrace-defined default security context or a custom one is applied. Possible values: `custom`, `default`",
			Required:    true,
		},
		"values": {
			Type:        schema.TypeList,
			Description: "Security context",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *SecurityProblemSecurityContext) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":   me.Type,
		"values": me.Values,
	})
}

func (me *SecurityProblemSecurityContext) HandlePreconditions() error {
	// ---- Values []string -> {"expectedValue":"custom","property":"type","type":"EQUALS"}
	return nil
}

func (me *SecurityProblemSecurityContext) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":   &me.Type,
		"values": &me.Values,
	})
}
