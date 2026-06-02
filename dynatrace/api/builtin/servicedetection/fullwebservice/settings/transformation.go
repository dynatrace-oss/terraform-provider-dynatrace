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
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Transformations []*Transformation

func (me *Transformations) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transformation": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Transformation).Schema()},
		},
	}
}

func (me Transformations) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("transformation", me)
}

func (me *Transformations) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("transformation", me)
}

// Transformation. Transforms a detected contributor value before it contributes to the Service Id.
type Transformation struct {
	IncludeHexNumbers  *bool              `json:"includeHexNumbers,omitempty"` // Whether to also remove hexadecimal numbers (sequences of at least `minDigitCount` hexadecimal digits preceded by '0x'). It is used only when the transformation type is `REMOVE_NUMBERS`.
	MinDigitCount      *int               `json:"minDigitCount,omitempty"`     // The minimum number of digits that a numeric sequence must have to be removed. It is used only when the transformation type is `REMOVE_NUMBERS`.
	Prefix             *string            `json:"prefix,omitempty"`            // The part of the text that serves as a reference point for the transformation. Its use depends on the transformation type.
	ReplacementValue   *string            `json:"replacementValue,omitempty"`  // The text that replaces the part between `prefix` and `suffix`. It is used only when the transformation type is `REPLACE_BETWEEN`.
	SegmentCount       *int               `json:"segmentCount,omitempty"`      // How many segments should be taken.
	SelectIndex        *int               `json:"selectIndex,omitempty"`       // The index of the element to keep after splitting. The index is zero-based. It is used only when the transformation type is `SPLIT_SELECT`.
	SplitDelimiter     *string            `json:"splitDelimiter,omitempty"`    // The delimiter used for splitting the text. It is used only when the transformation type is `SPLIT_SELECT` or `TAKE_SEGMENTS`.
	Suffix             *string            `json:"suffix,omitempty"`            // The part of the text that serves as a reference point for the transformation. Its use depends on the transformation type.
	TakeFromEnd        *bool              `json:"takeFromEnd,omitempty"`       // Whether to take segments from the end of the text instead of the beginning. It is used only when the transformation type is `TAKE_SEGMENTS`.
	TransformationType TransformationType `json:"transformationType"`          // Defines what kind of transformation will be applied on the original value. Possible values: `AFTER`, `BEFORE`, `BETWEEN`, `REMOVE_CREDIT_CARDS`, `REMOVE_IBANS`, `REMOVE_IPS`, `REMOVE_NUMBERS`, `REPLACE_BETWEEN`, `SPLIT_SELECT`, `TAKE_SEGMENTS`
}

func (me *Transformation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"include_hex_numbers": {
			Type:        schema.TypeBool,
			Description: "Whether to also remove hexadecimal numbers (sequences of at least `minDigitCount` hexadecimal digits preceded by '0x'). It is used only when the transformation type is `REMOVE_NUMBERS`.",
			Optional:    true, // precondition
		},
		"min_digit_count": {
			Type:        schema.TypeInt,
			Description: "The minimum number of digits that a numeric sequence must have to be removed. It is used only when the transformation type is `REMOVE_NUMBERS`.",
			Optional:    true, // precondition
		},
		"prefix": {
			Type:        schema.TypeString,
			Description: "The part of the text that serves as a reference point for the transformation. Its use depends on the transformation type.",
			Optional:    true, // nullable & precondition
		},
		"replacement_value": {
			Type:        schema.TypeString,
			Description: "The text that replaces the part between `prefix` and `suffix`. It is used only when the transformation type is `REPLACE_BETWEEN`.",
			Optional:    true, // nullable & precondition
		},
		"segment_count": {
			Type:        schema.TypeInt,
			Description: "How many segments should be taken.",
			Optional:    true, // precondition
		},
		"select_index": {
			Type:        schema.TypeInt,
			Description: "The index of the element to keep after splitting. The index is zero-based. It is used only when the transformation type is `SPLIT_SELECT`.",
			Optional:    true, // precondition
		},
		"split_delimiter": {
			Type:        schema.TypeString,
			Description: "The delimiter used for splitting the text. It is used only when the transformation type is `SPLIT_SELECT` or `TAKE_SEGMENTS`.",
			Optional:    true, // nullable & precondition
		},
		"suffix": {
			Type:        schema.TypeString,
			Description: "The part of the text that serves as a reference point for the transformation. Its use depends on the transformation type.",
			Optional:    true, // nullable & precondition
		},
		"take_from_end": {
			Type:        schema.TypeBool,
			Description: "Whether to take segments from the end of the text instead of the beginning. It is used only when the transformation type is `TAKE_SEGMENTS`.",
			Optional:    true, // precondition
		},
		"transformation_type": {
			Type:        schema.TypeString,
			Description: "Defines what kind of transformation will be applied on the original value. Possible values: `AFTER`, `BEFORE`, `BETWEEN`, `REMOVE_CREDIT_CARDS`, `REMOVE_IBANS`, `REMOVE_IPS`, `REMOVE_NUMBERS`, `REPLACE_BETWEEN`, `SPLIT_SELECT`, `TAKE_SEGMENTS`",
			Required:    true,
		},
	}
}

