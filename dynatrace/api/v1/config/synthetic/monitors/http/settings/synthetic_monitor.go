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

package http

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SyntheticMonitor HTTP synthetic monitor update. Some fields are inherited from base `SyntheticMonitorUpdate` model
type SyntheticMonitor struct {
	monitors.SyntheticMonitor
	Script *Script `json:"script,omitempty"`
}

func (me *SyntheticMonitor) GetType() monitors.Type {
	return monitors.Types.HTTP
}

func (me *SyntheticMonitor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		// ID                   *string                     `json:"entityId,omitempty"`
		Name                 string                      `json:"name"`
		Type                 monitors.Type               `json:"type"`
		FrequencyMin         int32                       `json:"frequencyMin"`
		Enabled              bool                        `json:"enabled"`
		AnomalyDetection     *monitors.AnomalyDetection  `json:"anomalyDetection,omitempty"`
		Locations            []string                    `json:"locations"`
		Tags                 monitors.TagsWithSourceInfo `json:"tags"`
		ManuallyAssignedApps []string                    `json:"manuallyAssignedApps"`
		Script               *Script                     `json:"script,omitempty"`
	}{
		// me.ID,
		me.Name,
		me.GetType(),
		me.FrequencyMin,
		me.Enabled,
		me.AnomalyDetection,
		me.Locations,
		me.GetTags(),
		me.ManuallyAssignedApps,
		me.Script,
	})
}

func (me *SyntheticMonitor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the monitor.",
			Required:    true,
		},
		"frequency": {
			Type:        schema.TypeInt,
			Description: "The frequency of the monitor, in minutes.\n\nYou can use one of the following values: `5`, `10`, `15`, `30`, and `60`.",
			Required:    true,
		},
		"locations": {
			Type:        schema.TypeSet,
			Description: "A list of locations from which the monitor is executed.\n\nTo specify a location, use its entity ID.",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The monitor is enabled (`true`) or disabled (`false`).",
			Optional:    true,
		},
		"manually_assigned_apps": {
			Type:        schema.TypeSet,
			Description: "A set of manually assigned applications.",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tags": {
			Type:        schema.TypeList,
			Description: "A set of tags assigned to the monitor.\n\nYou can specify only the value of the tag here and the `CONTEXTLESS` context and source 'USER' will be added automatically.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(monitors.TagsWithSourceInfo).Schema(),
			},
		},
		"anomaly_detection": {
			Type:        schema.TypeList,
			Description: "The anomaly detection configuration.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(monitors.AnomalyDetection).Schema(),
			},
		},
		"script": {
			Type:        schema.TypeList,
			Description: "The HTTP Script",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: new(Script).Schema(),
			},
		},
	}
}

func (me *SyntheticMonitor) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("frequency", me.FrequencyMin); err != nil {
		return err
	}
	if err := properties.Encode("locations", me.Locations); err != nil {
		return err
	}
	if err := properties.Encode("enabled", me.Enabled); err != nil {
		return err
	}
	if err := properties.Encode("manually_assigned_apps", me.ManuallyAssignedApps); err != nil {
		return err
	}
	if err := properties.Encode("tags", me.Tags); err != nil {
		return err
	}
	if err := properties.Encode("anomaly_detection", me.AnomalyDetection); err != nil {
		return err
	}
	if err := properties.Encode("script", me.Script); err != nil {
		return err
	}
	return nil
}

func (me *SyntheticMonitor) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("frequency"); ok {
		me.FrequencyMin = int32(value.(int))
	}
	if err := decoder.Decode("locations", &me.Locations); err != nil {
		return err
	}
	if err := decoder.Decode("manually_assigned_apps", &me.ManuallyAssignedApps); err != nil {
		return err
	}
	if me.ManuallyAssignedApps == nil {
		me.ManuallyAssignedApps = []string{}
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if _, ok := decoder.GetOk("tags.#"); ok {
		me.Tags = monitors.TagsWithSourceInfo{}
		if err := me.Tags.UnmarshalHCL(hcl.NewDecoder(decoder, "tags", 0)); err != nil {
			return err
		}
	}
	if err := decoder.Decode("tags", &me.Tags); err != nil {
		return err
	}
	if err := decoder.Decode("anomaly_detection", &me.AnomalyDetection); err != nil {
		return err
	}
	if err := decoder.Decode("script", &me.Script); err != nil {
		return err
	}
	return nil
}
