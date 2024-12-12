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

package microsoftentraidentitydeveloperconnection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID string  `json:"applicationId"`         // Application (client) ID of your app registered in Microsoft Azure App registrations
	ClientSecret  string  `json:"clientSecret"`          // Client secret of your app registered in Microsoft Azure App registrations
	Description   *string `json:"description,omitempty"` // Description
	DirectoryID   string  `json:"directoryId"`           // Directory (tenant) ID of Microsoft Entra Identity Developer
	Name          string  `json:"name"`                  // The name of the Microsoft Entra Identity Developer connection
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "Application (client) ID of your app registered in Microsoft Azure App registrations",
			Required:    true,
			Sensitive:   true,
		},
		"client_secret": {
			Type:        schema.TypeString,
			Description: "Client secret of your app registered in Microsoft Azure App registrations",
			Required:    true,
			Sensitive:   true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"directory_id": {
			Type:        schema.TypeString,
			Description: "Directory (tenant) ID of Microsoft Entra Identity Developer",
			Required:    true,
			Sensitive:   true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the Microsoft Entra Identity Developer connection",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id": "${state.secret_value}",
		"client_secret":  "${state.secret_value}",
		"description":    me.Description,
		"directory_id":   "${state.secret_value}",
		"name":           me.Name,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id": &me.ApplicationID,
		"client_secret":  &me.ClientSecret,
		"description":    &me.Description,
		"directory_id":   &me.DirectoryID,
		"name":           &me.Name,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.DirectoryID != "" {
		me.DirectoryID = "#######"
	}
	if me.ApplicationID != "" {
		me.ApplicationID = "#######"
	}
	if me.ClientSecret != "" {
		me.ClientSecret = "#######"
	}
	return []string{"REST API didn't provide application/directory ID and client secret data"}
}
