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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Object scopes. Scopes controlling which objects are affected by the maintenance window.
type ObjectScopes struct {
	SuppressProblems  *bool                     `json:"suppressProblems,omitempty"` // When enabled, Davis does not open new problems in the window's scope while the maintenance window is running.
	SyntheticMonitors *SyntheticMonitorSettings `json:"syntheticMonitors"`          // Synthetic monitors
}

func (me *ObjectScopes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"suppress_problems": {
			Type:        schema.TypeBool,
			Description: "When enabled, Davis does not open new problems in the window's scope while the maintenance window is running.",
			Optional:    true, // nullable
		},
		"synthetic_monitors": {
			Type:        schema.TypeList,
			Description: "Synthetic monitors",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SyntheticMonitorSettings).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *ObjectScopes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"suppress_problems":  me.SuppressProblems,
		"synthetic_monitors": me.SyntheticMonitors,
	})
}

func (me *ObjectScopes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"suppress_problems":  &me.SuppressProblems,
		"synthetic_monitors": &me.SyntheticMonitors,
	})
}
