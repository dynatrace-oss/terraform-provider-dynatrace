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

package slo

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FilterSegmentVariables []*FilterSegmentVariable

func (me *FilterSegmentVariables) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter_segment_variable": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(FilterSegmentVariable).Schema()},
		},
	}
}

func (me FilterSegmentVariables) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter_segment_variable", me)
}

func (me *FilterSegmentVariables) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter_segment_variable", me)
}

type FilterSegmentVariable struct {
	Name   string   `json:"name"`
	Values []string `json:"values" maxlength:"1000"`
}

func (me *FilterSegmentVariable) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the filter segment variable",
			Required:    true,
		},
		"values": {
			Type:        schema.TypeSet,
			Description: "Values of the filter segment variable",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *FilterSegmentVariable) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":   me.Name,
		"values": me.Values,
	})
}

func (me *FilterSegmentVariable) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":   &me.Name,
		"values": &me.Values,
	})
}
