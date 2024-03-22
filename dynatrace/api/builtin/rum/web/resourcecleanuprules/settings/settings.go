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

package resourcecleanuprules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name              string `json:"name"`              // For example: *Mask journeyId*
	RegularExpression string `json:"regularExpression"` // For example: `(.*)(journeyId=)-?\\d+(.*)`
	ReplaceWith       string `json:"replaceWith"`       // For example: `$1$2\\*$3`
	InsertAfter       string `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "For example: *Mask journeyId*",
			Required:    true,
		},
		"regular_expression": {
			Type:        schema.TypeString,
			Description: "For example: `(.*)(journeyId=)-?\\d+(.*)`",
			Required:    true,
		},
		"replace_with": {
			Type:        schema.TypeString,
			Description: "For example: `$1$2\\*$3`",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":               me.Name,
		"regular_expression": me.RegularExpression,
		"replace_with":       me.ReplaceWith,
		"insert_after":       me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":               &me.Name,
		"regular_expression": &me.RegularExpression,
		"replace_with":       &me.ReplaceWith,
		"insert_after":       &me.InsertAfter,
	})
}
