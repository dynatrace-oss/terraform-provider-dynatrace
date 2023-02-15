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

package schemalesslogmetric

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Dimensions       []string `json:"dimensions"`
	Enabled          bool     `json:"enabled"`                    // This setting is enabled (`true`) or disabled (`false`)
	Key              string   `json:"key"`                        // Metric key
	Measure          Measure  `json:"measure"`                    // Possible Values: `OCCURRENCE`, `ATTRIBUTE`
	MeasureAttribute *string  `json:"measureAttribute,omitempty"` // Attribute
	Query            string   `json:"query"`                      // Matcher
}

func (me *Settings) Name() string {
	return me.Key
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dimensions": {
			Type:        schema.TypeSet,
			Description: "no documentation available",
			Required:    true,

			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Metric key",
			Required:    true,
		},
		"measure": {
			Type:        schema.TypeString,
			Description: "Possible Values: `OCCURRENCE`, `ATTRIBUTE`",
			Required:    true,
		},
		"measure_attribute": {
			Type:        schema.TypeString,
			Description: "Attribute",
			Optional:    true,
		},
		"query": {
			Type:        schema.TypeString,
			Description: "Matcher",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dimensions":        me.Dimensions,
		"enabled":           me.Enabled,
		"key":               me.Key,
		"measure":           me.Measure,
		"measure_attribute": me.MeasureAttribute,
		"query":             me.Query,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dimensions":        &me.Dimensions,
		"enabled":           &me.Enabled,
		"key":               &me.Key,
		"measure":           &me.Measure,
		"measure_attribute": &me.MeasureAttribute,
		"query":             &me.Query,
	})
}
