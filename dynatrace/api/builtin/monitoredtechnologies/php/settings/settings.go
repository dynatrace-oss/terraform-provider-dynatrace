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

package php

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	EnablePhpCliServerInstrumentation *bool   `json:"enablePhpCliServerInstrumentation,omitempty"` // Requires enabled PHP monitoring and Dynatrace OneAgent version 1.261 or later
	Enabled                           bool    `json:"enabled"`                                     // This setting is enabled (`true`) or disabled (`false`)
	EnabledFastCGI                    *bool   `json:"enabledFastCGI,omitempty"`                    // Requires PHP monitoring enabled and from Dynatrace OneAgent version 1.191 it's ignored and permanently enabled
	HostID                            *string `json:"-" scope:"hostId"`                            // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
}

func (me *Settings) Name() string {
	return *me.HostID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_php_cli_server": {
			Type:        schema.TypeBool,
			Description: "Requires enabled PHP monitoring and Dynatrace OneAgent version 1.261 or later",
			Optional:    true, // precondition
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"enabled_fast_cgi": {
			Type:        schema.TypeBool,
			Description: "Requires PHP monitoring enabled and from Dynatrace OneAgent version 1.191 it's ignored and permanently enabled",
			Optional:    true, // precondition
		},
		"host_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enable_php_cli_server": me.EnablePhpCliServerInstrumentation,
		"enabled":               me.Enabled,
		"enabled_fast_cgi":      me.EnabledFastCGI,
		"host_id":               me.HostID,
	})
}

func (me *Settings) HandlePreconditions() error {
	if me.EnablePhpCliServerInstrumentation == nil && me.Enabled {
		me.EnablePhpCliServerInstrumentation = opt.NewBool(false)
	}
	if me.EnabledFastCGI == nil && me.Enabled {
		me.EnabledFastCGI = opt.NewBool(false)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enable_php_cli_server": &me.EnablePhpCliServerInstrumentation,
		"enabled":               &me.Enabled,
		"enabled_fast_cgi":      &me.EnabledFastCGI,
		"host_id":               &me.HostID,
	})
}
