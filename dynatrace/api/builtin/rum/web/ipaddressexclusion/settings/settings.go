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

package ipaddressexclusion

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID             string                  `json:"-" scope:"applicationId"`   // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	IpAddressExclusionInclude bool                    `json:"ipAddressExclusionInclude"` // These are the only IP addresses that should be monitored
	IpExclusionList           IpAddressExclusionRules `json:"ipExclusionList,omitempty"` // **Examples:**\n\n   - 84.112.10.5\n   - fe80::10a1:c6b2:5f68:785d
}

func (me *Settings) Name() string {
	return me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
		},
		"ip_address_exclusion_include": {
			Type:        schema.TypeBool,
			Description: "These are the only IP addresses that should be monitored",
			Required:    true,
		},
		"ip_exclusion_list": {
			Type:        schema.TypeList,
			Description: "**Examples:**\n\n   - 84.112.10.5\n   - fe80::10a1:c6b2:5f68:785d",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(IpAddressExclusionRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":               me.ApplicationID,
		"ip_address_exclusion_include": me.IpAddressExclusionInclude,
		"ip_exclusion_list":            me.IpExclusionList,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":               &me.ApplicationID,
		"ip_address_exclusion_include": &me.IpAddressExclusionInclude,
		"ip_exclusion_list":            &me.IpExclusionList,
	})
}
