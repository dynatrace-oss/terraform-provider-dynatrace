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

package logdpprules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProcessorDefinition struct {
	Rule string `json:"rule"` // Processor definition
}

func (me *ProcessorDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:             schema.TypeString,
			Description:      "Processor definition",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
	}
}

func (me *ProcessorDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"rule": me.Rule,
	})
}

func (me *ProcessorDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"rule": &me.Rule,
	})
}
