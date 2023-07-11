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

package vmware

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled   bool    `json:"enabled"`          // This setting is enabled (`true`) or disabled (`false`)
	Filter    *string `json:"filter,omitempty"` // This string should have one of the following formats:\n- $prefix(parameter) - property value starting with 'parameter'\n- $eq(parameter) - property value exactly matching 'parameter'\n- $suffix(parameter) - property value ends with 'parameter'\n- $contains(parameter) - property value contains 'parameter'
	Ipaddress string  `json:"ipaddress"`        // Specify the IP address or name of the vCenter or standalone ESXi host:
	Label     string  `json:"label"`            // Name this connection
	Password  string  `json:"password"`
	Username  string  `json:"username"` // Provide user credentials for the vCenter or standalone ESXi host:
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"filter": {
			Type:        schema.TypeString,
			Description: "This string should have one of the following formats:\n- $prefix(parameter) - property value starting with 'parameter'\n- $eq(parameter) - property value exactly matching 'parameter'\n- $suffix(parameter) - property value ends with 'parameter'\n- $contains(parameter) - property value contains 'parameter'",
			Optional:    true, // nullable
		},
		"ipaddress": {
			Type:        schema.TypeString,
			Description: "Specify the IP address or name of the vCenter or standalone ESXi host:",
			Required:    true,
		},
		"label": {
			Type:        schema.TypeString,
			Description: "Name this connection",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
			Sensitive:   true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "Provide user credentials for the vCenter or standalone ESXi host:",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":   me.Enabled,
		"filter":    me.Filter,
		"ipaddress": me.Ipaddress,
		"label":     me.Label,
		"password":  "${state.secret_value}",
		"username":  me.Username,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":   &me.Enabled,
		"filter":    &me.Filter,
		"ipaddress": &me.Ipaddress,
		"label":     &me.Label,
		"password":  &me.Password,
		"username":  &me.Username,
	})
}

const credsNotProvided = "REST API didn't provide credential data"

func (me *Settings) FillDemoValues() []string {
	me.Password = "################"
	return []string{credsNotProvided}
}
