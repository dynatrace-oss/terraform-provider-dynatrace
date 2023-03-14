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

type Settings struct {
	Description       *string            `json:"description,omitempty"`       // Description
	Dimensions        Dimensions         `json:"dimensions"`                  // Define metadata per metric dimension.
	DisplayName       *string            `json:"displayName,omitempty"`       // Display name
	MetricProperties  *MetricProperties  `json:"metricProperties,omitempty"`  // Metric properties
	MetricID          string             `json:"-" scope:"metricId"`          // The scope of this setting (metric)
	SourceEntityType  *string            `json:"sourceEntityType,omitempty"`  // Specifies which entity dimension should be used as the primary dimension. The property can only be configured for metrics ingested with the Metrics API.
	Tags              []string           `json:"tags"`                        // Tags
	Unit              string             `json:"unit"`                        // Unit
	UnitDisplayFormat *UnitDisplayFormat `json:"unitDisplayFormat,omitempty"` // The raw value is stored in bits or bytes. The user interface can display it in these numeral systems:\n\nBinary: 1 MiB = 1024 KiB = 1,048,576 bytes\n\nDecimal: 1 MB = 1000 kB = 1,000,000 bytes\n\nIf not set, the decimal system is used.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true,
		},
		"dimensions": {
			Type:        schema.TypeList,
			Description: "Define metadata per metric dimension.",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(Dimensions).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name",
			Optional:    true,
		},
		"metric_properties": {
			Type:        schema.TypeList,
			Description: "Metric properties",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(MetricProperties).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"metric_id": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (metric)",
			Required:    true,
		},
		"source_entity_type": {
			Type:        schema.TypeString,
			Description: "Specifies which entity dimension should be used as the primary dimension. The property can only be configured for metrics ingested with the Metrics API.",
			Optional:    true,
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "Tags",
			Optional:    true,

			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"unit": {
			Type:        schema.TypeString,
			Description: "Unit",
			Required:    true,
		},
		"unit_display_format": {
			Type:        schema.TypeString,
			Description: "The raw value is stored in bits or bytes. The user interface can display it in these numeral systems:\n\nBinary: 1 MiB = 1024 KiB = 1,048,576 bytes\n\nDecimal: 1 MB = 1000 kB = 1,000,000 bytes\n\nIf not set, the decimal system is used.",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":         me.Description,
		"dimensions":          me.Dimensions,
		"display_name":        me.DisplayName,
		"metric_properties":   me.MetricProperties,
		"metric_id":           me.MetricID,
		"source_entity_type":  me.SourceEntityType,
		"tags":                me.Tags,
		"unit":                me.Unit,
		"unit_display_format": me.UnitDisplayFormat,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description":         &me.Description,
		"dimensions":          &me.Dimensions,
		"display_name":        &me.DisplayName,
		"metric_properties":   &me.MetricProperties,
		"metric_id":           &me.MetricID,
		"source_entity_type":  &me.SourceEntityType,
		"tags":                &me.Tags,
		"unit":                &me.Unit,
		"unit_display_format": &me.UnitDisplayFormat,
	})
}
