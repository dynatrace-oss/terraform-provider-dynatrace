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

package defaultrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Rule     *Rule         `json:"rule"`     // Rule:
	Settings *RuleSettings `json:"settings"` // Settings:
}

func (me *Settings) Name() string {
	return me.Rule.ID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Description: "Rule:",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Rule).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"settings": {
			Type:        schema.TypeList,
			Description: "Settings:",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(RuleSettings).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"rule":     me.Rule,
		"settings": me.Settings,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"rule":     &me.Rule,
		"settings": &me.Settings,
	})
}
