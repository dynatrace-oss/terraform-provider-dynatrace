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

package fullwebservice

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ReducedTransformations []*ReducedTransformation

func (me *ReducedTransformations) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transformation": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ReducedTransformation).Schema()},
		},
	}
}

func (me ReducedTransformations) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("transformation", me)
}

func (me *ReducedTransformations) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("transformation", me)
}

type ReducedTransformation struct {
	IncludeHexNumbers  *bool                         `json:"includeHexNumbers,omitempty"` // include hexadecimal numbers
	MinDigitCount      *int                          `json:"minDigitCount,omitempty"`     // min digit count
	Prefix             *string                       `json:"prefix,omitempty"`
	ReplacementValue   *string                       `json:"replacementValue,omitempty"` // replacement
	Suffix             *string                       `json:"suffix,omitempty"`
	TransformationType ContextRootTransformationType `json:"transformationType"` // Possible Values: `BEFORE`, `REMOVE_CREDIT_CARDS`, `REMOVE_IBANS`, `REMOVE_IPS`, `REMOVE_NUMBERS`, `REPLACE_BETWEEN`
}

func (me *ReducedTransformation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"include_hex_numbers": {
			Type:        schema.TypeBool,
			Description: "include hexadecimal numbers",
			Optional:    true, // precondition
		},
		"min_digit_count": {
			Type:        schema.TypeInt,
			Description: "min digit count",
			Optional:    true, // precondition
		},
		"prefix": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable & precondition
		},
		"replacement_value": {
			Type:        schema.TypeString,
			Description: "replacement",
			Optional:    true, // nullable & precondition
		},
		"suffix": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable & precondition
		},
		"transformation_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BEFORE`, `REMOVE_CREDIT_CARDS`, `REMOVE_IBANS`, `REMOVE_IPS`, `REMOVE_NUMBERS`, `REPLACE_BETWEEN`",
			Required:    true,
		},
	}
}

func (me *ReducedTransformation) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"include_hex_numbers": me.IncludeHexNumbers,
		"min_digit_count":     me.MinDigitCount,
		"prefix":              me.Prefix,
		"replacement_value":   me.ReplacementValue,
		"suffix":              me.Suffix,
		"transformation_type": me.TransformationType,
	})
}

func (me *ReducedTransformation) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"include_hex_numbers": &me.IncludeHexNumbers,
		"min_digit_count":     &me.MinDigitCount,
		"prefix":              &me.Prefix,
		"replacement_value":   &me.ReplacementValue,
		"suffix":              &me.Suffix,
		"transformation_type": &me.TransformationType,
	})
	if me.IncludeHexNumbers == nil && me.TransformationType == ContextRootTransformationTypes.RemoveNumbers {
		me.IncludeHexNumbers = opt.NewBool(false)
	}
	return err
}
