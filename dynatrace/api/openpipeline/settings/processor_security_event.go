package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecurityEventExtractionProcessor struct {
	Processor
	FieldExtraction FieldExtraction `json:"fieldExtraction"`
}

func (ep *SecurityEventExtractionProcessor) Schema() map[string]*schema.Schema {
	s := ep.Processor.Schema()

	s["field_extraction"] = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Elem:        &schema.Resource{Schema: new(FieldExtraction).Schema()},
		Description: "Definition of the field extraction",
		Optional:    true,
	}

	return s
}

func (ep *SecurityEventExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := ep.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"field_extraction": ep.FieldExtraction,
	})
}

func (ep *SecurityEventExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := ep.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.DecodeAll(map[string]any{
		"field_extraction": &ep.FieldExtraction,
	})
}

func (p SecurityEventExtractionProcessor) MarshalJSON() ([]byte, error) {
	type securityEventExtractionProcessor SecurityEventExtractionProcessor
	return MarshalAsJSONWithType((securityEventExtractionProcessor)(p), SecurityEventExtractionProcessorType)
}
