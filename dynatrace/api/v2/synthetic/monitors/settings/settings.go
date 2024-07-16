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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled               bool                   `json:"enabled"`                                          // If true, the monitor is enabled
	Name                  string                 `json:"name"`                                             // Name of the monitor
	Type                  MonitorType            `json:"type"`                                             // Type of the monitor, possible values: `MULTI_PROTOCOL`
	Description           *string                `json:"description,omitempty"`                            // Description of the monitor
	Tags                  TagsWithSourceInfo     `json:"tags,omitempty"`                                   // A set of tags assigned to the monitor.\n\nYou can specify only the value of the tag here and the CONTEXTLESS context and source 'USER' will be added automatically. But preferred option is usage of SyntheticTagWithSourceDto model.
	Steps                 Steps                  `json:"steps"`                                            // The steps of the monitor
	FrequencyMin          *int64                 `json:"frequencyMin,omitempty"`                           // Frequency of the monitor, in minutes
	Locations             []string               `json:"locations"`                                        // The locations to which the monitor is assigned
	OutageHandling        *OutageHandling        `json:"syntheticMonitorOutageHandlingSettings,omitempty"` // Outage handling configuration
	PerformanceThresholds *PerformanceThresholds `json:"performanceThresholds,omitempty"`                  // Performance thresholds configuration
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "If true, the monitor is enabled",
			Optional:    true,
			Default:     true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the monitor",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Type of the monitor, possible values: `MULTI_PROTOCOL`",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description of the monitor",
			Optional:    true,
		},
		"tags": {
			Type:        schema.TypeList,
			Description: "A set of tags assigned to the monitor.\n\nYou can specify only the value of the tag here and the CONTEXTLESS context and source 'USER' will be added automatically. But preferred option is usage of SyntheticTagWithSourceDto model.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TagsWithSourceInfo).Schema()},
		},
		"steps": {
			Type:        schema.TypeList,
			Description: "The steps of the monitor",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Steps).Schema()},
		},
		"frequency_min": {
			Type:        schema.TypeInt,
			Description: "Frequency of the monitor, in minutes",
			Optional:    true,
		},
		"locations": {
			Type:        schema.TypeSet,
			Description: "The locations to which the monitor is assigned",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"outage_handling": {
			Type:             schema.TypeList,
			Description:      "Outage handling configuration",
			Optional:         true,
			Elem:             &schema.Resource{Schema: new(OutageHandling).Schema()},
			MinItems:         1,
			MaxItems:         1,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return newValue == "0" },
		},
		"performance_thresholds": {
			Type:             schema.TypeList,
			Description:      "Performance thresholds configuration",
			Optional:         true,
			Elem:             &schema.Resource{Schema: new(PerformanceThresholds).Schema()},
			MinItems:         1,
			MaxItems:         1,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return newValue == "0" },
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":                me.Enabled,
		"name":                   me.Name,
		"type":                   me.Type,
		"description":            me.Description,
		"tags":                   me.Tags,
		"steps":                  me.Steps,
		"frequency_min":          me.FrequencyMin,
		"locations":              me.Locations,
		"outage_handling":        me.OutageHandling,
		"performance_thresholds": me.PerformanceThresholds,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":                &me.Enabled,
		"name":                   &me.Name,
		"type":                   &me.Type,
		"description":            &me.Description,
		"tags":                   &me.Tags,
		"steps":                  &me.Steps,
		"frequency_min":          &me.FrequencyMin,
		"locations":              &me.Locations,
		"outage_handling":        &me.OutageHandling,
		"performance_thresholds": &me.PerformanceThresholds,
	})
}
