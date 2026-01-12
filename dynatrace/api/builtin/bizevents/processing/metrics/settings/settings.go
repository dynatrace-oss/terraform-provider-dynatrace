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

package metrics

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Dimensions       []string `json:"dimensions,omitempty"`
	Enabled          bool     `json:"enabled"`                    // This setting is enabled (`true`) or disabled (`false`)
	Key              string   `json:"key"`                        // Key
	Matcher          string   `json:"matcher"`                    // [See our documentation](https://dt-url.net/bp234rv)
	Measure          Measure  `json:"measure"`                    // Possible Values: `ATTRIBUTE`, `OCCURRENCE`
	MeasureAttribute *string  `json:"measureAttribute,omitempty"` // Attribute
}

func (me *Settings) Deprecated() string {
	return "Classic bizevents processing rules have been deprecated in favor of OpenPipeline. Please migrate your OpenPipeline configurations and use `dynatrace_openpipeline_v2_*` instead."
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dimensions": {
			Type:        schema.TypeSet,
			Description: "no documentation available",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Key",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "[See our documentation](https://dt-url.net/bp234rv)",
			Required:    true,
		},
		"measure": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ATTRIBUTE`, `OCCURRENCE`",
			Required:    true,
		},
		"measure_attribute": {
			Type:        schema.TypeString,
			Description: "Attribute",
			Optional:    true, // precondition
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dimensions":        me.Dimensions,
		"enabled":           me.Enabled,
		"key":               me.Key,
		"matcher":           me.Matcher,
		"measure":           me.Measure,
		"measure_attribute": me.MeasureAttribute,
	})
}

func (me *Settings) HandlePreconditions() error {
	if me.MeasureAttribute == nil && string(me.Measure) == "ATTRIBUTE" {
		return fmt.Errorf("'measure_attribute' must be specified if 'measure' is set to '%v'", me.Measure)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dimensions":        &me.Dimensions,
		"enabled":           &me.Enabled,
		"key":               &me.Key,
		"matcher":           &me.Matcher,
		"measure":           &me.Measure,
		"measure_attribute": &me.MeasureAttribute,
	})
}
