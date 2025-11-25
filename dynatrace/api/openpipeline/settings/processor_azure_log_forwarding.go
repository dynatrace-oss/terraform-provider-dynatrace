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

type AzureLogForwardingProcessor struct {
	Processor
	FieldExtraction   FieldExtraction `json:"fieldExtraction"`
	ForwarderConfigId string          `json:"forwarderConfigId,omitempty"`
}

func (ep *AzureLogForwardingProcessor) Schema() map[string]*schema.Schema {
	s := ep.Processor.Schema()

	s["field_extraction"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(FieldExtraction).Schema()},
		Description: "Definition of the field extraction",
		Optional:    true,
	}

	s["forwarder_config_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "",
		Required:    true,
	}

	return s
}

func (ep *AzureLogForwardingProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := ep.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"field_extraction":    ep.FieldExtraction,
		"forwarder_config_id": ep.ForwarderConfigId,
	})
}

func (ep *AzureLogForwardingProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := ep.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.DecodeAll(map[string]any{
		"field_extraction":    &ep.FieldExtraction,
		"forwarder_config_id": &ep.ForwarderConfigId,
	})
}

func (p AzureLogForwardingProcessor) MarshalJSON() ([]byte, error) {
	type azureLogForwardingProcessor AzureLogForwardingProcessor
	return MarshalAsJSONWithType((azureLogForwardingProcessor)(p), AzureLogForwardingProcessorType)
}
