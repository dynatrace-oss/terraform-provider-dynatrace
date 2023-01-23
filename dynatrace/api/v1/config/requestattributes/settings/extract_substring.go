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

// ExtractSubstring Preprocess by extracting a substring from the original value.
type ExtractSubstring struct {
	EndDelimiter *string                    `json:"endDelimiter,omitempty"` // The end-delimiter string.   Required if the **position** value is `BETWEEN`. Otherwise not allowed.
	Position     Position                   `json:"position"`               // The position of the extracted string relative to delimiters.
	Delimiter    string                     `json:"delimiter"`              // The delimiter string.
	Unknowns     map[string]json.RawMessage `json:"-"`
}

func (me *ExtractSubstring) IsZero() bool {
	if me.EndDelimiter != nil && len(*me.EndDelimiter) > 0 {
		return false
	}
	if len(me.Position) > 0 {
		return false
	}
	if len(me.Delimiter) > 0 {
		return false
	}
	if len(me.Unknowns) > 0 {
		return false
	}
	return true
}

func (me *ExtractSubstring) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"end_delimiter": {
			Type:        schema.TypeString,
			Description: "The end-delimiter string.   Required if the **position** value is `BETWEEN`. Otherwise not allowed",
			Optional:    true,
		},
		"position": {
			Type:        schema.TypeString,
			Description: "The position of the extracted string relative to delimiters",
			Required:    true,
		},
		"delimiter": {
			Type:        schema.TypeString,
			Description: "The delimiter string",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *ExtractSubstring) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("end_delimiter", me.EndDelimiter); err != nil {
		return err
	}
	if err := properties.Encode("position", string(me.Position)); err != nil {
		return err
	}
	if err := properties.Encode("delimiter", me.Delimiter); err != nil {
		return err
	}
	return nil
}

func (me *ExtractSubstring) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "end_delimiter")
		delete(me.Unknowns, "position")
		delete(me.Unknowns, "delimiter")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("end_delimiter"); ok {
		me.EndDelimiter = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("position"); ok {
		me.Position = Position(value.(string))
	}
	if value, ok := decoder.GetOk("delimiter"); ok {
		me.Delimiter = value.(string)
	}
	return nil
}

func (me *ExtractSubstring) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("endDelimiter", me.EndDelimiter); err != nil {
		return nil, err
	}
	if err := m.Marshal("position", me.Position); err != nil {
		return nil, err
	}
	if err := m.Marshal("delimiter", me.Delimiter); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *ExtractSubstring) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("endDelimiter", &me.EndDelimiter); err != nil {
		return err
	}
	if err := m.Unmarshal("position", &me.Position); err != nil {
		return err
	}
	if err := m.Unmarshal("delimiter", &me.Delimiter); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
