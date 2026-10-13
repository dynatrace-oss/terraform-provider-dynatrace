/**
* @license
* Copyright 2026 Dynatrace LLC
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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MaskingRules []*MaskingRule

func (me *MaskingRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"masking_rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(MaskingRule).Schema()},
		},
	}
}

func (me MaskingRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("masking_rule", me)
}

func (me *MaskingRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("masking_rule", me)
}

type MaskingRule struct {
	DataType        string `json:"dataType"`        // Data type
	MaskedOutput    string `json:"maskedOutput"`    // Masked output
	MaskingStrategy string `json:"maskingStrategy"` // Masking strategy
}

func (me *MaskingRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data_type": {
			Type:        schema.TypeString,
			Description: "Data type",
			Required:    true,
		},
		"masked_output": {
			Type:        schema.TypeString,
			Description: "Masked output",
			Required:    true,
		},
		"masking_strategy": {
			Type:        schema.TypeString,
			Description: "Masking strategy",
			Required:    true,
		},
	}
}

func (me *MaskingRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"data_type":        me.DataType,
		"masked_output":    me.MaskedOutput,
		"masking_strategy": me.MaskingStrategy,
	})
}

func (me *MaskingRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"data_type":        &me.DataType,
		"masked_output":    &me.MaskedOutput,
		"masking_strategy": &me.MaskingStrategy,
	})
}
