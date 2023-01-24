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

package retransmission

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DetectionConfig Configuration of high retransmission rate detection.
type DetectionConfig struct {
	CustomThresholds *Thresholds `json:"customThresholds,omitempty"` // Custom thresholds for high retransmission rate. If not set, automatic mode is used.   **All** of these conditions must be met to trigger an alert.
	Enabled          bool        `json:"enabled"`                    // The detection is enabled (`true`) or disabled (`false`).
}

func (me *DetectionConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The detection is enabled (`true`) or disabled (`false`)",
		},
		"thresholds": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Custom thresholds for high retransmission rate. If not set, automatic mode is used.   **All** of these conditions must be met to trigger an alert",
			Elem:        &schema.Resource{Schema: new(Thresholds).Schema()},
		},
	}
}

func (me *DetectionConfig) MarshalHCL(properties hcl.Properties) error {
	if !me.Enabled {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"enabled":    me.Enabled,
		"thresholds": me.CustomThresholds,
	})
}

func (me *DetectionConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if _, ok := decoder.GetOk("thresholds.#"); ok {
		me.CustomThresholds = new(Thresholds)
		if err := me.CustomThresholds.UnmarshalHCL(hcl.NewDecoder(decoder, "thresholds", 0)); err != nil {
			return err
		}
	}
	return nil
}
