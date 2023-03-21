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

package general

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DefaultDashboardList UserGroupss `json:"defaultDashboardList"` // Configure home dashboard for selected user group. The selected preset dashboard will be loaded as default landing page for this environment.
	EnablePublicSharing  bool        `json:"enablePublicSharing"`  // Allow users to grant anonymous access to dashboards. No sign-in will be required to view those dashboards read-only.
}

func (me *Settings) Name() string {
	return "dashboards_general"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_dashboard_list": {
			Type:        schema.TypeList,
			Description: "Configure home dashboard for selected user group. The selected preset dashboard will be loaded as default landing page for this environment.",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(UserGroupss).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"enable_public_sharing": {
			Type:        schema.TypeBool,
			Description: "Allow users to grant anonymous access to dashboards. No sign-in will be required to view those dashboards read-only.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"default_dashboard_list": me.DefaultDashboardList,
		"enable_public_sharing":  me.EnablePublicSharing,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"default_dashboard_list": &me.DefaultDashboardList,
		"enable_public_sharing":  &me.EnablePublicSharing,
	})
}
