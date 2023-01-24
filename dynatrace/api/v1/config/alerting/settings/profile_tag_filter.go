package alerting

import (
	"encoding/json"
	"sort"

	common "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/common"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProfileTagFilter Configuration of the tag filtering of the alerting profile.
type ProfileTagFilter struct {
	IncludeMode IncludeMode                `json:"includeMode"`          // The filtering mode:  * `INCLUDE_ANY`: The rule applies to monitored entities that have at least one of the specified tags. You can specify up to 100 tags.  * `INCLUDE_ALL`: The rule applies to monitored entities that have **all** of the specified tags. You can specify up to 10 tags.  * `NONE`: The rule applies to all monitored entities.
	TagFilters  []*common.TagFilter        `json:"tagFilters,omitempty"` // A list of required tags.
	Unknowns    map[string]json.RawMessage `json:"-"`
}

func (me *ProfileTagFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"include_mode": {
			Type:        schema.TypeString,
			Description: "The filtering mode:  * `INCLUDE_ANY`: The rule applies to monitored entities that have at least one of the specified tags. You can specify up to 100 tags.  * `INCLUDE_ALL`: The rule applies to monitored entities that have **all** of the specified tags. You can specify up to 10 tags.  * `NONE`: The rule applies to all monitored entities",
			Required:    true,
		},
		"tag_filters": {
			Type:        schema.TypeList,
			Description: "A list of required tags",
			Optional:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: new(common.TagFilter).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *ProfileTagFilter) EnsurePredictableOrder() {
	if len(me.TagFilters) > 0 {
		conds := []*common.TagFilter{}
		condStrings := sort.StringSlice{}
		for _, entry := range me.TagFilters {
			condBytes, _ := json.Marshal(entry)
			condStrings = append(condStrings, string(condBytes))
		}
		condStrings.Sort()
		for _, condString := range condStrings {
			cond := common.TagFilter{}
			json.Unmarshal([]byte(condString), &cond)
			conds = append(conds, &cond)
		}
		me.TagFilters = conds
	}
}

func (me *ProfileTagFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"include_mode": string(me.IncludeMode),
		"tag_filters":  me.TagFilters,
	})
}

func (me *ProfileTagFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "include_mode")
		delete(me.Unknowns, "tag_filters")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("include_mode"); ok {
		me.IncludeMode = IncludeMode(value.(string))
	}
	if result, ok := decoder.GetOk("tag_filters.#"); ok {
		me.TagFilters = []*common.TagFilter{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(common.TagFilter)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "tag_filters", idx)); err != nil {
				return err
			}
			me.TagFilters = append(me.TagFilters, entry)
		}
	}
	return nil
}

func (me *ProfileTagFilter) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.IncludeMode)
		if err != nil {
			return nil, err
		}
		m["includeMode"] = rawMessage
	}
	if len(me.TagFilters) > 0 {
		rawMessage, err := json.Marshal(me.TagFilters)
		if err != nil {
			return nil, err
		}
		m["tagFilters"] = rawMessage
	}

	return json.Marshal(m)
}

func (me *ProfileTagFilter) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["includeMode"]; found {
		if err := json.Unmarshal(v, &me.IncludeMode); err != nil {
			return err
		}
	}
	if v, found := m["tagFilters"]; found {
		if err := json.Unmarshal(v, &me.TagFilters); err != nil {
			return err
		}
	}

	delete(m, "includeMode")
	delete(m, "tagFilters")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// IncludeMode The filtering mode:
// * `INCLUDE_ANY`: The rule applies to monitored entities that have at least one of the specified tags. You can specify up to 100 tags.
// * `INCLUDE_ALL`: The rule applies to monitored entities that have **all** of the specified tags. You can specify up to 10 tags.
// * `NONE`: The rule applies to all monitored entities.
type IncludeMode string

// IncludeModes offers the known enum values
var IncludeModes = struct {
	IncludeAll IncludeMode
	IncludeAny IncludeMode
	None       IncludeMode
}{
	"INCLUDE_ALL",
	"INCLUDE_ANY",
	"NONE",
}
