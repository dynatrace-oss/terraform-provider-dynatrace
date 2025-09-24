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
			Description: "Groups all processors applicable for the DataExtractionStage.\nApplicable processors are DavisEventExtractionProcessor and BizeventExtractionProcessor.",
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
	davisEventExtractionProcessor    *DavisEventExtractionProcessor
	bizEventExtractionProcessor      *BizEventExtractionProcessor
	azureLogForwardingProcessor      *AzureLogForwardingProcessor
	securityEventExtractionProcessor *SecurityEventExtractionProcessor
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
			Description: "Processor to extract a bizevent.\nFields event.type and event.provider can only be assigned to a constant or field value.\nA multi-value constant is not supported for those fields.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(BizEventExtractionProcessor).Schema()},
			Optional:    true,
		},
		"azure_log_forwarding_processor": {
			Type:        schema.TypeList,
			Description: "Processor to extract a Azure log.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(AzureLogForwardingProcessor).Schema()},
			Optional:    true,
		},
		"security_event_extraction_processor": {
			Type:        schema.TypeList,
			Description: "Processor to extract a security event.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SecurityEventExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *DataExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"davis_event_extraction_processor":    ep.davisEventExtractionProcessor,
		"bizevent_extraction_processor":       ep.bizEventExtractionProcessor,
		"azure_log_forwarding_processor":      ep.azureLogForwardingProcessor,
		"security_event_extraction_processor": ep.securityEventExtractionProcessor,
	})
}

func (ep *DataExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"davis_event_extraction_processor":    &ep.davisEventExtractionProcessor,
		"bizevent_extraction_processor":       &ep.bizEventExtractionProcessor,
		"azure_log_forwarding_processor":      &ep.azureLogForwardingProcessor,
		"security_event_extraction_processor": &ep.securityEventExtractionProcessor,
	})
}

func (ep DataExtractionProcessor) MarshalJSON() ([]byte, error) {
	if ep.bizEventExtractionProcessor != nil {
		return json.Marshal(ep.bizEventExtractionProcessor)
	}
	if ep.davisEventExtractionProcessor != nil {
		return json.Marshal(ep.davisEventExtractionProcessor)
	}
	if ep.azureLogForwardingProcessor != nil {
		return json.Marshal(ep.azureLogForwardingProcessor)
	}
	if ep.securityEventExtractionProcessor != nil {
		return json.Marshal(ep.securityEventExtractionProcessor)
	}

	return nil, errors.New("missing DataExtractionProcessor value")
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

	case AzureLogForwardingProcessorType:
		azureLogForwardingProcessor := AzureLogForwardingProcessor{}
		if err := json.Unmarshal(b, &azureLogForwardingProcessor); err != nil {
			return err
		}
		ep.azureLogForwardingProcessor = &azureLogForwardingProcessor

	case SecurityEventExtractionProcessorType:
		securityEventExtractionProcessor := SecurityEventExtractionProcessor{}
		if err := json.Unmarshal(b, &securityEventExtractionProcessor); err != nil {
			return err
		}
		ep.securityEventExtractionProcessor = &securityEventExtractionProcessor

	default:
		return fmt.Errorf("unknown DataExtractionProcessor type: %s", ttype)
	}

	return nil
}
