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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
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

func (me *ReducedTransformation) HandlePreconditions() error {
	if me.IncludeHexNumbers == nil && (string(me.TransformationType) == "REMOVE_NUMBERS") {
		me.IncludeHexNumbers = opt.NewBool(false)
	}
	if me.MinDigitCount == nil && (string(me.TransformationType) == "REMOVE_NUMBERS") {
		return fmt.Errorf("'min_digit_count' must be specified if 'transformation_type' is set to '%v'", me.TransformationType)
	}
	if me.Prefix == nil && slices.Contains([]string{"REPLACE_BETWEEN"}, string(me.TransformationType)) {
		return fmt.Errorf("'prefix' must be specified if 'transformation_type' is set to '%v'", me.TransformationType)
	}
	if me.ReplacementValue == nil && (string(me.TransformationType) == "REPLACE_BETWEEN") {
		return fmt.Errorf("'replacement_value' must be specified if 'transformation_type' is set to '%v'", me.TransformationType)
	}
	if me.Suffix == nil && slices.Contains([]string{"BEFORE", "REPLACE_BETWEEN"}, string(me.TransformationType)) {
		return fmt.Errorf("'suffix' must be specified if 'transformation_type' is set to '%v'", me.TransformationType)
	}
	return nil
}

func (me *ReducedTransformation) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"include_hex_numbers": &me.IncludeHexNumbers,
		"min_digit_count":     &me.MinDigitCount,
		"prefix":              &me.Prefix,
		"replacement_value":   &me.ReplacementValue,
		"suffix":              &me.Suffix,
		"transformation_type": &me.TransformationType,
	})
}
