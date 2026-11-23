/**
* @license
* Copyright 2026 Dynatrace LLC
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

package eventcorrelationrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CausalityRules []*CausalityRule

func (me *CausalityRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"causality_rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CausalityRule).Schema()},
		},
	}
}

func (me CausalityRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("causality_rule", me)
}

func (me *CausalityRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("causality_rule", me)
}

type CausalityRule struct {
	CauseEventMatcher  string `json:"causeEventMatcher"`  // Matching condition selecting the events that cause the effect events.
	EffectEventMatcher string `json:"effectEventMatcher"` // Matching condition selecting the events caused by the cause events.
}

func (me *CausalityRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cause_event_matcher": {
			Type:        schema.TypeString,
			Description: "Matching condition selecting the events that cause the effect events.",
			Required:    true,
		},
		"effect_event_matcher": {
			Type:        schema.TypeString,
			Description: "Matching condition selecting the events caused by the cause events.",
			Required:    true,
		},
	}
}

func (me *CausalityRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cause_event_matcher":  me.CauseEventMatcher,
		"effect_event_matcher": me.EffectEventMatcher,
	})
}

func (me *CausalityRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cause_event_matcher":  &me.CauseEventMatcher,
		"effect_event_matcher": &me.EffectEventMatcher,
	})
}
