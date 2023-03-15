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

package activegatetoken

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AuthTokenEnforcementManuallyEnabled bool `json:"authTokenEnforcementManuallyEnabled"` // Manually enforce ActiveGate token authentication
	ExpiringTokenNotificationsEnabled   bool `json:"expiringTokenNotificationsEnabled"`   // Note: ActiveGate tokens notifications are sent only when you deployed ActiveGate tokens with expiration dates in your environment and notifications are defined ([see notification settings](/ui/settings/builtin:problem.notifications))
}

func (me *Settings) Name() string {
	return "activegate_token"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auth_token_enforcement_manually_enabled": {
			Type:        schema.TypeBool,
			Description: "Manually enforce ActiveGate token authentication",
			Required:    true,
		},
		"expiring_token_notifications_enabled": {
			Type:        schema.TypeBool,
			Description: "Note: ActiveGate tokens notifications are sent only when you deployed ActiveGate tokens with expiration dates in your environment and notifications are defined ([see notification settings](/ui/settings/builtin:problem.notifications))",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auth_token_enforcement_manually_enabled": me.AuthTokenEnforcementManuallyEnabled,
		"expiring_token_notifications_enabled":    me.ExpiringTokenNotificationsEnabled,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auth_token_enforcement_manually_enabled": &me.AuthTokenEnforcementManuallyEnabled,
		"expiring_token_notifications_enabled":    &me.ExpiringTokenNotificationsEnabled,
	})
}
