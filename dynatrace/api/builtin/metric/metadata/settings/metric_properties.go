/**
* @license
* Copyright 2020 Dynatrace LLC
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

package metadata

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetricProperties struct {
	ImpactRelevant    *bool     `json:"impactRelevant,omitempty"`    // Whether (true or false) the metric is relevant to a problem's impact.\n\nAn impact-relevant metric is highly dependent on other metrics and changes because an underlying root-cause metric has changed.
	Latency           *int      `json:"latency,omitempty"`           // The latency of the metric, in minutes. \n\n The latency is the expected reporting delay (for example, caused by constraints of cloud vendors or other third-party data sources) between the observation of a metric data point and its availability in Dynatrace. \n\nThe allowed value range is from 1 to 60 minutes.
	MaxValue          *float64  `json:"maxValue,omitempty"`          // The maximum allowed value of the metric.
	MinValue          *float64  `json:"minValue,omitempty"`          // The minimum allowed value of the metric.
	RootCauseRelevant *bool     `json:"rootCauseRelevant,omitempty"` // Whether (true or false) the metric is related to a root cause of a problem.\n\nA root-cause relevant metric represents a strong indicator for a faulty component.
	ValueType         ValueType `json:"valueType"`                   // Possible Values: `Error`, `Score`, `Unknown`
}

func (me *MetricProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"impact_relevant": {
			Type:        schema.TypeBool,
			Description: "Whether (true or false) the metric is relevant to a problem's impact.\n\nAn impact-relevant metric is highly dependent on other metrics and changes because an underlying root-cause metric has changed.",
			Optional:    true,
		},
		"latency": {
			Type:        schema.TypeInt,
			Description: "The latency of the metric, in minutes. \n\n The latency is the expected reporting delay (for example, caused by constraints of cloud vendors or other third-party data sources) between the observation of a metric data point and its availability in Dynatrace. \n\nThe allowed value range is from 1 to 60 minutes.",
			Optional:    true,
		},
		"max_value": {
			Type:        schema.TypeFloat,
			Description: "The maximum allowed value of the metric.",
			Optional:    true,
		},
		"min_value": {
			Type:        schema.TypeFloat,
			Description: "The minimum allowed value of the metric.",
			Optional:    true,
		},
		"root_cause_relevant": {
			Type:        schema.TypeBool,
			Description: "Whether (true or false) the metric is related to a root cause of a problem.\n\nA root-cause relevant metric represents a strong indicator for a faulty component.",
			Optional:    true,
		},
		"value_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Error`, `Score`, `Unknown`",
			Required:    true,
		},
	}
}

func (me *MetricProperties) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"impact_relevant":     me.ImpactRelevant,
		"latency":             me.Latency,
		"max_value":           me.MaxValue,
		"min_value":           me.MinValue,
		"root_cause_relevant": me.RootCauseRelevant,
		"value_type":          me.ValueType,
	})
}

func (me *MetricProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"impact_relevant":     &me.ImpactRelevant,
		"latency":             &me.Latency,
		"max_value":           &me.MaxValue,
		"min_value":           &me.MinValue,
		"root_cause_relevant": &me.RootCauseRelevant,
		"value_type":          &me.ValueType,
	})
}
