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

package internetproxy

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The configuration of the user
type Settings struct {
	Scheme        Scheme   `json:"scheme"`                  // Protocol which proxy server uses
	Server        string   `json:"server"`                  // Address (either IP or Hostname) of proxy server
	Port          int      `json:"port"`                    // Port of proxy server
	User          *string  `json:"user,omitempty"`          // User of proxy server, null means do not change previous value
	Password      *string  `json:"password,omitempty"`      // Password of proxy server, null means do not change previous value
	NonProxyHosts []string `json:"nonProxyHosts,omitempty"` // Definition of hosts for which proxy won't be used. You can define multiple hosts. Each host can start or end with wildcard '*' for instance to match whole domain.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"scheme": {
			Type:        schema.TypeString,
			Description: "Protocol which proxy server uses",
			Required:    true,
		},
		"server": {
			Type:        schema.TypeString,
			Description: "Address (either IP or Hostname) of proxy server",
			Required:    true,
		},
		"port": {
			Type:        schema.TypeInt,
			Description: "Port of proxy server",
			Required:    true,
		},
		"user": {
			Type:        schema.TypeString,
			Description: "User of proxy server, null means do not change previous value",
			Optional:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Password of proxy server, null means do not change previous value",
			Optional:    true,
			Sensitive:   true,
		},
		"non_proxy_hosts": {
			Type:        schema.TypeSet,
			Description: "Definition of hosts for which proxy won't be used. You can define multiple hosts. Each host can start or end with wildcard '*' for instance to match whole domain.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"scheme":          me.Scheme,
		"server":          me.Server,
		"port":            me.Port,
		"non_proxy_hosts": me.NonProxyHosts,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"scheme":          &me.Scheme,
		"server":          &me.Server,
		"port":            &me.Port,
		"user":            &me.User,
		"password":        &me.Password,
		"non_proxy_hosts": &me.NonProxyHosts,
	})
}
