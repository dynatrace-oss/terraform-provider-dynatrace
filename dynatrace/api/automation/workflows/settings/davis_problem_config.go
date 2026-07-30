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

package workflows

import (
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DavisProblemConfig struct {
	EntityTagsMatch       *EntityTagsMatch        `json:"entityTagsMatch"`                 // Possible values: `all` and `any`
	EntityTags            map[string]StringArray  `json:"entityTags"`                      // key/value pairs for entity tags to match for. For tags that don't require a value, just specify an empty string as value. Multiple values can be provided separated by whitespace (e.g. \"val1 val2\") and will be parsed as multiple tag values. Omit this attribute if all entities should match
	OnProblemClose        bool                    `json:"onProblemClose" default:"false"`  //
	TriggerOn             string                  `json:"triggerOn,omitempty"`             // Problem state to trigger on. Possible values: `open`, `open-and-close`, `close`. When unset, falls back to `on_problem_close`
	Categories            *DavisProblemCategories `json:"categories"`                      //
	CustomFilter          string                  `json:"customFilter,omitempty"`          //
	AnalysisReady         bool                    `json:"analysisReady" default:"false"`   // If set to `true`, the workflow will only be triggered after the initial root cause analysis run is completed
	SeverityThreshold     *int                    `json:"severityThreshold,omitempty"`     // Triggers only for problems whose severity is this value or more severe. Lower numbers are more severe (1 = critical, 5 = informational)
	TriggerOnUpdateFields []string                `json:"triggerOnUpdateFields,omitempty"` // Problem event fields tracked for value changes. Changes to any selected field cause re-triggering
	ProblemOpenDuration   *int                    `json:"problemOpenDuration,omitempty"`   // Minimum problem duration in minutes before the trigger fires. Possible values: 5, 10, 15, 30, 60, 120, 240, 1440, 10080
}

func (me *DavisProblemConfig) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_tags_match": {
			Type:         schema.TypeString,
			Description:  "Specifies whether all or just any of the configured entity tags need to match. Possible values: `all` and `any`. Omit this attribute if all entities should match",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
			RequiredWith: []string{prefix + ".0.entity_tags"},
		},
		"entity_tags": {
			Type:         schema.TypeMap,
			Description:  "key/value pairs for entity tags to match for. For tags that don't require a value, just specify an empty string as value. Multiple values can be provided separated by whitespace (e.g. \"val1 val2\") and will be parsed as multiple tag values. Omit this attribute if all entities should match",
			Optional:     true,
			Elem:         &schema.Schema{Type: schema.TypeString},
			RequiredWith: []string{prefix + ".0.entity_tags_match"},
		},
		"on_problem_close": {
			Type:        schema.TypeBool,
			Description: "If set to `true` closing a problem also is considered an event that triggers the execution.",
			Deprecated:  "Use `trigger_on` instead",
			Optional:    true,
			Default:     false,
		},
		"trigger_on": {
			Type:        schema.TypeString,
			Description: "Problem state to trigger on. Possible values: `open` (active only), `open-and-close` (both phases), or `close` (closure only). When unset, falls back to `on_problem_close`",
			Optional:    true,
		},
		"categories": {
			Type:        schema.TypeList,
			Description: "",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DavisProblemCategories).Schema(prefix + ".0.categories")},
		},
		"custom_filter": {
			Type:        schema.TypeString,
			Description: "Additional DQL matcher expression to further filter events to match",
			Optional:    true,
		},
		"analysis_ready": {
			Type:        schema.TypeBool,
			Description: "If set to `true`, the workflow will only be triggered after the initial root cause analysis run is completed",
			Optional:    true,
			Default:     false,
		},
		"severity_threshold": {
			Type:        schema.TypeInt,
			Description: "Triggers only for problems whose severity is this value or more severe. Possible values: `1` (critical) to `5` (informational). Lower numbers are more severe, so 3 matches severities 1, 2, and 3",
			Optional:    true,
		},
		"trigger_on_update_fields": {
			Type:        schema.TypeSet,
			Description: "Problem event fields tracked for value changes. Changes to any selected field cause re-triggering. Possible values: `dt.davis.affected_users_count`, `dt.davis.impact_level`, `event.category`, `event.severity`, `root_cause_entity_id`, `smartscape.affected_entities`",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"problem_open_duration": {
			Type:        schema.TypeInt,
			Description: "Minimum problem duration in minutes before the trigger fires. Possible values: `5`, `10`, `15`, `30`, `60`, `120`, `240`, `1440`, `10080`",
			Optional:    true,
		},
	}
}

func (me *DavisProblemConfig) MarshalHCL(properties hcl.Properties) error {
	if err := me.MarshalEntityTagsHCL(properties); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"entity_tags_match":        me.EntityTagsMatch,
		"on_problem_close":         me.OnProblemClose,
		"trigger_on":               me.TriggerOn,
		"categories":               me.Categories,
		"custom_filter":            me.CustomFilter,
		"analysis_ready":           me.AnalysisReady,
		"severity_threshold":       me.SeverityThreshold,
		"trigger_on_update_fields": me.TriggerOnUpdateFields,
		"problem_open_duration":    me.ProblemOpenDuration,
	})
}

func (me *DavisProblemConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := me.UnmarshalEntityTagsHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"entity_tags_match":        &me.EntityTagsMatch,
		"on_problem_close":         &me.OnProblemClose,
		"trigger_on":               &me.TriggerOn,
		"categories":               &me.Categories,
		"custom_filter":            &me.CustomFilter,
		"analysis_ready":           &me.AnalysisReady,
		"severity_threshold":       &me.SeverityThreshold,
		"trigger_on_update_fields": &me.TriggerOnUpdateFields,
		"problem_open_duration":    &me.ProblemOpenDuration,
	})
}

func (me *DavisProblemConfig) MarshalEntityTagsHCL(properties hcl.Properties) error {
	entityTagsMap := map[string]string{}
	for k, v := range me.EntityTags {
		if len(k) == 0 {
			continue
		}
		if len(v) == 0 {
			continue
		}
		entityTagsMap[k] = strings.Join([]string(v), " ")
	}
	if len(entityTagsMap) > 0 {
		if err := properties.Encode("entity_tags", entityTagsMap); err != nil {
			return err
		}
	}
	return nil
}

func (me *DavisProblemConfig) UnmarshalEntityTagsHCL(decoder hcl.Decoder) error {
	entityTagsMap := map[string]string{}
	if err := decoder.Decode("entity_tags", &entityTagsMap); err != nil {
		return err
	}
	for k, v := range entityTagsMap {
		if len(k) == 0 {
			continue
		}
		if me.EntityTags == nil {
			me.EntityTags = map[string]StringArray{}
		}
		parts := strings.Split(v, " ")
		var sa StringArray
		for _, p := range parts {
			p = strings.TrimSpace(p)
			sa = append(sa, p)
		}
		me.EntityTags[k] = sa
	}
	return nil
}
