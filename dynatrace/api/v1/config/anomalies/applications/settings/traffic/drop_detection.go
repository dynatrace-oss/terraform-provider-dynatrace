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

// DropDetection The configuration of traffic drops detection.
type DropDetection struct {
	Enabled            bool   `json:"enabled"`                      // The detection is enabled (`true`) or disabled (`false`).
	TrafficDropPercent *int32 `json:"trafficDropPercent,omitempty"` // Alert if the observed traffic is less than *X* % of the expected value.
}

func (me *DropDetection) Schema() map[string]*schema.Schema {
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

func (me *DropDetection) MarshalHCL(properties hcl.Properties) error {
	if !me.Enabled {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"percent": me.TrafficDropPercent,
	})
}

func (me *DropDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("percent"); ok {
		me.TrafficDropPercent = opt.NewInt32(int32(value.(int)))
	}

	return nil
}
