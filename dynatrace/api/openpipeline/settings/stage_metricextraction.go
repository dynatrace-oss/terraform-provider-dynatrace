package openpipeline

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetricExtractionStage struct {
	Editable   *bool                        `json:"editable,omitempty"`
	Processors []*MetricExtractionProcessor `json:"processors"`
}

func (f *MetricExtractionStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Data extraction processor to use",
			Elem:        &schema.Resource{Schema: new(MetricExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (f *MetricExtractionStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", f.Processors)
}

func (f *MetricExtractionStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &f.Processors)
}

type MetricExtractionProcessor struct {
	counterMetricExtractionProcessor *CounterMetricExtractionProcessor
	valueMetricExtractionProcessor   *ValueMetricExtractionProcessor
}

func (ep *MetricExtractionProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"counter_metric_extraction_processor": {
			Type:        schema.TypeList,
			Description: "Processor to write the occurrences as a metric",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(CounterMetricExtractionProcessor).Schema()},
			Optional:    true,
		},
		"value_metric_extraction_processor": {
			Type:        schema.TypeList,
			Description: "Processor to extract a value from a field as a metric",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ValueMetricExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *MetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"counter_metric_extraction_processor": ep.counterMetricExtractionProcessor,
		"value_metric_extraction_processor":   ep.valueMetricExtractionProcessor,
	})
}

func (ep *MetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"counter_metric_extraction_processor": &ep.counterMetricExtractionProcessor,
		"value_metric_extraction_processor":   &ep.valueMetricExtractionProcessor,
	})
}

func (ep MetricExtractionProcessor) MarshalJSON() ([]byte, error) {
	if ep.counterMetricExtractionProcessor != nil {
		return json.Marshal(ep.counterMetricExtractionProcessor)
	}
	if ep.valueMetricExtractionProcessor != nil {
		return json.Marshal(ep.valueMetricExtractionProcessor)
	}

	return nil, errors.New("missing MetricExtractionProcessor value")
}

func (ep *MetricExtractionProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case CounterMetricProcessorType:
		counterMetricExtractionProcessor := CounterMetricExtractionProcessor{}
		if err := json.Unmarshal(b, &counterMetricExtractionProcessor); err != nil {
			return err
		}
		ep.counterMetricExtractionProcessor = &counterMetricExtractionProcessor

	case ValueMetricProcessorType:
		valueMetricExtractionProcessor := ValueMetricExtractionProcessor{}
		if err := json.Unmarshal(b, &valueMetricExtractionProcessor); err != nil {
			return err
		}
		ep.valueMetricExtractionProcessor = &valueMetricExtractionProcessor

	default:
		return fmt.Errorf("unknown MetricExtractionProcessor type: %s", ttype)
	}

	return nil
}
