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

type FederatedIdentityCredential struct {
	ApplicationID *string                                  `json:"applicationId,omitempty"` // Application (client) ID of your app registered in Microsoft Azure App registrations
	Consumers     []ConsumersOfFederatedIdentityCredential `json:"consumers,omitempty"`     // Consumers that can use the connection. Possible Values: `APP:dynatrace.microsoft.azure.connector`, `DA`, `NONE`, `SVC:com.dynatrace.da`, `SVC:com.dynatrace.openpipeline`
	DirectoryID   *string                                  `json:"directoryId,omitempty"`   // Directory (tenant) ID of Microsoft Entra ID
}

func (me *FederatedIdentityCredential) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"consumers": {
			Type:        schema.TypeList,
			Description: "Consumers that can use the connection. Possible Values: `APP:dynatrace.microsoft.azure.connector`, `DA`, `NONE`, `SVC:com.dynatrace.da`, `SVC:com.dynatrace.openpipeline`",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
			ForceNew:    true,
		},
	}
}

func (me *FederatedIdentityCredential) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"consumers": me.Consumers,
	})
}

func (me *FederatedIdentityCredential) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"consumers": &me.Consumers,
	})
}
