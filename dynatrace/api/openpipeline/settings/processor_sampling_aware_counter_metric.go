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

type SamplingAwareCounterMetricExtractionProcessor struct {
	Processor
	Aggregation string   `json:"aggregation,omitempty"`
	Dimensions  []string `json:"dimensions,omitempty"`
	MetricKey   string   `json:"metricKey"`
	Sampling    string   `json:"sampling,omitempty"`
}

func (p *SamplingAwareCounterMetricExtractionProcessor) Schema() map[string]*schema.Schema {
	s := p.Processor.Schema()

	s["aggregation"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Indicates whether aggregation is considered for metric extraction",
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

	s["metric_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The key of the metric to write",
		Required:    true,
	}

	s["sampling"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Indicates whether sampling is considered for metric extraction. Possible values: 'DISABLED', 'ENABLED'",
		Optional:    true,
	}

	return s
}

func (p *SamplingAwareCounterMetricExtractionProcessor) MarshalHCL(properties hcl.Properties) error {
	if err := p.Processor.MarshalHCL(properties); err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"aggregation": p.Aggregation,
		"dimensions":  p.Dimensions,
		"metric_key":  p.MetricKey,
		"sampling":    p.Sampling,
	})
}

func (p *SamplingAwareCounterMetricExtractionProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := p.Processor.UnmarshalHCL(decoder); err != nil {
		return err
	}

	return decoder.DecodeAll(map[string]any{
		"aggregation": &p.Aggregation,
		"dimensions":  &p.Dimensions,
		"metric_key":  &p.MetricKey,
		"sampling":    &p.Sampling,
	})
}

func (ep SamplingAwareCounterMetricExtractionProcessor) MarshalJSON() ([]byte, error) {
	type samplingAwareCounterMetricExtractionProcessor SamplingAwareCounterMetricExtractionProcessor
	return MarshalAsJSONWithType((samplingAwareCounterMetricExtractionProcessor)(ep), SamplingAwareCounterMetricExtractionProcessorType)
}
