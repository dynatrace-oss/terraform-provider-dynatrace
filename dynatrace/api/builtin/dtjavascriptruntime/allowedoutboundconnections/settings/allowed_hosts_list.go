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

package allowedoutboundconnections

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AllowedHostsList. Specifies allowed hosts and if the allow list should be enforced.
type AllowedHostsList struct {
	Enforced bool     `json:"enforced"`           // If enabled, the Dynatrace JavaScript runtime will only be able to connect to the specified hosts.
	HostList []string `json:"hostList,omitempty"` // The Dynatrace JavaScript runtime will only be to connect to these hosts.
}

func (me *AllowedHostsList) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enforced": {
			Type:        schema.TypeBool,
			Description: "If enabled, the Dynatrace JavaScript runtime will only be able to connect to the specified hosts.",
			Required:    true,
		},
		"host_list": {
			Type:        schema.TypeSet,
			Description: "The Dynatrace JavaScript runtime will only be to connect to these hosts.",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *AllowedHostsList) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enforced":  me.Enforced,
		"host_list": me.HostList,
	})
}

func (me *AllowedHostsList) HandlePreconditions() error {
	// ---- HostList []string -> {"expectedValue":true,"property":"enforced","type":"EQUALS"}
	return nil
}

func (me *AllowedHostsList) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enforced":  &me.Enforced,
		"host_list": &me.HostList,
	})
}
