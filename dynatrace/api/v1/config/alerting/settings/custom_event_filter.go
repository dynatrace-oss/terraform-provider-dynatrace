package alerting

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomEventFilter Configuration of a custom event filter.
// Filters custom events by title or description. If both specified, the AND logic applies.
type CustomEventFilter struct {
	Description *CustomTextFilter          `json:"customDescriptionFilter,omitempty"` // Configuration of a matching filter.
	Title       *CustomTextFilter          `json:"customTitleFilter,omitempty"`       // Configuration of a matching filter.
	Unknowns    map[string]json.RawMessage `json:"-"`
}

func (me *CustomEventFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_description_filter": {
			Type:        schema.TypeList,
			Description: "Configuration of a matching filter",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomTextFilter).Schema()},
		},
		"custom_title_filter": {
			Type:        schema.TypeList,
			Description: "Configuration of a matching filter",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomTextFilter).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomEventFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("custom_description_filter", me.Description); err != nil {
		return err
	}
	if err := properties.Encode("custom_title_filter", me.Title); err != nil {
		return err
	}
	return nil
}

func (me *CustomEventFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "custom_description_filter")
		delete(me.Unknowns, "custom_title_filter")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if _, ok := decoder.GetOk("custom_description_filter.#"); ok {
		me.Description = new(CustomTextFilter)
		if err := me.Description.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_description_filter", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("custom_title_filter.#"); ok {
		me.Title = new(CustomTextFilter)
		if err := me.Title.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_title_filter", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *CustomEventFilter) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.Description != nil {
		rawMessage, err := json.Marshal(me.Description)
		if err != nil {
			return nil, err
		}
		m["customDescriptionFilter"] = rawMessage
	}
	if me.Title != nil {
		rawMessage, err := json.Marshal(me.Title)
		if err != nil {
			return nil, err
		}
		m["customTitleFilter"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *CustomEventFilter) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["customDescriptionFilter"]; found {
		if err := json.Unmarshal(v, &me.Description); err != nil {
			return err
		}
	}
	if v, found := m["customTitleFilter"]; found {
		if err := json.Unmarshal(v, &me.Title); err != nil {
			return err
		}
	}

	delete(m, "customDescriptionFilter")
	delete(m, "customTitleFilter")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
