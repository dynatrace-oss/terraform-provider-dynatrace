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

package attribute

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (me *Settings) Deprecated() string {
	return "This resource API endpoint has been deprecated, please use `dynatrace_attribute_allow_list` and `dynatrace_attribute_masking` instead."
}

type Settings struct {
	Key     string      `json:"key"`     // Key of the span attribute to store
	Masking MaskingType `json:"masking"` // Possible Values: `MASK_ENTIRE_VALUE`, `MASK_ONLY_CONFIDENTIAL_DATA`, `NOT_MASKED`
}

func (me *Settings) Name() string {
	return me.Key
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "Key of the span attribute to store",
			Required:    true,
		},
		"masking": {
			Type:        schema.TypeString,
			Description: "Possible Values: `MASK_ENTIRE_VALUE`, `MASK_ONLY_CONFIDENTIAL_DATA`, `NOT_MASKED`",
			Required:    true,
		},
		"persistent": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Prevents the Span Attribute from getting deleted when running `terraform destroy` - to be used for Span Attributes that are defined by default on every Dynatrace environment.",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":     me.Key,
		"masking": string(me.Masking),
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if key, ok := decoder.GetOk("key"); ok {
		me.Key = key.(string)
	}
	if value, ok := decoder.GetOk("masking"); ok {
		me.Masking = MaskingType(value.(string))
	}
	return nil
}
