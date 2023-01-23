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

package connection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// LostDetectionConfig Configuration of lost connection detection.
type LostDetectionConfig struct {
	EnabledOnGracefulShutdowns bool `json:"enabledOnGracefulShutdowns"` // Alert (`true`) on graceful host shutdowns.
	Enabled                    bool `json:"enabled"`                    // The detection is enabled (`true`) or disabled (`false`).
}

func (me *LostDetectionConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The detection is enabled (`true`) or disabled (`false`)",
		},
		"enabled_on_graceful_shutdowns": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Alert (`true`) on graceful host shutdowns",
		},
	}
}

func (me *LostDetectionConfig) MarshalHCL(properties hcl.Properties) error {
	if !me.Enabled {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"enabled":                       me.Enabled,
		"enabled_on_graceful_shutdowns": me.EnabledOnGracefulShutdowns,
	})
}

func (me *LostDetectionConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("enabled_on_graceful_shutdowns"); ok {
		me.EnabledOnGracefulShutdowns = value.(bool)
	}
	return nil
}
