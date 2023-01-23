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

package load

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SpikeDetection The configuration of load spikes detection.
type SpikeDetection struct {
	Enabled          bool                       `json:"enabled"`                                     // The detection is enabled (`true`) or disabled (`false`).
	LoadSpikePercent *int32                     `json:"loadSpikePercent,omitempty"`                  // Alert if the observed load is more than *X* % of the expected value.
	AbnormalMinutes  *int32                     `json:"minAbnormalStateDurationInMinutes,omitempty"` // Alert if the service stays in abnormal state for at least *X* minutes.
	Unknowns         map[string]json.RawMessage `json:"-"`
}

func (me *SpikeDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"percent": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Alert if the observed load is more than *X* % of the expected value",
		},
		"minutes": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Alert if the service stays in abnormal state for at least *X* minutes",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *SpikeDetection) MarshalHCL(properties hcl.Properties) error {
	if !me.Enabled {
		return nil
	}
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"percent": me.LoadSpikePercent,
		"minutes": me.AbnormalMinutes,
	})
}

func (me *SpikeDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "loadSpikePercent")
		delete(me.Unknowns, "minAbnormalStateDurationInMinutes")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	me.Enabled = true
	if value, ok := decoder.GetOk("percent"); ok {
		me.LoadSpikePercent = opt.NewInt32(int32(value.(int)))
	}
	if value, ok := decoder.GetOk("minutes"); ok {
		me.AbnormalMinutes = opt.NewInt32(int32(value.(int)))
	}
	return nil
}

func (me *SpikeDetection) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"enabled":                           me.Enabled,
		"loadSpikePercent":                  me.LoadSpikePercent,
		"minAbnormalStateDurationInMinutes": me.AbnormalMinutes,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *SpikeDetection) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"enabled":                           &me.Enabled,
		"loadSpikePercent":                  &me.LoadSpikePercent,
		"minAbnormalStateDurationInMinutes": &me.AbnormalMinutes,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
