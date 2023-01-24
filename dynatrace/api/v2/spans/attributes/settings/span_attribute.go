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

package attributes

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SpanAttribute has no documentation
type SpanAttribute struct {
	Key     string      `json:"key"`
	Masking MaskingType `json:"masking"`
}

func (me *SpanAttribute) Name() string {
	return me.Key
}

func (me *SpanAttribute) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "the key of the attribute to capture",
		},
		"masking": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "granular control over the visibility of attribute values",
		},
		"persistent": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Prevents the Span Attribute from getting deleted when running `terraform destroy` - to be used for Span Attributes that are defined by default on every Dynatrace environment.",
		},
	}
}

func (me *SpanAttribute) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":     me.Key,
		"masking": string(me.Masking),
	})
}

func (me *SpanAttribute) UnmarshalHCL(decoder hcl.Decoder) error {
	if key, ok := decoder.GetOk("key"); ok {
		me.Key = key.(string)
	}
	if value, ok := decoder.GetOk("masking"); ok {
		me.Masking = MaskingType(value.(string))
	}
	return nil
}

type MaskingType string

var MaskingTypes = struct {
	NotMasked    MaskingType
	Confidential MaskingType
	EntireValue  MaskingType
}{
	MaskingType("NOT_MASKED"),
	MaskingType("MASK_ONLY_CONFIDENTIAL_DATA"),
	MaskingType("MASK_ENTIRE_VALUE"),
}
