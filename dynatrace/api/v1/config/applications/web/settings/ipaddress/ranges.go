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

package ipaddress

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Ranges []*Range

func (me *Ranges) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"range": {
			Type:        schema.TypeList,
			Description: "The IP address or the IP address range to be mapped to the location",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Range).Schema()},
		},
	}
}

func (me Ranges) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("range", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *Ranges) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("range", me)
}

// Range The IP address or the IP address range to be mapped to the location
type Range struct {
	SubNetMask *int32  `json:"subnetMask,omitempty"` // The subnet mask of the IP address range. Valid values range from 0 to 128.
	Address    string  `json:"address"`              // The IP address to be mapped. \n\nFor an IP address range, this is the **from** address.
	ToAddress  *string `json:"addressTo,omitempty"`  // The **to** address of the IP address range.
}

func (me *Range) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_mask": {
			Type:        schema.TypeInt,
			Description: "The subnet mask of the IP address range. Valid values range from 0 to 128.",
			Optional:    true,
		},
		"address": {
			Type:        schema.TypeString,
			Description: "The IP address to be mapped. \n\nFor an IP address range, this is the **from** address.",
			Required:    true,
		},
		"address_to": {
			Type:        schema.TypeString,
			Description: "The **to** address of the IP address range.",
			Optional:    true,
		},
	}
}

func (me *Range) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"subnet_mask": me.SubNetMask,
		"address":     me.Address,
		"address_to":  me.ToAddress,
	})
}

func (me *Range) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"subnet_mask": &me.SubNetMask,
		"address":     &me.Address,
		"address_to":  &me.ToAddress,
	})
}
