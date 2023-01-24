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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings/restriction"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RestrictionSettings Settings for restricting certain ip addresses and for introducing subnet mask. It also restricts the mode
type RestrictionSettings struct {
	Mode         restriction.Mode `json:"mode"`                            // The mode of the list of ip address restrictions. Possible values area `EXCLUDE` and `INCLUDE`.
	Restrictions Ranges           `json:"ipAddressRestrictions,omitempty"` // The IP addresses or the IP address ranges to be mapped to the location
}

func (me *RestrictionSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mode": {
			Type:        schema.TypeString,
			Description: "The mode of the list of ip address restrictions. Possible values area `EXCLUDE` and `INCLUDE`.",
			Required:    true,
		},
		"restrictions": {
			Type:        schema.TypeList,
			Description: "The IP addresses or the IP address ranges to be mapped to the location",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Ranges).Schema()},
		},
	}
}

func (me *RestrictionSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"mode":         me.Mode,
		"restrictions": me.Restrictions,
	})
}

func (me *RestrictionSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"mode":         &me.Mode,
		"restrictions": &me.Restrictions,
	})
}
