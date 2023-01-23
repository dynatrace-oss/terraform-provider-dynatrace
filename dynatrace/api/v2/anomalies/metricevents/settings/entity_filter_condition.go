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

package metricevents

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EntityFilterConditions []*EntityFilterCondition // Entity filter conditions

func (me *EntityFilterConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Entity filter conditions",
			Elem:        &schema.Resource{Schema: new(EntityFilterCondition).Schema()},
		},
	}
}

func (me EntityFilterConditions) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("condition", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *EntityFilterConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("condition"); ok {

		entrySet := value.(*schema.Set)

		for _, entryMap := range entrySet.List() {
			hash := entrySet.F(entryMap)
			entry := new(EntityFilterCondition)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "condition", hash)); err != nil {
				return err
			}
			*me = append(*me, entry)
		}
	}
	return nil
}

type EntityFilterCondition struct {
	Type     EntityFilterType     `json:"type"`
	Operator EntityFilterOperator `json:"operator"`
	Value    string               `json:"value"`
}

func (me *EntityFilterCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "",
			Required:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "",
			Required:    true,
		},
	}
}

func (me *EntityFilterCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":     me.Type,
		"operator": me.Operator,
		"value":    me.Value,
	})
}

func (me *EntityFilterCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":     &me.Type,
		"operator": &me.Operator,
		"value":    &me.Value,
	})
}
