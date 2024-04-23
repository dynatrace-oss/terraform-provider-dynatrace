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

type AnalyzerInputFields []*AnalyzerInputField

func (me *AnalyzerInputFields) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"analyzer_input_field": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(AnalyzerInputField).Schema()},
		},
	}
}

func (me AnalyzerInputFields) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("analyzer_input_field", me)
}

func (me *AnalyzerInputFields) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("analyzer_input_field", me)
}

type AnalyzerInputField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (me *AnalyzerInputField) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *AnalyzerInputField) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":   me.Key,
		"value": me.Value,
	})
}

func (me *AnalyzerInputField) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":   &me.Key,
		"value": &me.Value,
	})
}
