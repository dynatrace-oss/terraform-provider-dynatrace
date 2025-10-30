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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SamplingAwareCounterMetricAttributes struct {
	Aggregation *Aggregation           `json:"aggregation,omitempty"` // Possible Values: `disabled`, `enabled`.
	Dimensions  FieldExtractionEntries `json:"dimensions,omitempty"`  // List of dimensions
	MetricKey   string                 `json:"metricKey"`             // Metric key
	Sampling    *Sampling              `json:"sampling,omitempty"`    // Possible Values: `disabled`, `enabled`.
}

func (me *SamplingAwareCounterMetricAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aggregation": {
			Type:        schema.TypeString,
			Description: "Possible Values: `disabled`, `enabled`.",
			Optional:    true, // nullable
		},
		"dimensions": {
			Type:        schema.TypeList,
			Description: "List of dimensions",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(FieldExtractionEntries).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"metric_key": {
			Type:        schema.TypeString,
			Description: "Metric key",
			Required:    true,
		},
		"sampling": {
			Type:        schema.TypeString,
			Description: "Possible Values: `disabled`, `enabled`.",
			Optional:    true, // nullable
		},
	}
}

func (me *SamplingAwareCounterMetricAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"aggregation": me.Aggregation,
		"dimensions":  me.Dimensions,
		"metric_key":  me.MetricKey,
		"sampling":    me.Sampling,
	})
}

func (me *SamplingAwareCounterMetricAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"aggregation": &me.Aggregation,
		"dimensions":  &me.Dimensions,
		"metric_key":  &me.MetricKey,
		"sampling":    &me.Sampling,
	})
}
