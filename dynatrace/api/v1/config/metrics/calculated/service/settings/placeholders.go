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

package service

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Placeholders []*Placeholder

func (me *Placeholders) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"placeholder": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A custom placeholder to be used in a dimension value pattern",
			Elem:        &schema.Resource{Schema: new(Placeholder).Schema()},
		},
	}
}

func (me Placeholders) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("placeholder", me)
}

func (me *Placeholders) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("placeholder", me)
}
