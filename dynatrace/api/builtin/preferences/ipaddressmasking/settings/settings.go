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

package ipaddressmasking

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled bool                    `json:"enabled"`         // This setting is enabled (`true`) or disabled (`false`)
	Scope   *string                 `json:"-" scope:"scope"` // The scope of this setting (MOBILE_APPLICATION, CUSTOM_APPLICATION, APPLICATION). Omit this property if you want to cover the whole environment.
	Type    *IpAddressMaskingOption `json:"type,omitempty"`  // Possible Values: `All`, `Public`
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (MOBILE_APPLICATION, CUSTOM_APPLICATION, APPLICATION). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `all`, `public`",
			Optional:    true, // precondition
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"scope":   me.Scope,
		"type":    me.Type,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Type == nil) && (me.Enabled) {
		return fmt.Errorf("'type' must be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	if (me.Type != nil) && (!me.Enabled) {
		return fmt.Errorf("'type' must not be specified if 'enabled' is set to '%v'", me.Enabled)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled": &me.Enabled,
		"scope":   &me.Scope,
		"type":    &me.Type,
	})
}
