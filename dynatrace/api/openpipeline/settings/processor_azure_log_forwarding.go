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
