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
	Name     string  `json:"name"`               // A unique and clearly identifiable connection name to your ServiceNow instance.
	Password *string `json:"password,omitempty"` // Password of the ServiceNow user.
	Type     Type    `json:"type"`               // Possible Values: `basic`
	Url      string  `json:"url"`                // URL of the ServiceNow instance.
	User     *string `json:"user,omitempty"`     // Username or Email address.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "A unique and clearly identifiable connection name to your ServiceNow instance.",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Password of the ServiceNow user.",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `basic`",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "URL of the ServiceNow instance.",
			Required:    true,
		},
		"user": {
			Type:        schema.TypeString,
			Description: "Username or Email address.",
			Optional:    true, // precondition
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":     me.Name,
		"password": "${state.secret_value}",
		"type":     me.Type,
		"url":      me.Url,
		"user":     me.User,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Password == nil) && (slices.Contains([]string{"basic"}, string(me.Type))) {
		return fmt.Errorf("'password' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.User == nil) && (slices.Contains([]string{"basic"}, string(me.Type))) {
		return fmt.Errorf("'user' must be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":     &me.Name,
		"password": &me.Password,
		"type":     &me.Type,
		"url":      &me.Url,
		"user":     &me.User,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.Password != nil {
		me.Password = opt.NewString("#######")
	}
	return []string{"REST API didn't provide password data"}
}
