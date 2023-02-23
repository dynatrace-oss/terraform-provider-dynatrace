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

package workloaddetection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type InclusionToggles struct {
	IncBasepod   bool `json:"incBasepod"`   // E.g. \"cloud-credential-operator-\" for \"cloud-credential-operator-5ff6dbff57-gszgq\"
	IncContainer bool `json:"incContainer"` // Container name
	IncNamespace bool `json:"incNamespace"` // Namespace name
	IncProduct   bool `json:"incProduct"`   // If Product is enabled and has no value, it defaults to Base pod name
	IncStage     bool `json:"incStage"`     // Stage
}

func (me *InclusionToggles) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"inc_basepod": {
			Type:        schema.TypeBool,
			Description: "E.g. \"cloud-credential-operator-\" for \"cloud-credential-operator-5ff6dbff57-gszgq\"",
			Required:    true,
		},
		"inc_container": {
			Type:        schema.TypeBool,
			Description: "Container name",
			Required:    true,
		},
		"inc_namespace": {
			Type:        schema.TypeBool,
			Description: "Namespace name",
			Required:    true,
		},
		"inc_product": {
			Type:        schema.TypeBool,
			Description: "If Product is enabled and has no value, it defaults to Base pod name",
			Required:    true,
		},
		"inc_stage": {
			Type:        schema.TypeBool,
			Description: "Stage",
			Required:    true,
		},
	}
}

func (me *InclusionToggles) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"inc_basepod":   me.IncBasepod,
		"inc_container": me.IncContainer,
		"inc_namespace": me.IncNamespace,
		"inc_product":   me.IncProduct,
		"inc_stage":     me.IncStage,
	})
}

func (me *InclusionToggles) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"inc_basepod":   &me.IncBasepod,
		"inc_container": &me.IncContainer,
		"inc_namespace": &me.IncNamespace,
		"inc_product":   &me.IncProduct,
		"inc_stage":     &me.IncStage,
	})
}
