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

package traffic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SpikeDetection The configuration of traffic spikes detection.
type SpikeDetection struct {
	Enabled             bool   `json:"enabled"`                       // The detection is enabled (`true`) or disabled (`false`).
	TrafficSpikePercent *int32 `json:"trafficSpikePercent,omitempty"` // Alert if the observed traffic is more than *X* % of the expected value.
}

func (me *SpikeDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The detection is enabled (`true`) or disabled (`false`)",
		},
		"percent": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Alert if the observed traffic is less than *X* % of the expected value",
		},
	}
}

func (me *SpikeDetection) MarshalHCL(properties hcl.Properties) error {
	if !me.Enabled {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"percent": me.TrafficSpikePercent,
	})
}

func (me *SpikeDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("percent"); ok {
		me.TrafficSpikePercent = opt.NewInt32(int32(value.(int)))
	}

	return nil
}
