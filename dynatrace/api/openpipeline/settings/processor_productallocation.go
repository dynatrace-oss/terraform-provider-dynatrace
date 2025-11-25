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

type ProductAllocationProcessor struct {
	Processor
	Value *ValueAssignment `json:"value,omitempty"`
}

func (p *ProductAllocationProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()
	s["value"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(ValueAssignment).Schema()},
		Description: "Strategy to assign a value",
		Required:    true,
	}

	return s
}

func (p *ProductAllocationProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.Encode("value", p.Value)
}

func (p *ProductAllocationProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.Decode("value", &p.Value)
}

func (ep ProductAllocationProcessor) MarshalJSON() ([]byte, error) {
	type costAllocationProcessor CostAllocationProcessor
	return MarshalAsJSONWithType((costAllocationProcessor)(ep), ProductAllocationProcessorType)
}
