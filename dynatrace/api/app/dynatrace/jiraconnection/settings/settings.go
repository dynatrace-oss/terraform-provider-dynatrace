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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type Settings struct {
	Name        string  `json:"name"`               // The name of the Jira connection
	Password    *string `json:"password,omitempty"` // Password of the Jira user
	Token       *string `json:"token,omitempty"`    // Token for the selected authentication type
	Type        Type    `json:"type"`               // Possible Values: `Basic`, `Cloud_token`, `Pat`
	Url         string  `json:"url"`                // URL of the Jira server
	User        *string `json:"user,omitempty"`     // Username or E-Mail address
	InsertAfter string  `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the Jira connection",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Password of the Jira user",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "Token for the selected authentication type",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Basic`, `Cloud_token`, `Pat`",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "URL of the Jira server",
			Required:    true,
		},
		"user": {
			Type:        schema.TypeString,
			Description: "Username or E-Mail address",
			Optional:    true, // precondition
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":         me.Name,
		"password":     "${state.secret_value}",
		"token":        "${state.secret_value}",
		"type":         me.Type,
		"url":          me.Url,
		"user":         me.User,
		"insert_after": me.InsertAfter,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Password == nil) && (string(me.Type) == "basic") {
		return fmt.Errorf("'password' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Token == nil) && (slices.Contains([]string{"pat", "cloud-token"}, string(me.Type))) {
		return fmt.Errorf("'token' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.User == nil) && (slices.Contains([]string{"basic", "cloud-token"}, string(me.Type))) {
		return fmt.Errorf("'user' must be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.Name,
		"password":     &me.Password,
		"token":        &me.Token,
		"type":         &me.Type,
		"url":          &me.Url,
		"user":         &me.User,
		"insert_after": &me.InsertAfter,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.Password != nil {
		me.Password = opt.NewString("#######")
	}
	if me.Token != nil {
		me.Token = opt.NewString("#######")
	}
	return []string{"REST API didn't provide password/token data"}
}
