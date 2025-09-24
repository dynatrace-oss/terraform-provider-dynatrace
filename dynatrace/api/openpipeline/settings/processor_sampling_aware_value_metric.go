package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SamplingAwareValueMetricExtractionProcessor struct {
	Processor
	Aggregation  string   `json:"aggregation,omitempty"`
	DefaultValue string   `json:"defaultValue,omitempty"`
	Dimensions   []string `json:"dimensions,omitempty"`
	Field        string   `json:"field,omitempty"`
	Measurement  string   `json:"measurement"`
	MetricKey    string   `json:"metricKey"`
	Sampling     string   `json:"sampling,omitempty"`
}

func (p *SamplingAwareValueMetricExtractionProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()

	s["aggregation"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Indicates whether aggregation is considered for metric extraction. Possible values: 'DISABLED', 'ENABLED'",
		Optional:    true,
	}

	s["default_value"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The default value to use if the given field value could not be extracted. String value must be a number",
		Optional:    true,
	}

	s["dimensions"] = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{Type: schema.TypeString},
		Description: "List of unique dimensions to add to the metric.\nDimensions are represented in the format '({dimensionName}=)?{sourceField}'.\n" +
			"'{dimensionName}' is optional if {sourceField} represents a valid metric dimension name.\n" +
			"'{sourceField}' has to represent a valid DQL field accessor and it can access a nested field (for example, 'field[field2][0]')",
		Optional: true,
	}

	s["field"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The field from which to extract the value for the metric, or 'null' for a duration measurement",
		Optional:    true,
	}

	s["metric_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The key of the metric to write",
		Required:    true,
	}

	s["measurement"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Indicates whether the metric value should be the specified field or duration. Possible values: 'DURATION', 'FIELD'",
		Required:    true,
	}

	s["sampling"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Indicates whether sampling is considered for metric extraction. Possible values: 'DISABLED', 'ENABLED'",
		Optional:    true,
	}

	return s
}

func (p *SamplingAwareValueMetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"aggregation":   p.Aggregation,
		"default_value": p.DefaultValue,
		"dimensions":    p.Dimensions,
		"field":         p.Field,
		"measurement":   p.Measurement,
		"metric_key":    p.MetricKey,
		"sampling":      p.Sampling,
	})
}

func (p *SamplingAwareValueMetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.DecodeAll(map[string]any{
		"aggregation":   &p.Aggregation,
		"default_value": &p.DefaultValue,
		"dimensions":    &p.Dimensions,
		"field":         &p.Field,
		"measurement":   &p.Measurement,
		"metric_key":    &p.MetricKey,
		"sampling":      &p.Sampling,
	})
}

func (ep SamplingAwareValueMetricExtractionProcessor) MarshalJSON() ([]byte, error) {
	type samplingAwareValueMetricExtractionProcessor SamplingAwareValueMetricExtractionProcessor
	return MarshalAsJSONWithType((samplingAwareValueMetricExtractionProcessor)(ep), SamplingAwareValueMetricExtractionProcessorType)
}
