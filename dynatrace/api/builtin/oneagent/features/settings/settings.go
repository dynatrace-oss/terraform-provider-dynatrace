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

package features

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled         bool    `json:"enabled"`                   // This setting is enabled (`true`) or disabled (`false`)
	Forcible        *bool   `json:"forcible,omitempty"`        // Activate this feature also in OneAgents only fulfilling the minimum Opt-In version
	Instrumentation *bool   `json:"instrumentation,omitempty"` // Instrumentation enabled (change needs a process restart)
	Key             string  `json:"key"`                       // Feature
	Scope           *string `json:"-" scope:"scope"`           // The scope of this setting (PROCESS_GROUP_INSTANCE, PROCESS_GROUP). Omit this property if you want to cover the whole environment.
	// OneAgent Features cannot get deleted if their scope is the environment
	// Settings with this exact Property are being treated special by the provider
	// Upon Create the current state of the setting is getting stored within the state file
	// Upon Delete instead of deleting the setting, that stored state is getting reconstructed
	// Also the method `Match` needs to be present in order for that functionality to kick in
	RestoreOnDelete *string `json:"-"` // For Terraform internal use. Do not populate.
}

func (me *Settings) Name() string {
	if me.Scope == nil || len(*me.Scope) == 0 || *me.Scope == "environment" {
		return me.Key
	}
	return me.Key + "_" + *me.Scope
}

func (me *Settings) Match(o any) bool {
	if other, ok := o.(*Settings); ok {
		if other == nil {
			return false
		}
		if me == nil {
			return false
		}
		if other.Key != me.Key {
			return false
		}
		scope := "environment"
		otherScope := "environment"
		if other.Scope != nil && len(*other.Scope) > 0 {
			otherScope = *other.Scope
		}
		if me.Scope != nil && len(*me.Scope) > 0 {
			scope = *me.Scope
		}
		return scope == otherScope
	}
	return false
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"forcible": {
			Type:        schema.TypeBool,
			Description: "Activate this feature also in OneAgents only fulfilling the minimum Opt-In version",
			Optional:    true,
		},
		"instrumentation": {
			Type:        schema.TypeBool,
			Description: "Instrumentation enabled (change needs a process restart)",
			Optional:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Feature",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (PROCESS_GROUP_INSTANCE, PROCESS_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"_restore_": {
			Type:        schema.TypeString,
			Description: "Used internally by the terraform provider. Do not populate",
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":         me.Enabled,
		"forcible":        me.Forcible,
		"instrumentation": me.Instrumentation,
		"key":             me.Key,
		"scope":           me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":         &me.Enabled,
		"forcible":        &me.Forcible,
		"instrumentation": &me.Instrumentation,
		"key":             &me.Key,
		"scope":           &me.Scope,
	})
}
