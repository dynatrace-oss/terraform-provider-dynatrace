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
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	EventStreamEnabled *bool   `json:"eventStreamEnabled,omitempty"` // Flag if Red Hat Event Stream is use for Event-Driven Ansible
	Name               string  `json:"name"`                         // A unique and clearly identifiable connection name.
	Token              *string `json:"token,omitempty"`              // API access token for the Event-Driven Ansible Controller. Please note that this token is not refreshed and can expire.
	Type               Type    `json:"type"`                         // Possible Values: `Api_token`
	Url                string  `json:"url"`                          // URL of the Event-Driven Ansible source plugin webhook. For example, https://eda.yourdomain.com:5010
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_stream_enabled": {
			Type:        schema.TypeBool,
			Description: "Flag if Red Hat Event Stream is use for Event-Driven Ansible",
			Optional:    true, // nullable
		},
		"name": {
			Type:        schema.TypeString,
			Description: "A unique and clearly identifiable connection name.",
			Required:    true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "API access token for the Event-Driven Ansible Controller. Please note that this token is not refreshed and can expire.",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Api_token`",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "URL of the Event-Driven Ansible source plugin webhook. For example, https://eda.yourdomain.com:5010",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_stream_enabled": me.EventStreamEnabled,
		"name":                 me.Name,
		"token":                "${state.secret_value}",
		"type":                 me.Type,
		"url":                  me.Url,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Token == nil) && (slices.Contains([]string{"api-token"}, string(me.Type))) {
		return fmt.Errorf("'token' must be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_stream_enabled": &me.EventStreamEnabled,
		"name":                 &me.Name,
		"token":                &me.Token,
		"type":                 &me.Type,
		"url":                  &me.Url,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.Token != nil {
		me.Token = opt.NewString("#######")
	}
	return []string{"REST API didn't provide token data"}
}
