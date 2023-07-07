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

package monitoring

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventComplexTypes []*EventComplexType

func (me *EventComplexTypes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_pattern": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(EventComplexType).Schema()},
		},
	}
}

func (me EventComplexTypes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("event_pattern", me)
}

func (me *EventComplexTypes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("event_pattern", me)
}

type EventComplexType struct {
	Active  bool   `json:"active"`  // Activate
	Label   string `json:"label"`   // Field selector name
	Pattern string `json:"pattern"` // The set of allowed characters for this field has been extended with ActiveGate version 1.259. For more details, see the [documentation](https://dt-url.net/7h23wuk#set-up-event-field-selectors).
}

func (me *EventComplexType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"active": {
			Type:        schema.TypeBool,
			Description: "Activate",
			Required:    true,
		},
		"label": {
			Type:        schema.TypeString,
			Description: "Field selector name",
			Required:    true,
		},
		"pattern": {
			Type:        schema.TypeString,
			Description: "The set of allowed characters for this field has been extended with ActiveGate version 1.259. For more details, see the [documentation](https://dt-url.net/7h23wuk#set-up-event-field-selectors).",
			Required:    true,
		},
	}
}

func (me *EventComplexType) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"active":  me.Active,
		"label":   me.Label,
		"pattern": me.Pattern,
	})
}

func (me *EventComplexType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"active":  &me.Active,
		"label":   &me.Label,
		"pattern": &me.Pattern,
	})
}
