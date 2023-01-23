package alerting

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProfileSeverityRule A severity rule of the alerting profile.
//
//	A severity rule defines the level of severity that must be met before an alert is sent our for a detected problem. Additionally it restricts the alerting to certain monitored entities.
type ProfileSeverityRule struct {
	TagFilter      *ProfileTagFilter          `json:"tagFilter"`      // Configuration of the tag filtering of the alerting profile.
	DelayInMinutes int32                      `json:"delayInMinutes"` // Send a notification if a problem remains open longer than *X* minutes.
	SeverityLevel  SeverityLevel              `json:"severityLevel"`  // The severity level to trigger the alert.
	Unknowns       map[string]json.RawMessage `json:"-"`
}

func (me *ProfileSeverityRule) EnsurePredictableOrder() {
	if me.TagFilter != nil {
		me.TagFilter.EnsurePredictableOrder()
	}
}

func (me *ProfileSeverityRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"severity_level": {
			Type:        schema.TypeString,
			Description: "The severity level to trigger the alert. Possible values are `AVAILABILITY`,	`CUSTOM_ALERT`,	`ERROR`,`MONITORING_UNAVAILABLE`,`PERFORMANCE` and `RESOURCE_CONTENTION`.",
			Required:    true,
		},
		"delay_in_minutes": {
			Type:        schema.TypeInt,
			Description: "Send a notification if a problem remains open longer than *X* minutes",
			Required:    true,
		},
		"tag_filter": {
			Type:        schema.TypeList,
			Description: "Configuration of the tag filtering of the alerting profile",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ProfileTagFilter).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *ProfileSeverityRule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"tag_filter":       me.TagFilter,
		"delay_in_minutes": int(me.DelayInMinutes),
		"severity_level":   string(me.SeverityLevel),
	})
}

func (me *ProfileSeverityRule) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "severity_level")
		delete(me.Unknowns, "delay_in_minutes")
		delete(me.Unknowns, "tag_filter")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("severity_level"); ok {
		me.SeverityLevel = SeverityLevel(value.(string))
	}
	if value, ok := decoder.GetOk("delay_in_minutes"); ok {
		me.DelayInMinutes = int32(value.(int))
	}
	if _, ok := decoder.GetOk("tag_filter.#"); ok {
		me.TagFilter = new(ProfileTagFilter)
		if err := me.TagFilter.UnmarshalHCL(hcl.NewDecoder(decoder, "tag_filter", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *ProfileSeverityRule) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.TagFilter != nil {
		rawMessage, err := json.Marshal(me.TagFilter)
		if err != nil {
			return nil, err
		}
		m["tagFilter"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.DelayInMinutes)
		if err != nil {
			return nil, err
		}
		m["delayInMinutes"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.SeverityLevel)
		if err != nil {
			return nil, err
		}
		m["severityLevel"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *ProfileSeverityRule) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["tagFilter"]; found {
		if err := json.Unmarshal(v, &me.TagFilter); err != nil {
			return err
		}
	}
	if v, found := m["delayInMinutes"]; found {
		if err := json.Unmarshal(v, &me.DelayInMinutes); err != nil {
			return err
		}
	}
	if v, found := m["severityLevel"]; found {
		if err := json.Unmarshal(v, &me.SeverityLevel); err != nil {
			return err
		}
	}

	delete(m, "tagFilter")
	delete(m, "delayInMinutes")
	delete(m, "severityLevel")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// SeverityLevel The severity level to trigger the alert.
type SeverityLevel string

// SeverityLevels offers the known enum values
var SeverityLevels = struct {
	Availability          SeverityLevel
	CustomAlert           SeverityLevel
	Error                 SeverityLevel
	MonitoringUnavailable SeverityLevel
	Performance           SeverityLevel
	ResourceContention    SeverityLevel
}{
	"AVAILABILITY",
	"CUSTOM_ALERT",
	"ERROR",
	"MONITORING_UNAVAILABLE",
	"PERFORMANCE",
	"RESOURCE_CONTENTION",
}
