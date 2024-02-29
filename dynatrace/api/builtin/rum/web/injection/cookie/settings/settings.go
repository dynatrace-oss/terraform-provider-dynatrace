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

package cookie

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID            string                  `json:"-" scope:"applicationId"`         // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	CookiePlacementDomain    *string                 `json:"cookiePlacementDomain,omitempty"` // Specify an alternative domain for cookies set by Dynatrace. Keep in mind that your browser may not allow placement of cookies on certain domains (for example, top-level domains). Before typing a domain name here, confirm that the domain will accept cookies from your browser. For details, see the list of [forbidden top-level domains](https://dt-url.net/9n6b0pfz).
	SameSiteCookieAttribute  SameSiteCookieAttribute `json:"sameSiteCookieAttribute"`         // Possible Values: `LAX`, `NONE`, `NOTSET`, `STRICT`
	UseSecureCookieAttribute bool                    `json:"useSecureCookieAttribute"`        // If your application is only accessible via SSL, you can add the Secure attribute to all cookies set by Dynatrace. This setting prevents the display of warnings from PCI-compliance security scanners. Be aware that with this setting enabled Dynatrace correlation of user actions with server-side web requests is only possible over SSL connections.
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
		"cookie_placement_domain": {
			Type:        schema.TypeString,
			Description: "Specify an alternative domain for cookies set by Dynatrace. Keep in mind that your browser may not allow placement of cookies on certain domains (for example, top-level domains). Before typing a domain name here, confirm that the domain will accept cookies from your browser. For details, see the list of [forbidden top-level domains](https://dt-url.net/9n6b0pfz).",
			Optional:    true, // nullable
		},
		"same_site_cookie_attribute": {
			Type:        schema.TypeString,
			Description: "Possible Values: `LAX`, `NONE`, `NOTSET`, `STRICT`",
			Required:    true,
		},
		"use_secure_cookie_attribute": {
			Type:        schema.TypeBool,
			Description: "If your application is only accessible via SSL, you can add the Secure attribute to all cookies set by Dynatrace. This setting prevents the display of warnings from PCI-compliance security scanners. Be aware that with this setting enabled Dynatrace correlation of user actions with server-side web requests is only possible over SSL connections.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":              me.ApplicationID,
		"cookie_placement_domain":     me.CookiePlacementDomain,
		"same_site_cookie_attribute":  me.SameSiteCookieAttribute,
		"use_secure_cookie_attribute": me.UseSecureCookieAttribute,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":              &me.ApplicationID,
		"cookie_placement_domain":     &me.CookiePlacementDomain,
		"same_site_cookie_attribute":  &me.SameSiteCookieAttribute,
		"use_secure_cookie_attribute": &me.UseSecureCookieAttribute,
	})
}