func (me *Transformation) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"include_hex_numbers": me.IncludeHexNumbers,
		"min_digit_count":     me.MinDigitCount,
		"prefix":              me.Prefix,
		"replacement_value":   me.ReplacementValue,
		"segment_count":       me.SegmentCount,
		"select_index":        me.SelectIndex,
		"split_delimiter":     me.SplitDelimiter,
		"suffix":              me.Suffix,
		"take_from_end":       me.TakeFromEnd,
		"transformation_type": me.TransformationType,
	})
}

func (me *Transformation) HandlePreconditions() error {
	if (me.IncludeHexNumbers == nil) && (string(me.TransformationType) == "REMOVE_NUMBERS") {
		me.IncludeHexNumbers = new(false)
	}
	if (me.TakeFromEnd == nil) && (string(me.TransformationType) == "TAKE_SEGMENTS") {
		me.TakeFromEnd = new(false)
	}
	if (me.IncludeHexNumbers != nil) && (string(me.TransformationType) != "REMOVE_NUMBERS") {
		return fmt.Errorf("'include_hex_numbers' must not be specified unless 'transformation_type' is set to 'REMOVE_NUMBERS'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.MinDigitCount != nil) && (string(me.TransformationType) != "REMOVE_NUMBERS") {
		return fmt.Errorf("'min_digit_count' must not be specified unless 'transformation_type' is set to 'REMOVE_NUMBERS'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.MinDigitCount == nil) && (string(me.TransformationType) == "REMOVE_NUMBERS") {
		return fmt.Errorf("'min_digit_count' must be specified when 'transformation_type' is set to 'REMOVE_NUMBERS'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.Prefix != nil) && (!slices.Contains([]string{"AFTER", "BETWEEN", "REPLACE_BETWEEN"}, string(me.TransformationType))) {
		return fmt.Errorf("'prefix' must not be specified unless 'transformation_type' is one of ['AFTER', 'BETWEEN', 'REPLACE_BETWEEN']; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.ReplacementValue != nil) && (string(me.TransformationType) != "REPLACE_BETWEEN") {
		return fmt.Errorf("'replacement_value' must not be specified unless 'transformation_type' is set to 'REPLACE_BETWEEN'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.SegmentCount != nil) && (string(me.TransformationType) != "TAKE_SEGMENTS") {
		return fmt.Errorf("'segment_count' must not be specified unless 'transformation_type' is set to 'TAKE_SEGMENTS'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.SegmentCount == nil) && (string(me.TransformationType) == "TAKE_SEGMENTS") {
		return fmt.Errorf("'segment_count' must be specified when 'transformation_type' is set to 'TAKE_SEGMENTS'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.SelectIndex != nil) && (string(me.TransformationType) != "SPLIT_SELECT") {
		return fmt.Errorf("'select_index' must not be specified unless 'transformation_type' is set to 'SPLIT_SELECT'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.SelectIndex == nil) && (string(me.TransformationType) == "SPLIT_SELECT") {
		return fmt.Errorf("'select_index' must be specified when 'transformation_type' is set to 'SPLIT_SELECT'; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.SplitDelimiter != nil) && (!slices.Contains([]string{"SPLIT_SELECT", "TAKE_SEGMENTS"}, string(me.TransformationType))) {
		return fmt.Errorf("'split_delimiter' must not be specified unless 'transformation_type' is one of ['SPLIT_SELECT', 'TAKE_SEGMENTS']; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.Suffix != nil) && (!slices.Contains([]string{"BEFORE", "BETWEEN", "REPLACE_BETWEEN"}, string(me.TransformationType))) {
		return fmt.Errorf("'suffix' must not be specified unless 'transformation_type' is one of ['BEFORE', 'BETWEEN', 'REPLACE_BETWEEN']; got 'transformation_type'='%v'", me.TransformationType)
	}
	if (me.TakeFromEnd != nil) && (string(me.TransformationType) != "TAKE_SEGMENTS") {
		return fmt.Errorf("'take_from_end' must not be specified unless 'transformation_type' is set to 'TAKE_SEGMENTS'; got 'transformation_type'='%v'", me.TransformationType)
	}
	return nil
}

func (me *Transformation) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"include_hex_numbers": &me.IncludeHexNumbers,
		"min_digit_count":     &me.MinDigitCount,
		"prefix":              &me.Prefix,
		"replacement_value":   &me.ReplacementValue,
		"segment_count":       &me.SegmentCount,
		"select_index":        &me.SelectIndex,
		"split_delimiter":     &me.SplitDelimiter,
		"suffix":              &me.Suffix,
		"take_from_end":       &me.TakeFromEnd,
		"transformation_type": &me.TransformationType,
	})
}
