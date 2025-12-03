/**
* @license
* Copyright 2025 Dynatrace LLC
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

type EventDataFieldComplexes []*EventDataFieldComplex

func (me *EventDataFieldComplexes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_data_field_complex": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(EventDataFieldComplex).Schema()},
		},
	}
}

func (me EventDataFieldComplexes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("event_data_field_complex", me)
}

func (me *EventDataFieldComplexes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("event_data_field_complex", me)
}

type EventDataFieldComplex struct {
	Name   string                     `json:"name"` // Field name to be added to data.
	Source *EventDataAttributeComplex `json:"source"`
}

func (me *EventDataFieldComplex) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Field name to be added to data.",
			Required:    true,
		},
		"source": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventDataAttributeComplex).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *EventDataFieldComplex) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":   me.Name,
		"source": me.Source,
	})
}

func (me *EventDataFieldComplex) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":   &me.Name,
		"source": &me.Source,
	})
}
