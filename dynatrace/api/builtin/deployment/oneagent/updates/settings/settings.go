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

package updates

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type Settings struct {
	MaintenanceWindows MaintenanceWindows `json:"maintenanceWindows,omitempty"` // Maintenance windows
	Revision           *string            `json:"revision,omitempty"`           // Revision
	Scope              *string            `json:"-" scope:"scope"`              // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	TargetVersion      *string            `json:"targetVersion,omitempty"`      // Target version
	UpdateMode         UpdateMode         `json:"updateMode"`                   // Possible Values: `AUTOMATIC`, `AUTOMATIC_DURING_MW`, `MANUAL`
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"maintenance_windows": {
			Type:        schema.TypeList,
			Description: "Maintenance windows",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(MaintenanceWindows).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"revision": {
			Type:        schema.TypeString,
			Description: "Revision",
			Optional:    true, // precondition
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"target_version": {
			Type:        schema.TypeString,
			Description: "Target version",
			Optional:    true, // precondition
		},
		"update_mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `AUTOMATIC`, `AUTOMATIC_DURING_MW`, `MANUAL`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"maintenance_windows": me.MaintenanceWindows,
		"revision":            me.Revision,
		"scope":               me.Scope,
		"target_version":      me.TargetVersion,
		"update_mode":         me.UpdateMode,
	})
}

func (me *Settings) HandlePreconditions() error {
	if me.TargetVersion == nil && slices.Contains([]string{"AUTOMATIC", "AUTOMATIC_DURING_MW"}, string(me.UpdateMode)) {
		return fmt.Errorf("'target_version' must be specified if 'update_mode' is set to '%v'", me.UpdateMode)
	}
	// ---- MaintenanceWindows MaintenanceWindows -> {"expectedValue":"AUTOMATIC_DURING_MW","property":"updateMode","type":"EQUALS"}
	// ---- Revision *string -> {"preconditions":[{"precondition":{"expectedValues":["latest","previous","older"],"property":"targetVersion","type":"IN"},"type":"NOT"},{"expectedValues":["AUTOMATIC","AUTOMATIC_DURING_MW"],"property":"updateMode","type":"IN"}],"type":"AND"}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"maintenance_windows": &me.MaintenanceWindows,
		"revision":            &me.Revision,
		"scope":               &me.Scope,
		"target_version":      &me.TargetVersion,
		"update_mode":         &me.UpdateMode,
	})
}
