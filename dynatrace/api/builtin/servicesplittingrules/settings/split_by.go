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

package servicesplittingrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SplitBies []*SplitBy

func (me *SplitBies) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_splitting_attribute": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SplitBy).Schema()},
		},
	}
}

func (me SplitBies) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("service_splitting_attribute", me)
}

func (me *SplitBies) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("service_splitting_attribute", me)
}

type SplitBy struct {
	Key string `json:"key"` // Attribute key
}

func (me *SplitBy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "Attribute key",
			Required:    true,
		},
	}
}

func (me *SplitBy) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key": me.Key,
	})
}

func (me *SplitBy) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key": &me.Key,
	})
}
