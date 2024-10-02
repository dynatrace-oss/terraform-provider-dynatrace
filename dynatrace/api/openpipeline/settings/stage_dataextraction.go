package openpipeline

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataExtractionStage struct {
	Editable   *bool                      `json:"editable,omitempty"`
	Processors []*DataExtractionProcessor `json:"processors"`
}

func (f *DataExtractionStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(DataExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (f *DataExtractionStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", f.Processors)
}

func (f *DataExtractionStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &f.Processors)
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
		"davis_event_extraction_processor": &ep.davisEventExtractionProcessor,
		"bizevent_extraction_processor":    &ep.bizEventExtractionProcessor,
	})
}

func (ep DataExtractionProcessor) MarshalJSON() ([]byte, error) {
	if ep.bizEventExtractionProcessor != nil {
		return json.Marshal(ep.bizEventExtractionProcessor)
	}
	if ep.davisEventExtractionProcessor != nil {
		return json.Marshal(ep.davisEventExtractionProcessor)
	}

	return nil, errors.New("missing MetricExtractionProcessor value")
}

func (ep *DataExtractionProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case DavisEventExtractionProcessorType:
		davisEventExtractionProcessor := DavisEventExtractionProcessor{}
		if err := json.Unmarshal(b, &davisEventExtractionProcessor); err != nil {
			return err
		}
		ep.davisEventExtractionProcessor = &davisEventExtractionProcessor

	case BizEventExtractionProcessorType:
		bizEventExtractionProcessor := BizEventExtractionProcessor{}
		if err := json.Unmarshal(b, &bizEventExtractionProcessor); err != nil {
			return err
		}
		ep.bizEventExtractionProcessor = &bizEventExtractionProcessor

	default:
		return fmt.Errorf("unknown DataExtractionProcessor type: %s", ttype)
	}

	return nil
}
