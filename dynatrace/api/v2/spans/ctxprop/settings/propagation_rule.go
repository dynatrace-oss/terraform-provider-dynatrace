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

package ctxprop

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/spans/match"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PropagationRule has no documentation
type PropagationRule struct {
	Name     string             `json:"ruleName"`
	Action   PropagationAction  `json:"ruleAction"`
	Matchers match.SpanMatchers `json:"matchers" min:"1" max:"100"`
}

func (me *PropagationRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the rule",
		},
		"action": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Whether to create an entry point or not",
		},
		"matches": {
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Required:    true,
			Description: "Matching strategies for the Span",
			Elem:        &schema.Resource{Schema: new(match.SpanMatchers).Schema()},
		},
	}
}

func (me *PropagationRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"action":  me.Action,
		"matches": me.Matchers,
	})
}

func (me *PropagationRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"action":  &me.Action,
		"matches": &me.Matchers,
	})
}
