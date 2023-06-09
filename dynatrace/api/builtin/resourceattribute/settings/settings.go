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

package resourceattribute

import (
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AttributeKeys RuleItems `json:"attributeKeys,omitempty"` // Attribute key allow-list
}

func (me *Settings) Name() string {
	return "dynatrace_resource_attributes"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"keys": {
			Type:        schema.TypeList,
			Description: "Attribute key allow-list",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(RuleItems).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) EnsurePredictableOrder() {
	if len(me.AttributeKeys) == 0 {
		return
	}
	conds := RuleItems{}
	condStrings := sort.StringSlice{}
	for _, entry := range me.AttributeKeys {
		condBytes, _ := json.Marshal(entry)
		condStrings = append(condStrings, string(condBytes))
	}
	condStrings.Sort()
	for _, condString := range condStrings {
		cond := RuleItem{}
		json.Unmarshal([]byte(condString), &cond)
		conds = append(conds, &cond)
	}
	me.AttributeKeys = conds
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if len(me.AttributeKeys) > 0 {
		me.EnsurePredictableOrder()
		marshalled := hcl.Properties{}
		if err := me.AttributeKeys.MarshalHCL(marshalled); err != nil {
			return err
		}
		properties["keys"] = []any{marshalled}
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("keys.#"); ok {
		me.AttributeKeys = RuleItems{}
		if err := me.AttributeKeys.UnmarshalHCL(hcl.NewDecoder(decoder, "keys", 0)); err != nil {
			return err
		}
	}
	return nil
}
