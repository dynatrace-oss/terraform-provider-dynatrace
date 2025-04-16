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

type FilterSegments []*FilterSegment

func (me *FilterSegments) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter_segment": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(FilterSegment).Schema()},
		},
	}
}

func (me FilterSegments) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter_segment", me)
}

func (me *FilterSegments) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter_segment", me)
}

type FilterSegment struct {
	Id        string                 `json:"id"`
	Variables FilterSegmentVariables `json:"value,omitempty"`
}

func (me *FilterSegment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of the filter segment",
			Required:    true,
		},
		"variables": {
			Type:        schema.TypeList,
			Description: "Defines a variable with a name and a list of values",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(FilterSegmentVariables).Schema()},
		},
	}
}

func (me *FilterSegment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":        me.Id,
		"variables": me.Variables,
	})
}

func (me *FilterSegment) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":        &me.Id,
		"variables": &me.Variables,
	})
}
