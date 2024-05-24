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

package anomalydetectors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Analyzer          *AnalyzerInput      `json:"analyzer"`          // Analyzer input
	Description       string              `json:"description"`       // The description of the anomaly detector
	Enabled           bool                `json:"enabled"`           // This setting is enabled (`true`) or disabled (`false`)
	EventTemplate     *DavisEventTemplate `json:"eventTemplate"`     // Event template
	ExecutionSettings *ExecutionSettings  `json:"executionSettings"` // Execution settings
	Source            string              `json:"source"`            // Source
	Title             string              `json:"title"`             // The title of the anomaly detector
}

func (me *Settings) Name() string {
	return me.Title
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"analyzer": {
			Type:        schema.TypeList,
			Description: "Analyzer input",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(AnalyzerInput).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the anomaly detector",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"event_template": {
			Type:        schema.TypeList,
			Description: "Event template",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DavisEventTemplate).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"execution_settings": {
			Type:        schema.TypeList,
			Description: "Execution settings",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ExecutionSettings).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"source": {
			Type:        schema.TypeString,
			Description: "Source",
			Required:    true,
		},
		"title": {
			Type:        schema.TypeString,
			Description: "The title of the anomaly detector",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"analyzer":           me.Analyzer,
		"description":        me.Description,
		"enabled":            me.Enabled,
		"event_template":     me.EventTemplate,
		"execution_settings": me.ExecutionSettings,
		"source":             me.Source,
		"title":              me.Title,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"analyzer":           &me.Analyzer,
		"description":        &me.Description,
		"enabled":            &me.Enabled,
		"event_template":     &me.EventTemplate,
		"execution_settings": &me.ExecutionSettings,
		"source":             &me.Source,
		"title":              &me.Title,
	})
}
