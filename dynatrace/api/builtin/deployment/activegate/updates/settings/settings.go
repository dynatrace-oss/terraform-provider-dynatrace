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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Scope         *string       `json:"-" scope:"scope"`         // The scope of this setting (ENVIRONMENT_ACTIVE_GATE). Omit this property if you want to cover the whole environment.
	TargetVersion string        `json:"targetVersion"`           // Target version
	UpdateMode    UpdateMode    `json:"updateMode"`              // Update mode. Possible values: `AUTOMATIC`, `AUTOMATIC_DURING_UW`, `MANUAL`
	UpdateWindows UpdateWindows `json:"updateWindows,omitempty"` // Update windows
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (ENVIRONMENT_ACTIVE_GATE). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"target_version": {
			Type:        schema.TypeString,
			Description: "Target version",
			Required:    true,
		},
		"update_mode": {
			Type:        schema.TypeString,
			Description: "Update mode. Possible values: `AUTOMATIC`, `AUTOMATIC_DURING_UW`, `MANUAL`",
			Required:    true,
		},
		"update_windows": {
			Type:        schema.TypeList,
			Description: "Update windows",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(UpdateWindows).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"scope":          me.Scope,
		"target_version": me.TargetVersion,
		"update_mode":    me.UpdateMode,
		"update_windows": me.UpdateWindows,
	})
}

func (me *Settings) HandlePreconditions() error {
	// ---- UpdateWindows UpdateWindows -> {"expectedValue":"AUTOMATIC_DURING_UW","property":"updateMode","type":"EQUALS"}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"scope":          &me.Scope,
		"target_version": &me.TargetVersion,
		"update_mode":    &me.UpdateMode,
		"update_windows": &me.UpdateWindows,
	})
}
