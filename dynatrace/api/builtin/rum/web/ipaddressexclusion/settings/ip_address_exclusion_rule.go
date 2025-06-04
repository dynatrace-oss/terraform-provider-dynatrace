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

type IpAddressExclusionRules []*IpAddressExclusionRule

func (me *IpAddressExclusionRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_exclusion": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(IpAddressExclusionRule).Schema()},
		},
	}
}

func (me IpAddressExclusionRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("ip_exclusion", me)
}

func (me *IpAddressExclusionRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("ip_exclusion", me)
}

type IpAddressExclusionRule struct {
	Ip   string  `json:"ip"`             // Single IP or IP range start address
	IpTo *string `json:"ipTo,omitempty"` // IP range end
}

func (me *IpAddressExclusionRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip": {
			Type:        schema.TypeString,
			Description: "Single IP or IP range start address",
			Required:    true,
		},
		"ip_to": {
			Type:        schema.TypeString,
			Description: "IP range end",
			Optional:    true, // nullable
		},
	}
}

func (me *IpAddressExclusionRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"ip":    me.Ip,
		"ip_to": me.IpTo,
	})
}

func (me *IpAddressExclusionRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"ip":    &me.Ip,
		"ip_to": &me.IpTo,
	})
}
