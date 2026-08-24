/**
* @license
* Copyright 2026 Dynatrace LLC
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

package maintenancewindows

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Synthetic monitor disablement. Configuration for pausing synthetic monitors during the maintenance window.
type SyntheticMonitorSettings struct {
	DisableSyntheticMonitorFilter *string `json:"disableSyntheticMonitorFilter,omitempty"` // DQL filter selecting which synthetic monitors to pause. Required when synthetic monitors are disabled.
	DisableSyntheticMonitors      bool    `json:"disableSyntheticMonitors"`                // When enabled, synthetic monitors matching the filter are paused during the maintenance window.
}

func (me *SyntheticMonitorSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disable_synthetic_monitor_filter": {
			Type:        schema.TypeString,
			Description: "DQL filter selecting which synthetic monitors to pause. Required when synthetic monitors are disabled.",
			Optional:    true, // nullable & precondition
		},
		"disable_synthetic_monitors": {
			Type:        schema.TypeBool,
			Description: "When enabled, synthetic monitors matching the filter are paused during the maintenance window.",
			Required:    true,
		},
	}
}

func (me *SyntheticMonitorSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"disable_synthetic_monitor_filter": me.DisableSyntheticMonitorFilter,
		"disable_synthetic_monitors":       me.DisableSyntheticMonitors,
	})
}

func (me *SyntheticMonitorSettings) HandlePreconditions() error {
	if (me.DisableSyntheticMonitorFilter != nil) && (!me.DisableSyntheticMonitors) {
		return fmt.Errorf("'disable_synthetic_monitor_filter' must not be specified unless 'disable_synthetic_monitors' is set to 'true'; got 'disable_synthetic_monitors'='%v'", me.DisableSyntheticMonitors)
	}
	return nil
}

func (me *SyntheticMonitorSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"disable_synthetic_monitor_filter": &me.DisableSyntheticMonitorFilter,
		"disable_synthetic_monitors":       &me.DisableSyntheticMonitors,
	})
}
