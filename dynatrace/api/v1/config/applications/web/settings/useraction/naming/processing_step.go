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

package naming

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProcessingSteps []*ProcessingStep

func (me *ProcessingSteps) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"step": {
			Type:        schema.TypeList,
			Description: "The processing step",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ProcessingStep).Schema()},
		},
	}
}

func (me ProcessingSteps) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("step", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *ProcessingSteps) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("step", me)
}

type ProcessingStep struct {
	Type                    ProcessingStepType `json:"type"`                              // An action to be taken by the processing: \n\n* `SUBSTRING`: Extracts the string between **patternBefore** and **patternAfter**. \n* `REPLACEMENT`: Replaces the string between **patternBefore** and **patternAfter** with the specified **replacement**.\n* `REPLACE_WITH_PATTERN`: Replaces the **patternToReplace** with the specified **replacement**. \n* `EXTRACT_BY_REGULAR_EXPRESSION`: Extracts the part of the string that matches the **regularExpression**. \n* `REPLACE_WITH_REGULAR_EXPRESSION`: Replaces all occurrences that match **regularExpression** with the specified **replacement**. \n* `REPLACE_IDS`: Replaces all IDs and UUIDs with the specified **replacement**. Possible values are `EXTRACT_BY_REGULAR_EXPRESSION`, `REPLACEMENT`, `REPLACE_IDS`, `REPLACE_WITH_PATTERN`, `REPLACE_WITH_REGULAR_EXPRESSION` and `SUBSTRING`.
	PatternBefore           *string            `json:"patternBefore,omitempty"`           // The pattern before the required value. It will be removed
	PatternBeforeSearchType *PatternSearchType `json:"patternBeforeSearchType,omitempty"` // The required occurrence of **patternBefore**. Possible values are `FIRST` and `LAST`.
	PatternAfter            *string            `json:"patternAfter,omitempty"`            // The pattern after the required value. It will be removed.
	PatternAfterSearchType  *PatternSearchType `json:"patternAfterSearchType,omitempty"`  // The required occurrence of **patternAfter**. Possible values are `FIRST` and `LAST`.
	Replacement             *string            `json:"replacement,omitempty"`             // Replacement for the original value
	PatternToReplace        *string            `json:"patternToReplace,omitempty"`        // The pattern to be replaced. \n\n Only applicable if the `type` is `REPLACE_WITH_PATTERN`.
	RegularExpression       *string            `json:"regularExpression,omitempty"`       // A regular expression for the string to be extracted or replaced. Only applicable if the `type` is `EXTRACT_BY_REGULAR_EXPRESSION` or `REPLACE_WITH_REGULAR_EXPRESSION`.
	FallbackToInput         bool               `json:"fallbackToInput,omitempty"`         // If set to `true`: Returns the input if `patternBefore` or `patternAfter` cannot be found and the `type` is `SUBSTRING`. Returns the input if `regularExpression` doesn't match and `type` is `EXTRACT_BY_REGULAR_EXPRESSION`. \n\n Otherwise `null` is returned.
}

func (me *ProcessingStep) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "An action to be taken by the processing: \n\n* `SUBSTRING`: Extracts the string between `patternBefore` and `patternAfter`. \n* `REPLACEMENT`: Replaces the string between `patternBefore` and `patternAfter` with the specified `replacement`.\n* `REPLACE_WITH_PATTERN`: Replaces the **patternToReplace** with the specified **replacement**. \n* `EXTRACT_BY_REGULAR_EXPRESSION`: Extracts the part of the string that matches the **regularExpression**. \n* `REPLACE_WITH_REGULAR_EXPRESSION`: Replaces all occurrences that match **regularExpression** with the specified **replacement**. \n* `REPLACE_IDS`: Replaces all IDs and UUIDs with the specified **replacement**. Possible values are `EXTRACT_BY_REGULAR_EXPRESSION`, `REPLACEMENT`, `REPLACE_IDS`, `REPLACE_WITH_PATTERN`, `REPLACE_WITH_REGULAR_EXPRESSION` and `SUBSTRING`.",
			Required:    true,
		},
		"pattern_before": {
			Type:        schema.TypeString,
			Description: "The pattern before the required value. It will be removed.",
			Optional:    true,
		},
		"pattern_before_search_type": {
			Type:        schema.TypeString,
			Description: "The required occurrence of **patternBefore**. Possible values are `FIRST` and `LAST`.",
			Optional:    true,
		},
		"pattern_after": {
			Type:        schema.TypeString,
			Description: "The pattern after the required value. It will be removed.",
			Optional:    true,
		},
		"pattern_after_search_type": {
			Type:        schema.TypeString,
			Description: "The required occurrence of **patternAfter**. Possible values are `FIRST` and `LAST`.",
			Optional:    true,
		},
		"replacement": {
			Type:        schema.TypeString,
			Description: "Replacement for the original value",
			Optional:    true,
		},
		"pattern_to_replace": {
			Type:        schema.TypeString,
			Description: "The pattern to be replaced. \n\n Only applicable if the `type` is `REPLACE_WITH_PATTERN`.",
			Optional:    true,
		},
		"regular_expression": {
			Type:        schema.TypeString,
			Description: "A regular expression for the string to be extracted or replaced. Only applicable if the `type` is `EXTRACT_BY_REGULAR_EXPRESSION` or `REPLACE_WITH_REGULAR_EXPRESSION`.",
			Optional:    true,
		},
		"fallback_to_input": {
			Type:        schema.TypeBool,
			Description: "If set to `true`: Returns the input if `patternBefore` or `patternAfter` cannot be found and the `type` is `SUBSTRING`. Returns the input if `regularExpression` doesn't match and `type` is `EXTRACT_BY_REGULAR_EXPRESSION`. \n\n Otherwise `null` is returned.",
			Optional:    true,
		},
	}
}

func (me *ProcessingStep) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":                       me.Type,
		"pattern_before":             me.PatternBefore,
		"pattern_before_search_type": me.PatternBeforeSearchType,
		"pattern_after":              me.PatternAfter,
		"pattern_after_search_type":  me.PatternAfterSearchType,
		"replacement":                me.Replacement,
		"pattern_to_replace":         me.PatternToReplace,
		"regular_expression":         me.RegularExpression,
		"fallback_to_input":          me.FallbackToInput,
	})
}

func (me *ProcessingStep) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":                       &me.Type,
		"pattern_before":             &me.PatternBefore,
		"pattern_before_search_type": &me.PatternBeforeSearchType,
		"pattern_after":              &me.PatternAfter,
		"pattern_after_search_type":  &me.PatternAfterSearchType,
		"replacement":                &me.Replacement,
		"pattern_to_replace":         &me.PatternToReplace,
		"regular_expression":         &me.RegularExpression,
		"fallback_to_input":          &me.FallbackToInput,
	})
}
