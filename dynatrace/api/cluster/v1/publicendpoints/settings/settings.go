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

package publicendpoints

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The configuration of the user
type Settings struct {
	WebUiAddress             *string
	AdditionalWebUiAddresses []string
	BeaconForwarderAddress   *string
	CDNAddress               *string
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"web_ui_address": {
			Type:        schema.TypeString,
			Description: "Web UI address",
			Optional:    true,
		},
		"additional_web_ui_addresses": {
			Type:        schema.TypeSet,
			Description: "Additional web UI addresses",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"beacon_forwarder_address": {
			Type:        schema.TypeString,
			Description: "Beacon forwarder address",
			Optional:    true,
		},
		"cdn_address": {
			Type:        schema.TypeString,
			Description: "CDN address",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"web_ui_address":              me.WebUiAddress,
		"additional_web_ui_addresses": me.AdditionalWebUiAddresses,
		"beacon_forwarder_address":    me.BeaconForwarderAddress,
		"cdn_address":                 me.CDNAddress,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"web_ui_address":              &me.WebUiAddress,
		"additional_web_ui_addresses": &me.AdditionalWebUiAddresses,
		"beacon_forwarder_address":    &me.BeaconForwarderAddress,
		"cdn_address":                 &me.CDNAddress,
	})
}
