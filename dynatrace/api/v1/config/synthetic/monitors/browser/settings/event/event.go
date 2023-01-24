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

package event

import (
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Events []Event

type EventWrapper struct {
	Event Event
}

func (me *EventWrapper) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the event to appear in the UI",
			Required:    true,
		},
		"select": {
			Type:        schema.TypeList,
			Description: "Properties specified for a key strokes event. ",
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(SelectOption).Schema()},
		},
		"navigate": {
			Type:        schema.TypeList,
			Description: "Properties specified for a navigation event",
			Optional:    true,
			MaxItems:    1,
			// ConflictsWith: []string{"select"},
			Elem: &schema.Resource{Schema: new(Navigate).Schema()},
		},
		"keystrokes": {
			Type:        schema.TypeList,
			Description: "Properties specified for a key strokes event",
			Optional:    true,
			MaxItems:    1,
			// ConflictsWith: []string{"navigate", "select"},
			Elem: &schema.Resource{Schema: new(KeyStrokes).Schema()},
		},
		"javascript": {
			Type:        schema.TypeList,
			Description: "Properties specified for a javascript event",
			Optional:    true,
			MaxItems:    1,
			// ConflictsWith: []string{"keystrokes", "navigate", "select"},
			Elem: &schema.Resource{Schema: new(Javascript).Schema()},
		},
		"cookie": {
			Type:        schema.TypeList,
			Description: "Properties specified for a cookie event",
			Optional:    true,
			MaxItems:    1,
			// ConflictsWith: []string{"javascript", "keystrokes", "navigate", "select"},
			Elem: &schema.Resource{Schema: new(Cookie).Schema()},
		},
		"tap": {
			Type:        schema.TypeList,
			Description: "Properties specified for a tap event",
			Optional:    true,
			MaxItems:    1,
			// ConflictsWith: []string{"cookie", "javascript", "keystrokes", "navigate", "select"},
			Elem: &schema.Resource{Schema: new(Tap).Schema()},
		},
		"click": {
			Type:        schema.TypeList,
			Description: "Properties specified for a click event",
			Optional:    true,
			MaxItems:    1,
			// ConflictsWith: []string{"tap", "cookie", "javascript", "keystrokes", "navigate", "select"},
			Elem: &schema.Resource{Schema: new(Click).Schema()},
		},
	}
}

func (me *EventWrapper) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("description", me.Event.GetDescription()); err != nil {
		return err
	}
	marshalled := hcl.Properties{}
	if err := me.Event.MarshalHCL(marshalled); err == nil {
		if me.Event.GetType() == Types.Click {
			properties["click"] = []any{marshalled}
		} else if me.Event.GetType() == Types.Tap {
			properties["tap"] = []any{marshalled}
		} else if me.Event.GetType() == Types.Cookie {
			properties["cookie"] = []any{marshalled}
		} else if me.Event.GetType() == Types.Javascript {
			properties["javascript"] = []any{marshalled}
		} else if me.Event.GetType() == Types.KeyStrokes {
			properties["keystrokes"] = []any{marshalled}
		} else if me.Event.GetType() == Types.Navigate {
			properties["navigate"] = []any{marshalled}
		} else if me.Event.GetType() == Types.SelectOption {
			properties["select"] = []any{marshalled}
		} else {
			return fmt.Errorf("events of type %s are not supported", me.Event.GetType())
		}
	}
	return nil
}

