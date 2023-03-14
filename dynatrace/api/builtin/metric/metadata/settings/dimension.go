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

package metadata

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Dimensions []*Dimension

func (me *Dimensions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dimension": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Dimension).Schema()},
		},
	}
}

func (me Dimensions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("dimension", me)
}

func (me *Dimensions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("dimension", me)
}

type Dimension struct {
	DisplayName *string `json:"displayName,omitempty"` // Display name
	Key         string  `json:"key"`                   // Dimension key
}

func (me *Dimension) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name",
			Optional:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Dimension key",
			Required:    true,
		},
	}
}

func (me *Dimension) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"display_name": me.DisplayName,
		"key":          me.Key,
	})
}

func (me *Dimension) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"display_name": &me.DisplayName,
		"key":          &me.Key,
	})
}
