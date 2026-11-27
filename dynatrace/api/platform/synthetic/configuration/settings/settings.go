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

package configuration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Settings holds the tenant-wide synthetic configuration. Unset fields are nil and keep the value from the tenant.
type Settings struct {
	BmMonitorTimeout                           *int64 `json:"bmMonitorTimeout,omitempty"`                           // Browser monitor execution timeout (ms)
	BmStepTimeout                              *int64 `json:"bmStepTimeout,omitempty"`                              // Browser monitor single step execution timeout (ms)
	MaintenanceWindowsExcludedFromAvailability *bool  `json:"maintenanceWindowsExcludedFromAvailability,omitempty"` // Whether periods with maintenance windows are excluded from the availability calculation
}

func (me *Settings) Name() string {
	return "synthetic_v2_configuration"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bm_monitor_timeout": {
			Type:         schema.TypeInt,
			Description:  "Browser monitor execution timeout (ms). Allowed range is 60000 to 900000. If not set, the value of the environment is kept.",
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(60000, 900000),
		},
		"bm_step_timeout": {
			Type:         schema.TypeInt,
			Description:  "Browser monitor single step execution timeout (ms). Allowed range is 10000 to 180000. If not set, the value of the environment is kept.",
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(10000, 180000),
		},
		"maintenance_windows_excluded_from_availability": {
			Type:        schema.TypeBool,
			Description: "Whether periods with maintenance windows are excluded from the availability calculation. If not set, the value of the environment is kept.",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"bm_monitor_timeout": me.BmMonitorTimeout,
		"bm_step_timeout":    me.BmStepTimeout,
		"maintenance_windows_excluded_from_availability": me.MaintenanceWindowsExcludedFromAvailability,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"bm_monitor_timeout": &me.BmMonitorTimeout,
		"bm_step_timeout":    &me.BmStepTimeout,
	}); err != nil {
		return err
	}
	// GetOk reports an explicit false as unset, so the flag is decoded via GetOkExists.
	if value, ok := decoder.GetOkExists("maintenance_windows_excluded_from_availability"); ok {
		flag := value.(bool)
		me.MaintenanceWindowsExcludedFromAvailability = &flag
	}
	return nil
}
