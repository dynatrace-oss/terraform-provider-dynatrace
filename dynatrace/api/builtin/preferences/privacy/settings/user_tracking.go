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

type UserTracking struct {
	PersistentCookieEnabled bool `json:"persistentCookieEnabled"` // When enabled, Dynatrace places a [persistent cookie](https://dt-url.net/313o0p4n) on all end-user devices to identify returning users.
}

func (me *UserTracking) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"persistent_cookie_enabled": {
			Type:        schema.TypeBool,
			Description: "When enabled, Dynatrace places a [persistent cookie](https://dt-url.net/313o0p4n) on all end-user devices to identify returning users.",
			Required:    true,
		},
	}
}

func (me *UserTracking) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"persistent_cookie_enabled": me.PersistentCookieEnabled,
	})
}

func (me *UserTracking) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"persistent_cookie_enabled": &me.PersistentCookieEnabled,
	})
}
