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

package approval

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	WorkflowAppAccessApprovalEnabled bool `json:"workflowAppAccessApprovalEnabled"` // Allow on tenant level anyone with access to the app can respond to requests via an approval link.
}

func (me *Settings) Name() string {
	return "environment"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"workflow_app_access_approval_enabled": {
			Type:        schema.TypeBool,
			Description: "Allow on tenant level anyone with access to the app can respond to requests via an approval link.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"workflow_app_access_approval_enabled": me.WorkflowAppAccessApprovalEnabled,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"workflow_app_access_approval_enabled": &me.WorkflowAppAccessApprovalEnabled,
	})
}
