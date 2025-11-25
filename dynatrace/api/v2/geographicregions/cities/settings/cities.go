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

package cities

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Cities []*City

func (me *Cities) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"city": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: new(City).Schema()},
		},
	}
}

func (me Cities) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("city", me)
}

func (me *Cities) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("city", me)
}

type City struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (me *City) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"latitude": {
			Type:     schema.TypeFloat,
			Required: true,
		},
		"longitude": {
			Type:     schema.TypeFloat,
			Required: true,
		},
	}
}

func (me *City) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":      me.Name,
		"latitude":  me.Latitude,
		"longitude": me.Longitude,
	})
}

func (me *City) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":      &me.Name,
		"latitude":  &me.Latitude,
		"longitude": &me.Longitude,
	})
}
