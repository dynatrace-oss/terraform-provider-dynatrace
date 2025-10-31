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

type ValueMetricAttributes struct {
	DefaultValue *string                `json:"defaultValue,omitempty"` // Default value with metric value
	Dimensions   FieldExtractionEntries `json:"dimensions,omitempty"`   // List of dimensions
	Field        string                 `json:"field"`                  // Field with metric value
	MetricKey    string                 `json:"metricKey"`              // Metric key
}

func (me *ValueMetricAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_value": {
			Type:        schema.TypeString,
			Description: "Default value with metric value",
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
		"field": {
			Type:        schema.TypeString,
			Description: "Field with metric value",
			Required:    true,
		},
		"metric_key": {
			Type:        schema.TypeString,
			Description: "Metric key",
			Required:    true,
		},
	}
}

func (me *ValueMetricAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"default_value": me.DefaultValue,
		"dimensions":    me.Dimensions,
		"field":         me.Field,
		"metric_key":    me.MetricKey,
	})
}

func (me *ValueMetricAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"default_value": &me.DefaultValue,
		"dimensions":    &me.Dimensions,
		"field":         &me.Field,
		"metric_key":    &me.MetricKey,
	})
}
