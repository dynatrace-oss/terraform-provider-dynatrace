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

package sitereliabilityguardian

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Variables []*Variable

func (me *Variables) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"variable": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Variable).Schema()},
		},
	}
}

func (me Variables) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("variable", me)
}

func (me *Variables) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("variable", me)
}

type Variable struct {
	Definition string `json:"definition"` // Value
	Name       string `json:"name"`
}

func (me *Variable) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"definition": {
			Type:        schema.TypeString,
			Description: "Value",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *Variable) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"definition": me.Definition,
		"name":       me.Name,
	})
}

func (me *Variable) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"definition": &me.Definition,
		"name":       &me.Name,
	})
}
