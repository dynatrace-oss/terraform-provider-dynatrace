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

package privacy

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Masking struct {
	PersonalDataUriMaskingEnabled bool `json:"personalDataUriMaskingEnabled"` // Dynatrace captures the URIs and request headers sent from desktop and mobile browsers. Dynatrace also captures full URIs on the server-side to enable detailed performance analysis of your applications. For complete details, visit [Mask personal data in URIs](https://dt-url.net/mask-personal-data-in-URIs).. URIs and request headers contain personal data. When this setting is enabled, Dynatrace automatically detects UUIDs, credit card numbers, email addresses, IP addresses, and other IDs and replaces those values with placeholders. The personal data is then masked in PurePath analysis, error analysis, user action naming for RUM, and elsewhere in Dynatrace.
	UserActionMaskingEnabled      bool `json:"userActionMaskingEnabled"`      // When Dynatrace detects a user action that triggers a page load or an AJAX/XHR action. To learn more about masking user actions, visit [Mask user actions](https://dt-url.net/mask-user-action).. When Dynatrace detects a user action that triggers a page load or an AJAX/XHR action, it constructs a name for the user action based on:\n\n- User event type (click on..., loading of page..., or keypress on...)\n- Title, caption, label, value, ID, className, or other available property of the related HTML element (for example, an image, button, checkbox, or text input field).\n\nIn most instances, the default approach to user-action naming works well, resulting in user-action names such as:\n\n- click on \"Search\" on page /search.html\n- keypress on \"Feedback\" on page /contact.html\n- touch on \"Homescreen\" of page /list.jsf\n\nIn rare circumstances, confidential data (for example, email addresses, usernames, or account numbers) can be unintentionally included in user action names because the confidential data itself is included in an HTML element label, attribute, or other value (for example, click on \"my Account Number: 1231231\"...). If such confidential data appears in your application's user action names, enable the Mask user action names setting. This setting replaces specific HTML element names and values with generic HTML element names. With user-action name masking enabled, the user action names listed above appear as:\n\n- click on INPUT on page /search.html\n- keypress on TEXTAREA on page /contact.html\n- touch on DIV of page /list.jsf
}

func (me *Masking) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_address_masking": {
			Type:        schema.TypeString,
			Description: "Possible Values: `All`, `Public`",
			Optional:    true,
			Deprecated:  "This property is not supported anymore by the Dynatrace REST API (since schema version 4)",
		},
		"ip_address_masking_enabled": {
			Type:        schema.TypeBool,
			Description: "Dynatrace captures the IP addresses of your end-users to determine the regions from which they access your application. To learn more, visit [Mask IPs and GPS coordinates](https://dt-url.net/mask-end-users-ip-addresses).. Dynatrace also captures GPS data from mobile apps that provide their users with the option of sharing geolocation data. On the server side, Dynatrace captures IP addresses to enable detailed troubleshooting for Dynatrace service calls.\n\nOnce enabled, IP address masking sets the last octet of monitored IPv4 addresses and the last 80 bits of IPv6 addresses to zeroes. GPS coordinates are rounded up to 1 decimal place (~10 km). This masking occurs in memory. Full IP addresses are never written to disk. Location lookups are made using anonymized IP addresses and GPS coordinates.",
			Optional:    true,
			Deprecated:  "This property is not supported anymore by the Dynatrace REST API (since schema version 4)",
		},
		"personal_data_uri_masking_enabled": {
			Type:        schema.TypeBool,
			Description: "Dynatrace captures the URIs and request headers sent from desktop and mobile browsers. Dynatrace also captures full URIs on the server-side to enable detailed performance analysis of your applications. For complete details, visit [Mask personal data in URIs](https://dt-url.net/mask-personal-data-in-URIs).. URIs and request headers contain personal data. When this setting is enabled, Dynatrace automatically detects UUIDs, credit card numbers, email addresses, IP addresses, and other IDs and replaces those values with placeholders. The personal data is then masked in PurePath analysis, error analysis, user action naming for RUM, and elsewhere in Dynatrace.",
			Required:    true,
		},
		"user_action_masking_enabled": {
			Type:        schema.TypeBool,
			Description: "When Dynatrace detects a user action that triggers a page load or an AJAX/XHR action. To learn more about masking user actions, visit [Mask user actions](https://dt-url.net/mask-user-action).. When Dynatrace detects a user action that triggers a page load or an AJAX/XHR action, it constructs a name for the user action based on:\n\n- User event type (click on..., loading of page..., or keypress on...)\n- Title, caption, label, value, ID, className, or other available property of the related HTML element (for example, an image, button, checkbox, or text input field).\n\nIn most instances, the default approach to user-action naming works well, resulting in user-action names such as:\n\n- click on \"Search\" on page /search.html\n- keypress on \"Feedback\" on page /contact.html\n- touch on \"Homescreen\" of page /list.jsf\n\nIn rare circumstances, confidential data (for example, email addresses, usernames, or account numbers) can be unintentionally included in user action names because the confidential data itself is included in an HTML element label, attribute, or other value (for example, click on \"my Account Number: 1231231\"...). If such confidential data appears in your application's user action names, enable the Mask user action names setting. This setting replaces specific HTML element names and values with generic HTML element names. With user-action name masking enabled, the user action names listed above appear as:\n\n- click on INPUT on page /search.html\n- keypress on TEXTAREA on page /contact.html\n- touch on DIV of page /list.jsf",
			Required:    true,
		},
	}
}

func (me *Masking) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"personal_data_uri_masking_enabled": me.PersonalDataUriMaskingEnabled,
		"user_action_masking_enabled":       me.UserActionMaskingEnabled,
	})
}

func (me *Masking) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"personal_data_uri_masking_enabled": &me.PersonalDataUriMaskingEnabled,
		"user_action_masking_enabled":       &me.UserActionMaskingEnabled,
	})
}