func (me *EventWrapper) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("click.#"); ok && result.(int) != 0 {
		evt := new(Click)
		if err := evt.UnmarshalHCL(hcl.NewDecoder(decoder, "click", 0)); err != nil {
			return err
		}
		evt.Type = Types.Click
		me.Event = evt
	}
	if result, ok := decoder.GetOk("tap.#"); ok && result.(int) != 0 {
		evt := new(Tap)
		if err := evt.UnmarshalHCL(hcl.NewDecoder(decoder, "tap", 0)); err != nil {
			return err
		}
		evt.Type = Types.Tap
		me.Event = evt
	}
	if result, ok := decoder.GetOk("cookie.#"); ok && result.(int) != 0 {
		evt := new(Cookie)
		if err := evt.UnmarshalHCL(hcl.NewDecoder(decoder, "cookie", 0)); err != nil {
			return err
		}
		evt.Type = Types.Cookie
		me.Event = evt
	}
	if result, ok := decoder.GetOk("javascript.#"); ok && result.(int) != 0 {
		evt := new(Javascript)
		if err := evt.UnmarshalHCL(hcl.NewDecoder(decoder, "javascript", 0)); err != nil {
			return err
		}
		evt.Type = Types.Javascript
		me.Event = evt
	}
	if result, ok := decoder.GetOk("keystrokes.#"); ok && result.(int) != 0 {
		evt := new(KeyStrokes)
		if err := evt.UnmarshalHCL(hcl.NewDecoder(decoder, "keystrokes", 0)); err != nil {
			return err
		}
		evt.Type = Types.KeyStrokes
		me.Event = evt
	}
	if result, ok := decoder.GetOk("navigate.#"); ok && result.(int) != 0 {
		evt := new(Navigate)
		if err := evt.UnmarshalHCL(hcl.NewDecoder(decoder, "navigate", 0)); err != nil {
			return err
		}
		evt.Type = Types.Navigate
		me.Event = evt
	}
	if result, ok := decoder.GetOk("select.#"); ok && result.(int) != 0 {
		evt := new(SelectOption)
		if err := evt.UnmarshalHCL(hcl.NewDecoder(decoder, "select", 0)); err != nil {
			return err
		}
		evt.Type = Types.SelectOption
		me.Event = evt
	}
	if me.Event != nil {
		if v, ok := decoder.GetOk("description"); ok {
			me.Event.SetDescription(v.(string))
		}
	}
	return nil
}

func (me *Events) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event": {
			Type:        schema.TypeList,
			Description: "An event",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(EventWrapper).Schema()},
		},
	}
}

func (me *Events) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("event.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			evtw := new(EventWrapper)
			if err := evtw.UnmarshalHCL(hcl.NewDecoder(decoder, "event", idx)); err != nil {
				return err
			}
			*me = append(*me, evtw.Event)
		}
	}
	return nil
}

func (me Events) MarshalHCL(properties hcl.Properties) error {
	entries := []any{}
	for _, event := range me {
		evtw := &EventWrapper{Event: event}
		marshalled := hcl.Properties{}
		if err := evtw.MarshalHCL(marshalled); err == nil {
			entries = append(entries, marshalled)
		} else {
			return err
		}
	}
	properties["event"] = entries
	return nil
}

type evt struct {
	Type Type `json:"type"`
}

func (me *Events) UnmarshalJSON(data []byte) error {
	records := []json.RawMessage{}
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}
	for _, record := range records {
		var e evt
		if err := json.Unmarshal(record, &e); err != nil {
			return err
		}
		if e.Type == Types.Click {
			var re Click
			if err := json.Unmarshal(record, &re); err != nil {
				return err
			}
			*me = append(*me, &re)
		} else if e.Type == Types.Cookie {
			var re Cookie
			if err := json.Unmarshal(record, &re); err != nil {
				return err
			}
			*me = append(*me, &re)
		} else if e.Type == Types.Javascript {
			var re Javascript
			if err := json.Unmarshal(record, &re); err != nil {
				return err
			}
			*me = append(*me, &re)
		} else if e.Type == Types.KeyStrokes {
			var re KeyStrokes
			if err := json.Unmarshal(record, &re); err != nil {
				return err
			}
			*me = append(*me, &re)
		} else if e.Type == Types.Navigate {
			var re Navigate
			if err := json.Unmarshal(record, &re); err != nil {
				return err
			}
			*me = append(*me, &re)
		} else if e.Type == Types.Tap {
			var re Tap
			if err := json.Unmarshal(record, &re); err != nil {
				return err
			}
			*me = append(*me, &re)
		}
	}

	return nil
}

// func (me Events) MarshalJSON() ([]byte, error) {
// }

type Event interface {
	GetType() Type
	GetDescription() string
	SetDescription(string)
	MarshalHCL(hcl.Properties) error
}

type EventBase struct {
	Type        Type   `json:"type"`        // The type of synthetic event
	Description string `json:"description"` // A short description of the event to appear in the UI
}

func (me *EventBase) GetType() Type {
	return me.Type
}

func (me *EventBase) GetDescription() string {
	return me.Description
}

func (me *EventBase) SetDescription(description string) {
	me.Description = description
}

func (me *EventBase) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the event to appear in the UI",
			Required:    true,
		},
	}
}
