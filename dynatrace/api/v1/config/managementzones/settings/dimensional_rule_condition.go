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

package managementzones

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DimensionalRuleCondition struct {
	Type     ConditionType              `json:"conditionType"`   // The type of the condition
	Match    RuleMatcher                `json:"ruleMatcher"`     // How we compare the values
	Key      string                     `json:"key"`             // The reference value for comparison. For conditions of the `DIMENSION` type, specify the key here
	Value    *string                    `json:"value,omitempty"` // The value of the dimension. Only applicable when the **conditionType** is set to `DIMENSION`
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *DimensionalRuleCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the condition. Possible values are \n   - `DIMENSION`\n   - `LOG_FILE_NAME`\n   - `METRIC_KEY`",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The reference value for comparison. For conditions of the `DIMENSION` type, specify the key here",
			Required:    true,
		},
		"match": {
			Type:        schema.TypeString,
			Description: "How to compare. Possible values are \n   - `BEGINS_WITH`\n   - `EQUALS`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the dimension. Only applicable when type is set to `DIMENSION`",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (me *DimensionalRuleCondition) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("key", string(me.Key)); err != nil {
		return err
	}
	if err := properties.Encode("match", string(me.Match)); err != nil {
		return err
	}
	if me.Value != nil && len(*me.Value) > 0 {
		if err := properties.Encode("value", me.Value); err != nil {
			return err
		}
	}

	return nil
}

func (me *DimensionalRuleCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "type")
		delete(me.Unknowns, "key")
		delete(me.Unknowns, "match")
		delete(me.Unknowns, "value")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = ConditionType(value.(string))
	}
	if value, ok := decoder.GetOk("key"); ok {
		me.Key = value.(string)
	}
	if value, ok := decoder.GetOk("match"); ok {
		me.Match = RuleMatcher(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = opt.NewString(value.(string))
	}

	return nil
}

func (me *DimensionalRuleCondition) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.Type)
		if err != nil {
			return nil, err
		}
		m["conditionType"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Key)
		if err != nil {
			return nil, err
		}
		m["key"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Match)
		if err != nil {
			return nil, err
		}
		m["ruleMatcher"] = rawMessage
	}
	if me.Value != nil {
		rawMessage, err := json.Marshal(me.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *DimensionalRuleCondition) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["conditionType"]; found {
		if err := json.Unmarshal(v, &me.Type); err != nil {
			return err
		}
	}
	if v, found := m["key"]; found {
		if err := json.Unmarshal(v, &me.Key); err != nil {
			return err
		}
	}
	if v, found := m["ruleMatcher"]; found {
		if err := json.Unmarshal(v, &me.Match); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &me.Value); err != nil {
			return err
		}
	}
	delete(m, "conditionType")
	delete(m, "key")
	delete(m, "ruleMatcher")
	delete(m, "value")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// ConditionType The value to compare to.
type ConditionType string

func (v ConditionType) Ref() *ConditionType {
	return &v
}

// ConditionTypes offers the known enum values
var ConditionTypes = struct {
	Dimension   ConditionType
	LogFileName ConditionType
	MetricKey   ConditionType
}{
	"DIMENSION",
	"LOG_FILE_NAME",
	"METRIC_KEY",
}

func (v *ConditionType) String() string {
	return string(*v)
}

// RuleMatcher The value to compare to.
type RuleMatcher string

func (v RuleMatcher) Ref() *RuleMatcher {
	return &v
}

// RuleMatchers offers the known enum values
var RuleMatchers = struct {
	BeginsWith RuleMatcher
	Equals     RuleMatcher
}{
	"BEGINS_WITH",
	"EQUALS",
}

func (v *RuleMatcher) String() string {
	return string(*v)
}
