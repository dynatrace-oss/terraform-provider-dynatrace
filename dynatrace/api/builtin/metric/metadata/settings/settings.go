/**
* @license
* Copyright 2026 Dynatrace LLC
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
	Dimensions        Dimensions         `json:"dimensions,omitempty"`        // Define metadata per metric dimension.
	DisplayName       *string            `json:"displayName,omitempty"`       // Display name
	MetricProperties  *MetricProperties  `json:"metricProperties,omitempty"`  // Metric properties
	MetricID          string             `json:"-" scope:"metricId"`          // The scope of this setting (metric)
	SourceEntityType  *string            `json:"sourceEntityType,omitempty"`  // Specifies which entity dimension should be used as the primary dimension. The property can only be configured for metrics ingested with the Metrics API.
	Tags              []string           `json:"tags,omitempty"`              // Tags
	Unit              string             `json:"unit"`                        // Unit
	UnitDisplayFormat *UnitDisplayFormat `json:"unitDisplayFormat,omitempty"` // The raw value is stored in bits or bytes. The user interface can display it in these numeral systems:\n\n  Binary: 1 MiB = 1024 KiB = 1,048,576 bytes\n\n  Decimal: 1 MB = 1000 kB = 1,000,000 bytes\n\n  If not set, the decimal system is used. Possible values: `binary`, `decimal`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:             schema.TypeString,
			Description:      "Description",
			Optional:         true, // nullable
			DiffSuppressFunc: hcl.SuppressEOT,
		},
		"dimensions": {
			Type:        schema.TypeList,
			Description: "Define metadata per metric dimension.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Dimensions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name",
			Optional:    true, // nullable
		},
		"metric_properties": {
			Type:        schema.TypeList,
			Description: "Metric properties",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(MetricProperties).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"metric_id": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (metric)",
			ForceNew:    true,
			Required:    true,
		},
		"source_entity_type": {
			Type:        schema.TypeString,
			Description: "Specifies which entity dimension should be used as the primary dimension. The property can only be configured for metrics ingested with the Metrics API.",
			Optional:    true, // nullable
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "Tags",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"unit": {
			Type:        schema.TypeString,
			Description: "Unit",
			Required:    true,
		},
		"unit_display_format": {
			Type:        schema.TypeString,
			Description: "The raw value is stored in bits or bytes. The user interface can display it in these numeral systems:\n\n  Binary: 1 MiB = 1024 KiB = 1,048,576 bytes\n\n  Decimal: 1 MB = 1000 kB = 1,000,000 bytes\n\n  If not set, the decimal system is used. Possible values: `binary`, `decimal`",
			Optional:    true, // nullable & precondition
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

func (me *Settings) HandlePreconditions() error {
	// ---- UnitDisplayFormat *UnitDisplayFormat -> {"property":"unit","type":"REGEX_MATCH","pattern":"^[bB][iI][tT]$|^[bB][iI][tT][pP][eE][rR][sS][eE][cC][oO][nN][dD]$|^[bB][iI][tT][pP][eE][rR][mM][iI][nN][uU][tT][eE]$|^[bB][iI][tT][pP][eE][rR][hH][oO][uU][rR]$|^[bB][yY][tT][eE]$|^[bB][yY][tT][eE][pP][eE][rR][sS][eE][cC][oO][nN][dD]$|^[bB][yY][tT][eE][pP][eE][rR][mM][iI][nN][uU][tT][eE]$|^[bB][yY][tT][eE][pP][eE][rR][hH][oO][uU][rR]$"}
	return nil
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
