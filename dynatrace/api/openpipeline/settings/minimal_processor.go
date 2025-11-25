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

package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MinimalProcessor is the minimal set or shared fields across different processor types
// This is different to the defined Processor as TechnologyProcessor doesn't have the required Matcher and Description
type MinimalProcessor struct {
	Editable   *bool   `json:"editable,omitempty"`
	Enabled    bool    `json:"enabled"`
	Id         string  `json:"id"`
	SampleData *string `json:"sampleData,omitempty"`
}

func (p *MinimalProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Identifier of the processor. Must be unique within a stage.",
			Required:    true,
		},
		"sample_data": {
			Type:        schema.TypeString,
			Description: "Sample data related to the processor for documentation or testing",
			Optional:    true,
		},
	}
}

func (p *MinimalProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":     p.Enabled,
		"id":          p.Id,
		"sample_data": p.SampleData,
	})
}

func (p *MinimalProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":     &p.Enabled,
		"id":          &p.Id,
		"sample_data": &p.SampleData,
	})
}
