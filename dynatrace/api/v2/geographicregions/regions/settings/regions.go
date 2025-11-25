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

package regions

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Regions []*Region

func (me *Regions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: new(Region).Schema()},
		},
	}
}

func (me Regions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("region", me)
}

func (me *Regions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("region", me)
}

type Region struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (me *Region) Schema() map[string]*schema.Schema {
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

func (me *Region) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name": me.Name,
		"code": me.Code,
	})
}

func (me *Region) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name": &me.Name,
		"code": &me.Code,
	})
}
