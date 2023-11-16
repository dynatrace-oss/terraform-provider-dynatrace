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

package advanced

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	HostID                string `json:"-" scope:"hostId"`      // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	ProcessAgentInjection bool   `json:"processAgentInjection"` // Disabling this via [oneagentctl](https://dt-url.net/oneagentctl) takes precedence over this setting and cannot be changed from the Dynatrace web UI.
}

func (me *Settings) Name() string {
	return me.HostID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
		},
		"process_agent_injection": {
			Type:        schema.TypeBool,
			Description: "Disabling this via [oneagentctl](https://dt-url.net/oneagentctl) takes precedence over this setting and cannot be changed from the Dynatrace web UI.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"host_id":                 me.HostID,
		"process_agent_injection": me.ProcessAgentInjection,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"host_id":                 &me.HostID,
		"process_agent_injection": &me.ProcessAgentInjection,
	})
}
