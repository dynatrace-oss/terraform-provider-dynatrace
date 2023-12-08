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

package cloudfoundry

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ActiveGateGroup *string `json:"activeGateGroup,omitempty"` // ActiveGate group
	ApiUrl          string  `json:"apiUrl"`                    // Cloud Foundry API Target
	Enabled         bool    `json:"enabled"`                   // This setting is enabled (`true`) or disabled (`false`)
	Label           string  `json:"label"`                     // Name this connection
	LoginUrl        string  `json:"loginUrl"`                  // Cloud Foundry Authentication Endpoint
	Password        string  `json:"password"`                  // Cloud Foundry Password
	Username        string  `json:"username"`                  // Cloud Foundry Username
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"active_gate_group": {
			Type:        schema.TypeString,
			Description: "ActiveGate group",
			Optional:    true, // nullable
		},
		"api_url": {
			Type:        schema.TypeString,
			Description: "Cloud Foundry API Target",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"label": {
			Type:        schema.TypeString,
			Description: "Name this connection",
			Required:    true,
		},
		"login_url": {
			Type:        schema.TypeString,
			Description: "Cloud Foundry Authentication Endpoint",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Cloud Foundry Password",
			Required:    true,
			Sensitive:   true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "Cloud Foundry Username",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMap(
		me.Schema(),
		map[string]any{
			"active_gate_group": me.ActiveGateGroup,
			"api_url":           me.ApiUrl,
			"enabled":           me.Enabled,
			"label":             me.Label,
			"login_url":         me.LoginUrl,
			"password":          "${state.secret_value}",
			"username":          me.Username,
		},
	))
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"active_gate_group": &me.ActiveGateGroup,
		"api_url":           &me.ApiUrl,
		"enabled":           &me.Enabled,
		"label":             &me.Label,
		"login_url":         &me.LoginUrl,
		"password":          &me.Password,
		"username":          &me.Username,
	})
}

func (me *Settings) FillDemoValues() []string {
	me.Password = "#######"
	return []string{"Please fill in the correct password"}
}
