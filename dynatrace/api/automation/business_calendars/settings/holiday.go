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

package business_calendars

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Holidays []*Holiday

func (me *Holidays) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"holiday": {
			Type:        schema.TypeSet,
			Description: "A (unordered) list of holidays valid in this calendar",
			MinItems:    1,
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Holiday).Schema()},
		},
	}
}

func (me Holidays) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("holiday", me)
}

func (me *Holidays) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("holiday", me)
}

type Holiday struct {
	Title string `json:"title"`
	Date  string `json:"date" format:"date"`
}

func (me *Holiday) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:        schema.TypeString,
			Description: "An official name for this holiday",
			Required:    true,
		},
		"date": {
			Type:        schema.TypeString,
			Description: "The date of this holiday: Example `2017-07-04` for July 4th 2017",
			Required:    true,
		},
	}
}

func (me *Holiday) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"title": me.Title,
		"date":  me.Date,
	})
}

func (me *Holiday) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"title": &me.Title,
		"date":  &me.Date,
	})
}
