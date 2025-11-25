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

type ProductAllocationStage struct {
	Editable   *bool                              `json:"editable,omitempty"`
	Processors []*ProductAllocationStageProcessor `json:"processors"`
}

func (f *ProductAllocationStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Product allocation processor to use",
			Elem:        &schema.Resource{Schema: new(ProductAllocationStageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (f *ProductAllocationStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", f.Processors)
}

func (f *ProductAllocationStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &f.Processors)
}

type ProductAllocationStageProcessor struct {
	productAllocationProcessor *ProductAllocationProcessor
}

func (ep *ProductAllocationStageProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"product_allocation_processor": {
			Type:        schema.TypeList,
			Description: "Processor to write the occurrences as a product allocation",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ProductAllocationProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *ProductAllocationStageProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"product_allocation_processor": ep.productAllocationProcessor,
	})
}

func (ep *ProductAllocationStageProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"product_allocation_processor": &ep.productAllocationProcessor,
	})
}

func (ep ProductAllocationStageProcessor) MarshalJSON() ([]byte, error) {
	if ep.productAllocationProcessor != nil {
		return json.Marshal(ep.productAllocationProcessor)
	}

	return nil, errors.New("missing ProductAllocationProcessor value")
}

func (ep *ProductAllocationStageProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case ProductAllocationProcessorType:
		productAllocationProcessor := ProductAllocationProcessor{}
		if err := json.Unmarshal(b, &productAllocationProcessor); err != nil {
			return err
		}
		ep.productAllocationProcessor = &productAllocationProcessor

	default:
		return fmt.Errorf("unknown ProductAllocationProcessor type: %s", ttype)
	}

	return nil
}
