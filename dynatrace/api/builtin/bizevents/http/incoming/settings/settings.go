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

package incoming

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled  bool             `json:"enabled"`         // This setting is enabled (`true`) or disabled (`false`)
	Event    *EventComplex    `json:"event"`           // Event meta data
	RuleName string           `json:"ruleName"`        // Rule name
	Scope    *string          `json:"-" scope:"scope"` // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	Triggers MatcherComplexes `json:"triggers"`        // Define conditions to trigger business events from incoming web requests. Triggers are connected by AND logic per capture rule. If you set multiple trigger rules, all of them need to be fulfilled to capture a business event.
}

func (me *Settings) Name() string {
	return me.RuleName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"event": {
			Type:        schema.TypeList,
			Description: "Event meta data",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventComplex).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"triggers": {
			Type:        schema.TypeList,
			Description: "Define conditions to trigger business events from incoming web requests. Triggers are connected by AND logic per capture rule. If you set multiple trigger rules, all of them need to be fulfilled to capture a business event.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(MatcherComplexes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":   me.Enabled,
		"event":     me.Event,
		"rule_name": me.RuleName,
		"scope":     me.Scope,
		"triggers":  me.Triggers,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":   &me.Enabled,
		"event":     &me.Event,
		"rule_name": &me.RuleName,
		"scope":     &me.Scope,
		"triggers":  &me.Triggers,
	})
}
