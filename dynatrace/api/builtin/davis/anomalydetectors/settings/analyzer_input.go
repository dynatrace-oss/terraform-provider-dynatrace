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

package anomalydetectors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AnalyzerInput struct {
	Input AnalyzerInputFields `json:"input,omitempty"` // Input fields for the specified analyzer
	Name  string              `json:"name"`            // Fully qualified name of the analyzer
}

func (me *AnalyzerInput) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"input": {
			Type:        schema.TypeList,
			Description: "Input fields for the specified analyzer",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(AnalyzerInputFields).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Fully qualified name of the analyzer",
			Required:    true,
		},
	}
}

func (me *AnalyzerInput) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"input": me.Input,
		"name":  me.Name,
	})
}

func (me *AnalyzerInput) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"input": &me.Input,
		"name":  &me.Name,
	})
}
