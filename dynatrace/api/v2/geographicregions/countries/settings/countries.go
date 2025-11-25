/**
* @license
* Copyright 2025 Dynatrace LLC
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

package countries

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Countries []*Country

func (me *Countries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"country": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: new(Country).Schema()},
		},
	}
}

func (me Countries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("country", me)
}

func (me *Countries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("country", me)
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (me *Country) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"code": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func (me *Country) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name": me.Name,
		"code": me.Code,
	})
}

func (me *Country) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name": &me.Name,
		"code": &me.Code,
	})
}
