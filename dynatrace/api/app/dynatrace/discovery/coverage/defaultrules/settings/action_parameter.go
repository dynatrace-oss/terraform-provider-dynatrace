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

package defaultrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ActionParameters []*ActionParameter

func (me *ActionParameters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"parameter": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ActionParameter).Schema()},
		},
	}
}

func (me ActionParameters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("parameter", me)
}

func (me *ActionParameters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("parameter", me)
}

type ActionParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (me *ActionParameter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *ActionParameter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"value": me.Value,
	})
}

func (me *ActionParameter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"value": &me.Value,
	})
}
