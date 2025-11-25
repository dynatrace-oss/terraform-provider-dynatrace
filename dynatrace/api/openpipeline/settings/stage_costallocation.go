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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CostAllocationStage struct {
	Editable   *bool                           `json:"editable,omitempty"`
	Processors []*CostAllocationStageProcessor `json:"processors"`
}

func (f *CostAllocationStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Cost allocation processor to use",
			Elem:        &schema.Resource{Schema: new(CostAllocationStageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (f *CostAllocationStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", f.Processors)
}

func (f *CostAllocationStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &f.Processors)
}

type CostAllocationStageProcessor struct {
	costAllocationProcessor *CostAllocationProcessor
}

func (ep *CostAllocationStageProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cost_allocation_processor": {
			Type:        schema.TypeList,
			Description: "Processor to write the occurrences as a cost allocation",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(CostAllocationProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *CostAllocationStageProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cost_allocation_processor": ep.costAllocationProcessor,
	})
}

func (ep *CostAllocationStageProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cost_allocation_processor": &ep.costAllocationProcessor,
	})
}

func (ep CostAllocationStageProcessor) MarshalJSON() ([]byte, error) {
	if ep.costAllocationProcessor != nil {
		return json.Marshal(ep.costAllocationProcessor)
	}

	return nil, errors.New("missing CostAllocationProcessor value")
}

func (ep *CostAllocationStageProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case CostAllocationProcessorType:
		costAllocationProcessor := CostAllocationProcessor{}
		if err := json.Unmarshal(b, &costAllocationProcessor); err != nil {
			return err
		}
		ep.costAllocationProcessor = &costAllocationProcessor

	default:
		return fmt.Errorf("unknown CostAllocationProcessor type: %s", ttype)
	}

	return nil
}
