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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DimensionConditions []*DimensionCondition

func (me *DimensionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "Dimension conditions",
			Elem:        &schema.Resource{Schema: new(DimensionCondition).Schema()},
		},
	}
}

func (me DimensionConditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("condition", me)
}

func (me *DimensionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("condition", me)
}

type DimensionCondition struct {
	ConditionType DimensionConditionType `json:"conditionType"` // Possible Values: `DIMENSION`, `LOG_FILE_NAME`, `METRIC_KEY`
	Key           *string                `json:"key,omitempty"`
	RuleMatcher   DimensionOperator      `json:"ruleMatcher"` // Possible Values: `BEGINS_WITH`, `EQUALS`
	Value         string                 `json:"value"`
}

func (me *DimensionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DIMENSION`, `LOG_FILE_NAME`, `METRIC_KEY`",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"rule_matcher": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BEGINS_WITH`, `EQUALS`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *DimensionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition_type": me.ConditionType,
		"key":            me.Key,
		"rule_matcher":   me.RuleMatcher,
		"value":          me.Value,
	})
}

func (me *DimensionCondition) HandlePreconditions() error {
	if me.Key == nil && (string(me.ConditionType) == "DIMENSION") {
		return fmt.Errorf("'key' must be specified if 'condition_type' is set to '%v'", me.ConditionType)
	}
	return nil
}

func (me *DimensionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition_type": &me.ConditionType,
		"key":            &me.Key,
		"rule_matcher":   &me.RuleMatcher,
		"value":          &me.Value,
	})
}
