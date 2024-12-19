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
)

type Settings struct {
	Client_id     string  `json:"client_id"`               // Application (client) ID of your app registered in Microsoft Azure App registrations
	Client_secret *string `json:"client_secret,omitempty"` // Client secret of your app registered in Microsoft Azure App registrations
	From_address  string  `json:"from_address"`            // The email address from which the messages will be sent
	Name          string  `json:"name"`                    // A unique name for the Microsoft 365 email connection
	Tenant_id     string  `json:"tenant_id"`               // Directory (tenant) ID of your Azure Active Directory
	Type          Type    `json:"type"`                    // Possible Values: `client_secret`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_id": {
			Type:        schema.TypeString,
			Description: "Application (client) ID of your app registered in Microsoft Azure App registrations",
			Required:    true,
		},
		"client_secret": {
			Type:        schema.TypeString,
			Description: "Client secret of your app registered in Microsoft Azure App registrations",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"from_address": {
			Type:        schema.TypeString,
			Description: "The email address from which the messages will be sent",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "A unique name for the Microsoft 365 email connection",
			Required:    true,
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Description: "Directory (tenant) ID of your Azure Active Directory",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `client_secret`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"client_id":     me.Client_id,
		"client_secret": "${state.secret_value}",
		"from_address":  me.From_address,
		"name":          me.Name,
		"tenant_id":     me.Tenant_id,
		"type":          me.Type,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Client_secret == nil) && (string(me.Type) == "client_secret") {
		return fmt.Errorf("'client_secret' must be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"client_id":     &me.Client_id,
		"client_secret": &me.Client_secret,
		"from_address":  &me.From_address,
		"name":          &me.Name,
		"tenant_id":     &me.Tenant_id,
		"type":          &me.Type,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.Client_secret != nil {
		me.Client_secret = opt.NewString("#######")
	}
	return []string{"REST API didn't provide secret data"}
}
