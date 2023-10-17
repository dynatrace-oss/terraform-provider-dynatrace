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

package osservicesmonitoring

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AlertActivationDuration    *int                       `json:"alertActivationDuration,omitempty"`    // The number of **10-second measurement cycles** before alerting is triggered
	Alerting                   bool                       `json:"alerting"`                             // Toggle the switch in order to enable or disable alerting for this policy
	DetectionConditionsLinux   LinuxDetectionConditions   `json:"detectionConditionsLinux,omitempty"`   // Detection rules
	DetectionConditionsWindows WindowsDetectionConditions `json:"detectionConditionsWindows,omitempty"` // Detection rules
	Enabled                    bool                       `json:"enabled"`                              // This setting is enabled (`true`) or disabled (`false`)
	Metadata                   MetadataItems              `json:"metadata,omitempty"`                   // Set of additional key-value properties to be attached to the triggered event.
	Monitoring                 bool                       `json:"monitoring"`                           // Toggle the switch in order to enable or disable availability metric monitoring for this policy. Availability metrics produce custom metrics. Refer to [documentation](https://dt-url.net/vl03xzk) for consumption examples. Each monitored service consumes one custom metric.
	Name                       string                     `json:"name"`                                 // Rule name
	NotInstalledAlerting       *bool                      `json:"notInstalledAlerting,omitempty"`       // By default, Dynatrace does not alert if the service is not installed. Toggle the switch to enable or disable this feature
	Scope                      *string                    `json:"-" scope:"scope"`                      // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	StatusConditionLinux       *string                    `json:"statusConditionLinux,omitempty"`       // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(failed)` – Matches services that are in failed state.\n\nAvailable logic operations:\n- `$not($eq(active))` – Matches services with state different from active.\n- `$or($eq(inactive),$eq(failed))` – Matches services that are either in inactive or failed state.\n\nUse one of the following values as a parameter for this condition:\n\n- `reloading`\n- `activating`\n- `deactivating`\n- `failed`\n- `inactive`\n- `active`
	StatusConditionWindows     *string                    `json:"statusConditionWindows,omitempty"`     // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(paused)` – Matches services that are in paused state.\n\nAvailable logic operations:\n- `$not($eq(paused))` – Matches services that are in state different from paused.\n- `$or($eq(paused),$eq(running))` – Matches services that are either in paused or running state.\n\nUse one of the following values as a parameter for this condition:\n\n- `running`\n- `stopped`\n- `start_pending`\n- `stop_pending`\n- `continue_pending`\n- `pause_pending`\n- `paused`
	System                     System                     `json:"system"`                               // Possible Values: `LINUX`, `WINDOWS`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alert_activation_duration": {
			Type:        schema.TypeInt,
			Description: "The number of **10-second measurement cycles** before alerting is triggered",
			Optional:    true, // precondition
		},
		"alerting": {
			Type:        schema.TypeBool,
			Description: "Toggle the switch in order to enable or disable alerting for this policy",
			Required:    true,
		},
		"detection_conditions_linux": {
			Type:        schema.TypeList,
			Description: "Detection rules",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(LinuxDetectionConditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"detection_conditions_windows": {
			Type:        schema.TypeList,
			Description: "Detection rules",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(WindowsDetectionConditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Set of additional key-value properties to be attached to the triggered event.",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(MetadataItems).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"monitoring": {
			Type:        schema.TypeBool,
			Description: "Toggle the switch in order to enable or disable availability metric monitoring for this policy. Availability metrics produce custom metrics. Refer to [documentation](https://dt-url.net/vl03xzk) for consumption examples. Each monitored service consumes one custom metric.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
		"not_installed_alerting": {
			Type:        schema.TypeBool,
			Description: "By default, Dynatrace does not alert if the service is not installed. Toggle the switch to enable or disable this feature",
			Optional:    true, // precondition
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"status_condition_linux": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(failed)` – Matches services that are in failed state.\n\nAvailable logic operations:\n- `$not($eq(active))` – Matches services with state different from active.\n- `$or($eq(inactive),$eq(failed))` – Matches services that are either in inactive or failed state.\n\nUse one of the following values as a parameter for this condition:\n\n- `reloading`\n- `activating`\n- `deactivating`\n- `failed`\n- `inactive`\n- `active`",
			Optional:    true, // precondition
		},
		"status_condition_windows": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(paused)` – Matches services that are in paused state.\n\nAvailable logic operations:\n- `$not($eq(paused))` – Matches services that are in state different from paused.\n- `$or($eq(paused),$eq(running))` – Matches services that are either in paused or running state.\n\nUse one of the following values as a parameter for this condition:\n\n- `running`\n- `stopped`\n- `start_pending`\n- `stop_pending`\n- `continue_pending`\n- `pause_pending`\n- `paused`",
			Optional:    true, // precondition
		},
		"system": {
			Type:        schema.TypeString,
			Description: "Possible Values: `LINUX`, `WINDOWS`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"alert_activation_duration":    me.AlertActivationDuration,
		"alerting":                     me.Alerting,
		"detection_conditions_linux":   me.DetectionConditionsLinux,
		"detection_conditions_windows": me.DetectionConditionsWindows,
		"enabled":                      me.Enabled,
		"metadata":                     me.Metadata,
		"monitoring":                   me.Monitoring,
		"name":                         me.Name,
		"not_installed_alerting":       me.NotInstalledAlerting,
		"scope":                        me.Scope,
		"status_condition_linux":       me.StatusConditionLinux,
		"status_condition_windows":     me.StatusConditionWindows,
		"system":                       me.System,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.NotInstalledAlerting == nil) && (me.Alerting) {
		me.NotInstalledAlerting = opt.NewBool(false)
	}
	if (me.AlertActivationDuration == nil) && (me.Alerting) {
		return fmt.Errorf("'alert_activation_duration' must be specified if 'alerting' is set to '%v'", me.Alerting)
	}
	if (me.StatusConditionLinux == nil) && ((string(me.System) == "LINUX") && (me.Alerting)) {
		return fmt.Errorf("'status_condition_linux' must be specified if 'system' is set to '%v'", me.System)
	}
	if (me.StatusConditionWindows == nil) && ((string(me.System) == "WINDOWS") && (me.Alerting)) {
		return fmt.Errorf("'status_condition_windows' must be specified if 'system' is set to '%v'", me.System)
	}
	// ---- DetectionConditionsLinux LinuxDetectionConditions -> {"expectedValues":["LINUX"],"property":"system","type":"IN"}
	// ---- DetectionConditionsWindows WindowsDetectionConditions -> {"expectedValues":["WINDOWS"],"property":"system","type":"IN"}
	// ---- Metadata MetadataItems -> {"expectedValue":true,"property":"alerting","type":"EQUALS"}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"alert_activation_duration":    &me.AlertActivationDuration,
		"alerting":                     &me.Alerting,
		"detection_conditions_linux":   &me.DetectionConditionsLinux,
		"detection_conditions_windows": &me.DetectionConditionsWindows,
		"enabled":                      &me.Enabled,
		"metadata":                     &me.Metadata,
		"monitoring":                   &me.Monitoring,
		"name":                         &me.Name,
		"not_installed_alerting":       &me.NotInstalledAlerting,
		"scope":                        &me.Scope,
		"status_condition_linux":       &me.StatusConditionLinux,
		"status_condition_windows":     &me.StatusConditionWindows,
		"system":                       &me.System,
	})
}
