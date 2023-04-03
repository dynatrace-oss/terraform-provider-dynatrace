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

package traffic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IpAddressForms []*IpAddressForm

func (me *IpAddressForms) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_address_form": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(IpAddressForm).Schema()},
		},
	}
}

func (me IpAddressForms) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("ip_address_form", me)
}

func (me *IpAddressForms) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("ip_address_form", me)
}

type IpAddressForm struct {
	IpAddress string `json:"ipAddress"` // IP address
}

func (me *IpAddressForm) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_address": {
			Type:        schema.TypeString,
			Description: "IP address",
			Required:    true,
		},
	}
}

func (me *IpAddressForm) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"ip_address": me.IpAddress,
	})
}

func (me *IpAddressForm) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"ip_address": &me.IpAddress,
	})
}
