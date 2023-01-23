package alerting

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Profile Configuration of an alerting profile.
type Profile struct {
	ID               *string                    `json:"id,omitempty"`               // The ID of the alerting profile.
	DisplayName      string                     `json:"displayName"`                // The name of the alerting profile, displayed in the UI.
	MzID             *string                    `json:"mzId,omitempty"`             // The ID of the management zone to which the alerting profile applies.
	Rules            []*ProfileSeverityRule     `json:"rules,omitempty"`            // A list of severity rules.   The rules are evaluated from top to bottom. The first matching rule applies and further evaluation stops.  If you specify both severity rule and event filter, the AND logic applies.
	EventTypeFilters []*EventTypeFilter         `json:"eventTypeFilters,omitempty"` // The list of event filters.  For all filters that are *negated* inside of these event filters, that is all "Predefined" as well as "Custom" (Title and/or Description) ones the AND logic applies. For all *non-negated* ones the OR logic applies. Between these two groups, negated and non-negated, the AND logic applies.  If you specify both severity rule and event filter, the AND logic applies.
	Metadata         *ConfigMetadata            `json:"metadata,omitempty"`         // Metadata useful for debugging
	Unknowns         map[string]json.RawMessage `json:"-"`
}

func (me *Profile) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Description: "The name of the alerting profile, displayed in the UI",
			Required:    true,
		},
		"mz_id": {
			Type:        schema.TypeString,
			Description: "The ID of the management zone to which the alerting profile applies",
			Optional:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "A list of rules for management zone usage.  Each rule is evaluated independently of all other rules",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ProfileSeverityRule).Schema()},
		},
		"event_type_filters": {
			Type:        schema.TypeList,
			Description: "The list of event filters.  For all filters that are *negated* inside of these event filters, that is all `Predefined` as well as `Custom` (Title and/or Description) ones the AND logic applies. For all *non-negated* ones the OR logic applies. Between these two groups, negated and non-negated, the AND logic applies.  If you specify both severity rule and event filter, the AND logic applies",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(EventTypeFilter).Schema()},
		},
		"metadata": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Description: "`metadata` exists for backwards compatibility but shouldn't get specified anymore",
			Deprecated:  "`metadata` exists for backwards compatibility but shouldn't get specified anymore",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(ConfigMetadata).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Profile) EnsurePredictableOrder() {
	if len(me.Rules) > 0 {
		conds := []*ProfileSeverityRule{}
		condStrings := sort.StringSlice{}
		for _, entry := range me.Rules {
			entry.EnsurePredictableOrder()
			condBytes, _ := json.Marshal(entry)
			condStrings = append(condStrings, string(condBytes))
		}
		condStrings.Sort()
		for _, condString := range condStrings {
			cond := ProfileSeverityRule{}
			json.Unmarshal([]byte(condString), &cond)
			conds = append(conds, &cond)
		}
		me.Rules = conds
	}
}

func (me *Profile) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}

	me.EnsurePredictableOrder()

	if me.EventTypeFilters != nil {
		filters := append([]*EventTypeFilter{}, me.EventTypeFilters...)
		sort.Slice(filters, func(i, j int) bool {
			d1, _ := json.Marshal(filters[i])
			d2, _ := json.Marshal(filters[j])
			cmp := strings.Compare(string(d1), string(d2))
			return (cmp == -1)
		})
		me.EventTypeFilters = filters
	}

	return properties.EncodeAll(map[string]any{
		"display_name":       me.DisplayName,
		"mz_id":              me.MzID,
		"rules":              me.Rules,
		"event_type_filters": me.EventTypeFilters,
	})

}

func (me *Profile) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "display_name")
		delete(me.Unknowns, "mz_id")
		delete(me.Unknowns, "rules")
		delete(me.Unknowns, "event_type_filters")
		delete(me.Unknowns, "metadata")
		delete(me.Unknowns, "managementZoneId")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("display_name"); ok {
		me.DisplayName = value.(string)
	}
	if value, ok := decoder.GetOk("mz_id"); ok {
		me.MzID = opt.NewString(value.(string))
	}
	if result, ok := decoder.GetOk("rules.#"); ok {
		me.Rules = []*ProfileSeverityRule{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(ProfileSeverityRule)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "rules", idx)); err != nil {
				return err
			}
			me.Rules = append(me.Rules, entry)
		}
	}
	if result, ok := decoder.GetOk("event_type_filters.#"); ok {
		me.EventTypeFilters = []*EventTypeFilter{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(EventTypeFilter)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "event_type_filters", idx)); err != nil {
				return err
			}
			me.EventTypeFilters = append(me.EventTypeFilters, entry)
		}
	}

	return nil
}

