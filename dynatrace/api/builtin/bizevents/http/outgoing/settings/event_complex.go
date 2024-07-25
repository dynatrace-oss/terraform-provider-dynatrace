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

package outgoing

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventComplex struct {
	Category *EventCategoryAttributeComplex `json:"category"`       // Event category
	Data     EventDataFieldComplexes        `json:"data,omitempty"` // Additional attributes for the business event.
	Provider *EventAttributeComplex         `json:"provider"`       // Event provider
	Type     *EventAttributeComplex         `json:"type"`           // Event type
}

func (me *EventComplex) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"category": {
			Type:        schema.TypeList,
			Description: "Event category",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventCategoryAttributeComplex).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"data": {
			Type:        schema.TypeList,
			Description: "Additional attributes for the business event.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(EventDataFieldComplexes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"provider": {
			Type:        schema.TypeList,
			Description: "Event provider",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventAttributeComplex).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeList,
			Description: "Event type",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventAttributeComplex).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *EventComplex) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"category": me.Category,
		"data":     me.Data,
		"provider": me.Provider,
		"type":     me.Type,
	})
}

func (me *EventComplex) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"category": &me.Category,
		"data":     &me.Data,
		"provider": &me.Provider,
		"type":     &me.Type,
	})
}
