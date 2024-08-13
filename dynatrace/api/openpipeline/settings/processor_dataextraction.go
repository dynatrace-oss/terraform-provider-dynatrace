package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataExtractionProcessors struct {
	Processors []DataExtractionProcessor
}

func (ep *DataExtractionProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(DataExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *DataExtractionProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("processors", ep.Processors)
}

func (ep *DataExtractionProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("processors", &ep.Processors)
}

type DataExtractionProcessor struct {
	davisEventExtractionProcessor *DavisEventExtractionProcessor
	bizEventExtractionProcessor   *BizEventExtractionProcessor
}

func (ep *DataExtractionProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"davis_event_extraction_processor": {
			Type:        schema.TypeList,
			Description: "Processor to apply a DQL script",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DavisEventExtractionProcessor).Schema()},
			Optional:    true,
		},
		"bizevent_extraction_processor": {
			Type:        schema.TypeList,
			Description: "",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(BizEventExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *DataExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"davis_event_extraction_processor": ep.davisEventExtractionProcessor,
		"bizevent_extraction_processor":    ep.bizEventExtractionProcessor,
	})
}

func (ep *DataExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"davis_event_extraction_processor": ep.davisEventExtractionProcessor,
		"bizevent_extraction_processor":    ep.bizEventExtractionProcessor,
	})
}
