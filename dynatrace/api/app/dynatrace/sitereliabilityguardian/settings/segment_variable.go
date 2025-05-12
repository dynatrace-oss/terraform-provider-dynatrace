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

type SegmentVariables []*SegmentVariable

func (me *SegmentVariables) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"variable": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SegmentVariable).Schema()},
		},
	}
}

func (me SegmentVariables) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("variable", me)
}

func (me *SegmentVariables) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("variable", me)
}

type SegmentVariable struct {
	Name   string   `json:"name"`             // Variable Name
	Values []string `json:"values,omitempty"` // Variable Values
}

func (me *SegmentVariable) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Variable Name",
			Required:    true,
		},
		"values": {
			Type:        schema.TypeList,
			Description: "Variable Values",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *SegmentVariable) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":   me.Name,
		"values": me.Values,
	})
}

func (me *SegmentVariable) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":   &me.Name,
		"values": &me.Values,
	})
}
