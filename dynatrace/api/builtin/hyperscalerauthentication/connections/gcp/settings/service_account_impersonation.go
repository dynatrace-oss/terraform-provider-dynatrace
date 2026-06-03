/**
* @license
* Copyright 2026 Dynatrace LLC
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

package gcp

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceAccountImpersonation struct {
	Consumers        []ConsumersOfServiceAccountImpersonation `json:"consumers,omitempty"`        // Dynatrace integrations that can use this connection. Possible values: `SVC:com.dynatrace.bo`, `SVC:com.dynatrace.da`, `SVC:com.dynatrace.openpipeline`
	ServiceAccountID *string                                  `json:"serviceAccountId,omitempty"` // Id of your Service Account
}

func (me *ServiceAccountImpersonation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"consumers": {
			Type:        schema.TypeList,
			Description: "Dynatrace integrations that can use this connection. Possible values: `SVC:com.dynatrace.bo`, `SVC:com.dynatrace.da`, `SVC:com.dynatrace.openpipeline`",
			Optional:    true, // minobjects == 0
			ForceNew:    true, // Changing consumers after creation would lead to validation errors on the API side
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"service_account_id": {
			Type:        schema.TypeString,
			Description: "Id of your Service Account",
			Required:    true, // Even though the schema does not require it, we want the user to provide it already on creation
			ForceNew:    true, // Changing the service account after it has been set once would lead to validation errors on the API side
		},
	}
}

func (me *ServiceAccountImpersonation) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"consumers":          me.Consumers,
		"service_account_id": me.ServiceAccountID,
	})
}

func (me *ServiceAccountImpersonation) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"consumers":          &me.Consumers,
		"service_account_id": &me.ServiceAccountID,
	})
}
