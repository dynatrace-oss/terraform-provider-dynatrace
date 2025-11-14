/**
* @license
* Copyright 2025 Dynatrace LLC
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

package azure

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ClientSecretConfig struct {
	ApplicationID string                    `json:"applicationId"`       // Application (client) ID of your app registered in Microsoft Azure App registrations
	ClientSecret  string                    `json:"clientSecret"`        // Client secret of your app registered in Microsoft Azure App registrations
	Consumers     []ConsumersOfClientSecret `json:"consumers,omitempty"` // Dynatrace integrations that can use this connection. Possible Values: `DA`, `NONE`, `SVC:com.dynatrace.da`
	DirectoryID   string                    `json:"directoryId"`         // Directory (tenant) ID of Microsoft Entra ID
}

func (me *ClientSecretConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "Application (client) ID of your app registered in Microsoft Azure App registrations",
			Required:    true,
			ForceNew:    true,
		},
		"client_secret": {
			Type:        schema.TypeString,
			Description: "Client secret of your app registered in Microsoft Azure App registrations",
			Required:    true,
			Sensitive:   true,
		},
		"consumers": {
			Type:        schema.TypeList,
			Description: "Dynatrace integrations that can use this connection. Possible Values: `DA`, `NONE`, `SVC:com.dynatrace.da`",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
			ForceNew:    true,
		},
		"directory_id": {
			Type:        schema.TypeString,
			Description: "Directory (tenant) ID of Microsoft Entra ID",
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (me *ClientSecretConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id": me.ApplicationID,
		"client_secret":  me.ClientSecret,
		"consumers":      me.Consumers,
		"directory_id":   me.DirectoryID,
	})
}

func (me *ClientSecretConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id": &me.ApplicationID,
		"client_secret":  &me.ClientSecret,
		"consumers":      &me.Consumers,
		"directory_id":   &me.DirectoryID,
	})
}