func (me *Profile) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.ID != nil {
		rawMessage, err := json.Marshal(me.ID)
		if err != nil {
			return nil, err
		}
		m["id"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.DisplayName)
		if err != nil {
			return nil, err
		}
		m["displayName"] = rawMessage
	}
	if me.MzID != nil {
		rawMessage, err := json.Marshal(me.MzID)
		if err != nil {
			return nil, err
		}
		m["mzId"] = rawMessage
	}
	if len(me.Rules) > 0 {
		rawMessage, err := json.Marshal(me.Rules)
		if err != nil {
			return nil, err
		}
		m["rules"] = rawMessage
	}
	if len(me.EventTypeFilters) > 0 {
		rawMessage, err := json.Marshal(me.EventTypeFilters)
		if err != nil {
			return nil, err
		}
		m["eventTypeFilters"] = rawMessage
	}
	if me.Metadata != nil {
		rawMessage, err := json.Marshal(me.Metadata)
		if err != nil {
			return nil, err
		}
		m["metadata"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *Profile) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["id"]; found {
		if err := json.Unmarshal(v, &me.ID); err != nil {
			return err
		}
	}
	if v, found := m["displayName"]; found {
		if err := json.Unmarshal(v, &me.DisplayName); err != nil {
			return err
		}
	}
	if v, found := m["mzId"]; found {
		if err := json.Unmarshal(v, &me.MzID); err != nil {
			return err
		}
	}
	if v, found := m["rules"]; found {
		if err := json.Unmarshal(v, &me.Rules); err != nil {
			return err
		}
	}
	if v, found := m["eventTypeFilters"]; found {
		if err := json.Unmarshal(v, &me.EventTypeFilters); err != nil {
			return err
		}
	}
	if v, found := m["metadata"]; found {
		if err := json.Unmarshal(v, &me.Metadata); err != nil {
			return err
		}
	}
	delete(m, "id")
	delete(m, "displayName")
	delete(m, "mzId")
	delete(m, "managementZoneId")
	delete(m, "rules")
	delete(m, "eventTypeFilters")
	delete(m, "metadata")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// ConfigMetadata Metadata useful for debugging
type ConfigMetadata struct {
	ClusterVersion               *string  `json:"clusterVersion,omitempty"`               // Dynatrace server version.
	ConfigurationVersions        []int64  `json:"configurationVersions,omitempty"`        // A Sorted list of the version numbers of the configuration.
	CurrentConfigurationVersions []string `json:"currentConfigurationVersions,omitempty"` // A Sorted list of string version numbers of the configuration.
}

func (me *ConfigMetadata) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_version": {
			Type:        schema.TypeString,
			Description: "Dynatrace server version",
			Optional:    true,
		},
		"configuration_versions": {
			Type:        schema.TypeList,
			Description: "A Sorted list of the version numbers of the configuration",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"current_configuration_versions": {
			Type:        schema.TypeList,
			Description: "A Sorted list of the version numbers of the configuration",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *ConfigMetadata) MarshalHCL(properties hcl.Properties) error {
	if me.ClusterVersion != nil && len(*me.ClusterVersion) > 0 {
		if err := properties.Encode("cluster_version", me.ClusterVersion); err != nil {
			return err
		}
	}
	if err := properties.Encode("configuration_versions", me.ConfigurationVersions); err != nil {
		return err
	}
	if err := properties.Encode("current_configuration_versions", me.CurrentConfigurationVersions); err != nil {
		return err
	}
	return nil
}

func (me *ConfigMetadata) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("cluster_version"); ok {
		me.ClusterVersion = opt.NewString(value.(string))
	}
	if _, ok := decoder.GetOk("configuration_versions.#"); ok {
		me.ConfigurationVersions = []int64{}
		if entries, ok := decoder.GetOk("configuration_versions"); ok {
			for _, entry := range entries.([]any) {
				me.ConfigurationVersions = append(me.ConfigurationVersions, int64(entry.(int)))
			}
		}
	}
	if _, ok := decoder.GetOk("current_configuration_versions.#"); ok {
		me.CurrentConfigurationVersions = []string{}
		if entries, ok := decoder.GetOk("current_configuration_versions"); ok {
			for _, entry := range entries.([]any) {
				me.CurrentConfigurationVersions = append(me.CurrentConfigurationVersions, entry.(string))
			}
		}
	}
	return nil
}
