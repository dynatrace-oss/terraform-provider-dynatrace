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

package connection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ChannelName *string `json:"channelName,omitempty"` // Optional
	Name        string  `json:"name"`                  // The name of the Microsoft Teams connection
	TeamName    *string `json:"teamName,omitempty"`    // Optional
	Webhook     string  `json:"webhook"`               // The Webhook URL that links to the channel
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"channel_name": {
			Type:        schema.TypeString,
			Description: "Optional",
			Optional:    true, // nullable
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the Microsoft Teams connection",
			Required:    true,
		},
		"team_name": {
			Type:        schema.TypeString,
			Description: "Optional",
			Optional:    true, // nullable
		},
		"webhook": {
			Type:        schema.TypeString,
			Description: "The Webhook URL that links to the channel",
			Required:    true,
			Sensitive:   true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"channel_name": me.ChannelName,
		"name":         me.Name,
		"team_name":    me.TeamName,
		"webhook":      "${state.secret_value}",
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"channel_name": &me.ChannelName,
		"name":         &me.Name,
		"team_name":    &me.TeamName,
		"webhook":      &me.Webhook,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.Webhook != "" {
		me.Webhook = "#######"
	}
	return []string{"REST API didn't provide token data"}
}
