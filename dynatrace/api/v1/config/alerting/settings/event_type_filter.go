package alerting

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// EventTypeFilter Configuration of the event filter for the alerting profile.
// You have two mutually exclusive options:
// * Select an event type from the list of the predefined events. Specify it in the **predefinedEventFilter** field.
// * Set a rule for custom events. Specify it in the **customEventFilter** field.
type EventTypeFilter struct {
	CustomEventFilter     *CustomEventFilter         `json:"customEventFilter,omitempty"`     // Configuration of a custom event filter.  Filters custom events by title or description. If both specified, the AND logic applies.
	PredefinedEventFilter *PredefinedEventFilter     `json:"predefinedEventFilter,omitempty"` // Configuration of a predefined event filter.
	Unknowns              map[string]json.RawMessage `json:"-"`
}

func (me *EventTypeFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_event_filter": {
			Type:        schema.TypeList,
			Description: "Configuration of a custom event filter. Filters custom events by title or description. If both specified, the AND logic applies",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomEventFilter).Schema()},
		},
		"predefined_event_filter": {
			Type:        schema.TypeList,
			Description: "Configuration of a custom event filter. Filters custom events by title or description. If both specified, the AND logic applies",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(PredefinedEventFilter).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *EventTypeFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"custom_event_filter":     me.CustomEventFilter,
		"predefined_event_filter": me.PredefinedEventFilter,
	})
}

func (me *EventTypeFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "custom_event_filter")
		delete(me.Unknowns, "predefined_event_filter")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if _, ok := decoder.GetOk("custom_event_filter.#"); ok {
		me.CustomEventFilter = new(CustomEventFilter)
		if err := me.CustomEventFilter.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_event_filter", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("predefined_event_filter.#"); ok {
		me.PredefinedEventFilter = new(PredefinedEventFilter)
		if err := me.PredefinedEventFilter.UnmarshalHCL(hcl.NewDecoder(decoder, "predefined_event_filter", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *EventTypeFilter) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.CustomEventFilter != nil {
		rawMessage, err := json.Marshal(me.CustomEventFilter)
		if err != nil {
			return nil, err
		}
		m["customEventFilter"] = rawMessage
	}
	if me.PredefinedEventFilter != nil {
		rawMessage, err := json.Marshal(me.PredefinedEventFilter)
		if err != nil {
			return nil, err
		}
		m["predefinedEventFilter"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *EventTypeFilter) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["customEventFilter"]; found {
		if err := json.Unmarshal(v, &me.CustomEventFilter); err != nil {
			return err
		}
	}
	if v, found := m["predefinedEventFilter"]; found {
		if err := json.Unmarshal(v, &me.PredefinedEventFilter); err != nil {
			return err
		}
	}

	delete(m, "customEventFilter")
	delete(m, "predefinedEventFilter")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
