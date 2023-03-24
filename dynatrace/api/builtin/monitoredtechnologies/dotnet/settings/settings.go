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

package dotnet

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled           bool    `json:"enabled"`                     // This setting is enabled (`true`) or disabled (`false`)
	EnabledDotNetCore *bool   `json:"enabledDotNetCore,omitempty"` // Requires Dynatrace OneAgent version 1.117 or later on Windows and version 1.127 or later on Linux and .NET monitoring enabled
	HostID            *string `json:"-" scope:"hostId"`            // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
}

func (me *Settings) Name() string {
	return *me.HostID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"enabled_dot_net_core": {
			Type:        schema.TypeBool,
			Description: "Requires Dynatrace OneAgent version 1.117 or later on Windows and version 1.127 or later on Linux and .NET monitoring enabled",
			Optional:    true,
		},
		"host_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":              me.Enabled,
		"enabled_dot_net_core": me.EnabledDotNetCore,
		"host_id":              me.HostID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"enabled":              &me.Enabled,
		"enabled_dot_net_core": &me.EnabledDotNetCore,
		"host_id":              &me.HostID,
	})
	if me.EnabledDotNetCore == nil && me.Enabled {
		me.EnabledDotNetCore = opt.NewBool(false)
	}
	return err
}
