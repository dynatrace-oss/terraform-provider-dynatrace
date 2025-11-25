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

package mgmz

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Permission struct {
	EnvironmentID    string
	GroupID          string
	ManagementZoneID string
	Permissions      []string
}

func (me *Permission) Name() string {
	return fmt.Sprintf("%s#-#%s#-#%s", me.EnvironmentID, me.GroupID, me.ManagementZoneID)
}

func (me *Permission) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"environment": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The UUID of the environment",
		},
		"group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the group the permissions are valid for. You may refer to the id of a resource `dynatrace_user_group` here",
		},
		"management_zone": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the management zone the permissions are valid for. When referring to resource `dynatrace_management_zone_v2` or data source `dynatrace_management_zone` you need to refer to the attribute `legacy_id`.",
		},
		"permissions": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "The permissions to assign for that management zone. Allowed values are `DEMO_USER`, `LOG_VIEWER`, `MANAGE_SECURITY_PROBLEMS`, `MANAGE_SETTINGS`, `REPLAY_SESSION_DATA`, `REPLAY_SESSION_DATA_WITHOUT_MASKING`, `VIEWER`, `VIEW_SENSITIVE_REQUEST_DATA`.\nNote: In order to produce non-empty plans specifying at least the permission `VIEWER` is recommended. Your Dynatrace Cluster will enforce that permission, regardless of whether it has been specified or not.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
				/* Validation disabled because of #563. There exist meanwhile many more permissions - some of them enforced by the CLUSTER REST API */
				// ValidateFunc: validation.StringInSlice([]string{"DEMO_USER", "LOG_VIEWER", "MANAGE_SECURITY_PROBLEMS", "MANAGE_SETTINGS", "REPLAY_SESSION_DATA", "REPLAY_SESSION_DATA_WITHOUT_MASKING", "VIEWER", "VIEW_SENSITIVE_REQUEST_DATA"}, false),
			},
		},
	}
}

func (me *Permission) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"environment":     me.EnvironmentID,
		"group":           me.GroupID,
		"management_zone": me.ManagementZoneID,
		"permissions":     me.Permissions,
	})
}

func (me *Permission) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"environment":     &me.EnvironmentID,
		"group":           &me.GroupID,
		"management_zone": &me.ManagementZoneID,
		"permissions":     &me.Permissions,
	})
}
