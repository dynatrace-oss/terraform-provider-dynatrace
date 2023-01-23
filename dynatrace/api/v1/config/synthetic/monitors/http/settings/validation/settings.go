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

package validation

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Settings helps you verify that your HTTP monitor loads the expected content
type Settings struct {
	Rules Rules `json:"rules,omitempty"` // A list of validation rules
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Description: "A list of validation rules",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Rule).Schema()},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("rule", me.Rules); err != nil {
		return err
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("rule.#"); ok {
		me.Rules = Rules{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Rule)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "rule", idx)); err != nil {
				return err
			}
			me.Rules = append(me.Rules, entry)
		}
	}
	return nil
}
