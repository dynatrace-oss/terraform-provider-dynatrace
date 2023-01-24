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

// DropDetection The configuration of load drops detection.
type DropDetection struct {
	Enabled         bool   `json:"enabled"`                                     // The detection is enabled (`true`) or disabled (`false`).
	LoadDropPercent *int32 `json:"loadDropPercent,omitempty"`                   // Alert if the observed load is less than *X* % of the expected value.
	AbnormalMinutes *int32 `json:"minAbnormalStateDurationInMinutes,omitempty"` // Alert if the service stays in abnormal state for at least *X* minutes.
}

func (me *DropDetection) Schema() map[string]*schema.Schema {
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
	}
}

func (me *DropDetection) MarshalHCL(properties hcl.Properties) error {
	if !me.Enabled {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"percent": me.LoadDropPercent,
		"minutes": me.AbnormalMinutes,
	})
}

func (me *DropDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Enabled = true
	if value, ok := decoder.GetOk("percent"); ok {
		me.LoadDropPercent = opt.NewInt32(int32(value.(int)))
	}
	if value, ok := decoder.GetOk("minutes"); ok {
		me.AbnormalMinutes = opt.NewInt32(int32(value.(int)))
	}
	return nil
}

func (me *DropDetection) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"enabled":                           me.Enabled,
		"loadDropPercent":                   me.LoadDropPercent,
		"minAbnormalStateDurationInMinutes": me.AbnormalMinutes,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *DropDetection) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"enabled":                           &me.Enabled,
		"loadDropPercent":                   &me.LoadDropPercent,
		"minAbnormalStateDurationInMinutes": &me.AbnormalMinutes,
	}); err != nil {
		return err
	}
	return nil
}
