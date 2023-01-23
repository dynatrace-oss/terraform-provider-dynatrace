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

package alerting

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventFilters []*EventFilter

func (me *EventFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A conditions for the metric usage",
			Elem:        &schema.Resource{Schema: new(EventFilter).Schema()},
		},
	}
}

func (me EventFilters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *EventFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter", me)
}

type EventFilter struct {
	Type       EventFilterType        `json:"type"`             // The type of event to filter by
	Predefined *PredefinedEventFilter `json:"predefinedFilter"` // The predefined filter. Only valid if `type` is `PREDEFINED`
	Custom     *CustomEventFilter     `json:"customFilter"`     // The custom filter. Only valid if `type` is `CUSTOM`
}

func (me *EventFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom": {
			Type:        schema.TypeList,
			Description: "Configuration of a custom event filter. Filters custom events by title or description. If both specified, the AND logic applies",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomEventFilter).Schema()},
		},
		"predefined": {
			Type:        schema.TypeList,
			Description: "Configuration of a custom event filter. Filters custom events by title or description. If both specified, the AND logic applies",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(PredefinedEventFilter).Schema()},
		},
	}
}

func (me *EventFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("custom", me.Custom); err != nil {
		return err
	}
	if err := properties.Encode("predefined", me.Predefined); err != nil {
		return err
	}
	return nil
}

func (me *EventFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("custom.#"); ok {
		me.Custom = new(CustomEventFilter)
		if err := me.Custom.UnmarshalHCL(hcl.NewDecoder(decoder, "custom", 0)); err != nil {
			return err
		}
		me.Type = EventFilterTypes.Custom
	}
	if _, ok := decoder.GetOk("predefined.#"); ok {
		me.Predefined = new(PredefinedEventFilter)
		if err := me.Predefined.UnmarshalHCL(hcl.NewDecoder(decoder, "predefined", 0)); err != nil {
			return err
		}
		me.Type = EventFilterTypes.Predefined
	}
	return nil
}
