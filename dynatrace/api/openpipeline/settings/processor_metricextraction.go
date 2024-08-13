package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetricExtractionProcessors struct {
	Processors []MetricExtractionProcessor
}

func (ep *MetricExtractionProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeSet,
			Description: "Data extraction processor to use.",
			Elem:        &schema.Resource{Schema: new(DataExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *MetricExtractionProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("processors", ep.Processors)
}

func (ep *MetricExtractionProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("processors", &ep.Processors)
}

type MetricExtractionProcessor struct {
	counterMetricExtractionProcessor *CounterMetricExtractionProcessor
	valueMetricExtractionProcessor   *ValueMetricExtractionProcessor
}

func (ep *MetricExtractionProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"counter_metric_exctration_processor": {
			Type:        schema.TypeList,
			Description: "Processor to write the occurrences as a metric.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(CounterMetricExtractionProcessor).Schema()},
			Optional:    true,
		},
		"value_metric_extraction_processor": {
			Type:        schema.TypeList,
			Description: "Processor to extract a value from a field as a metric.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ValueMetricExtractionProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *MetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"counter_metric_exctration_processor": ep.counterMetricExtractionProcessor,
		"value_metric_extraction_processor":   ep.valueMetricExtractionProcessor,
	})
}

func (ep *MetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"counter_metric_exctration_processor": ep.counterMetricExtractionProcessor,
		"value_metric_extraction_processor":   ep.valueMetricExtractionProcessor,
	})
}
