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

package fullwebrequest

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// idContributorsType. Defines which detected values contribute to the Service Id for matching full web requests.
type IdContributorsType struct {
	ApplicationID *ServiceIdContributor `json:"applicationId"` // Contribute to the Service Id calculation from the detected application identifier.. You can keep the detected value, override it with a constant value, or apply transformations before it contributes to the Service Id.
	ContextRoot   *ContextIdContributor `json:"contextRoot"`   // The context root is the first segment of the request URL after the Server name. For example, in the `www.dynatrace.com/support/help/dynatrace-api/` URL the context root is `/support`. The context root value can be found on the **Service overview page** under **Properties and tags**.
	ServerName    *ServiceIdContributor `json:"serverName"`    // Contribute to the Service Id calculation from the detected server name.
}

func (me *IdContributorsType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeList,
			Description: "Contribute to the Service Id calculation from the detected application identifier.. You can keep the detected value, override it with a constant value, or apply transformations before it contributes to the Service Id.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"context_root": {
			Type:        schema.TypeList,
			Description: "The context root is the first segment of the request URL after the Server name. For example, in the `www.dynatrace.com/support/help/dynatrace-api/` URL the context root is `/support`. The context root value can be found on the **Service overview page** under **Properties and tags**.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ContextIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"server_name": {
			Type:        schema.TypeList,
			Description: "Contribute to the Service Id calculation from the detected server name.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *IdContributorsType) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id": me.ApplicationID,
		"context_root":   me.ContextRoot,
		"server_name":    me.ServerName,
	})
}

func (me *IdContributorsType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id": &me.ApplicationID,
		"context_root":   &me.ContextRoot,
		"server_name":    &me.ServerName,
	})
}
