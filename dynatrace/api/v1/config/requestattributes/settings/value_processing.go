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

package requestattributes

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValueProcessing Process values as specified.
type ValueProcessing struct {
	ExtractSubstring    *ExtractSubstring          `json:"extractSubstring,omitempty"`    // Preprocess by extracting a substring from the original value.
	SplitAt             *string                    `json:"splitAt,omitempty"`             // Split (preprocessed) string values at this separator.
	Trim                *bool                      `json:"trim"`                          // Prune Whitespaces. Defaults to false.
	ValueCondition      *ValueCondition            `json:"valueCondition,omitempty"`      // IBM integration bus label node name condition for which the value is captured.
	ValueExtractorRegex *string                    `json:"valueExtractorRegex,omitempty"` // Extract value from captured data per regex.
	Unknowns            map[string]json.RawMessage `json:"-"`
}

func (me *ValueProcessing) IsZero() bool {
	if me.ExtractSubstring != nil && !me.ExtractSubstring.IsZero() {
		return false
	}
	if me.SplitAt != nil && len(*me.SplitAt) > 0 {
		return false
	}
	if opt.Bool(me.Trim) {
		return false
	}
	if me.ValueCondition != nil && !me.ValueCondition.IsZero() {
		return false
	}
	if me.ValueExtractorRegex != nil && len(*me.ValueExtractorRegex) > 0 {
		return false
	}
	if len(me.Unknowns) > 0 {
		return false
	}
	return true
}

func (me *ValueProcessing) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"extract_substring": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Preprocess by extracting a substring from the original value",
			Elem: &schema.Resource{
				Schema: new(ExtractSubstring).Schema(),
			},
		},
		"split_at": {
			Type:        schema.TypeString,
			Description: "Split (preprocessed) string values at this separator",
			Optional:    true,
		},
		"trim": {
			Type:        schema.TypeBool,
			Description: "Prune Whitespaces. Defaults to false",
			Optional:    true,
		},
		"value_condition": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "IBM integration bus label node name condition for which the value is captured",
			Elem: &schema.Resource{
				Schema: new(ValueCondition).Schema(),
			},
		},
		"value_extractor_regex": {
			Type:        schema.TypeString,
			Description: "Extract value from captured data per regex",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *ValueProcessing) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("extract_substring", me.ExtractSubstring); err != nil {
		return err
	}
	if err := properties.Encode("split_at", me.SplitAt); err != nil {
		return err
	}
	if err := properties.Encode("trim", opt.Bool(me.Trim)); err != nil {
		return err
	}
	if err := properties.Encode("value_condition", me.ValueCondition); err != nil {
		return err
	}
	if err := properties.Encode("value_extractor_regex", me.ValueExtractorRegex); err != nil {
		return err
	}
	return nil
}

func (me *ValueProcessing) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "extract_substring")
		delete(me.Unknowns, "split_at")
		delete(me.Unknowns, "trim")
		delete(me.Unknowns, "value_condition")
		delete(me.Unknowns, "value_extractor_regex")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if _, ok := decoder.GetOk("extract_substring.#"); ok {
		me.ExtractSubstring = new(ExtractSubstring)
		if err := me.ExtractSubstring.UnmarshalHCL(hcl.NewDecoder(decoder, "extract_substring", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("split_at"); ok {
		me.SplitAt = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("trim"); ok {
		me.Trim = opt.NewBool(value.(bool))
	}
	if _, ok := decoder.GetOk("value_condition.#"); ok {
		me.ValueCondition = new(ValueCondition)
		if err := me.ValueCondition.UnmarshalHCL(hcl.NewDecoder(decoder, "value_condition", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("value_extractor_regex"); ok {
		me.ValueExtractorRegex = opt.NewString(value.(string))
	}
	return nil
}

func (me *ValueProcessing) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("extractSubstring", me.ExtractSubstring); err != nil {
		return nil, err
	}
	if err := m.Marshal("splitAt", me.SplitAt); err != nil {
		return nil, err
	}
	if err := m.Marshal("trim", opt.Bool(me.Trim)); err != nil {
		return nil, err
	}
	if err := m.Marshal("valueCondition", me.ValueCondition); err != nil {
		return nil, err
	}
	if err := m.Marshal("valueExtractorRegex", me.ValueExtractorRegex); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *ValueProcessing) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("extractSubstring", &me.ExtractSubstring); err != nil {
		return err
	}
	if err := m.Unmarshal("splitAt", &me.SplitAt); err != nil {
		return err
	}
	if err := m.Unmarshal("trim", &me.Trim); err != nil {
		return err
	}
	if err := m.Unmarshal("valueCondition", &me.ValueCondition); err != nil {
		return err
	}
	if err := m.Unmarshal("valueExtractorRegex", &me.ValueExtractorRegex); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
