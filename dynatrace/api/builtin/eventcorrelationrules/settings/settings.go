/**
* @license
* Copyright 2026 Dynatrace LLC
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

package eventcorrelationrules

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CausalityRules       CausalityRules `json:"causalityRules,omitempty"`       // Defines which correlated events cause which other events.
	CorrelationNamespace *string        `json:"correlationNamespace,omitempty"` // Optional namespace for cross-rule correlation. Rules sharing the same namespace produce compatible grouping keys, enabling events matched by different rules to correlate to the same problem.
	DisplayName          string         `json:"displayName"`                    // Name of the rule
	Enabled              bool           `json:"enabled"`                        // This setting is enabled (`true`) or disabled (`false`)
	GroupingFields       []string       `json:"groupingFields"`                 // Fields used to group events. Groups are created based on the values of the grouping fields.
	InsertAfter          *string        `json:"-"`                              // Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched
	MatchingCondition    string         `json:"matchingCondition"`              // Matching condition to check if an event matches the correlation rule.
	Source               *string        `json:"source,omitempty"`               // Name of the source that provides the rule.
	SourceType           *SourceType    `json:"sourceType,omitempty"`           // Where the rule comes from. Extension-provided and built-in rules are maintained by their source. Possible values: `BUILTIN`, `EXTENSION`, `USER`
	TimeWindow           string         `json:"timeWindow"`                     // Other events within this time window are considered for correlation.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"causality_rules": {
			Type:        schema.TypeList,
			Description: "Defines which correlated events cause which other events.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(CausalityRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"correlation_namespace": {
			Type:        schema.TypeString,
			Description: "Optional namespace for cross-rule correlation. Rules sharing the same namespace produce compatible grouping keys, enabling events matched by different rules to correlate to the same problem.",
			Optional:    true, // nullable
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Name of the rule",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"grouping_fields": {
			Type:        schema.TypeSet,
			Description: "Fields used to group events. Groups are created based on the values of the grouping fields.",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Computed:    true,
			Optional:    true,
		},
		"matching_condition": {
			Type:        schema.TypeString,
			Description: "Matching condition to check if an event matches the correlation rule.",
			Required:    true,
		},
		"source": {
			Type:        schema.TypeString,
			Description: "Name of the source that provides the rule.",
			ForceNew:    true,
			Optional:    true, // nullable & precondition
		},
		"source_type": {
			Type:        schema.TypeString,
			Description: "Where the rule comes from. Extension-provided and built-in rules are maintained by their source. Possible values: `BUILTIN`, `EXTENSION`, `USER`",
			ForceNew:    true,
			Optional:    true, // nullable
		},
		"time_window": {
			Type:        schema.TypeString,
			Description: "Other events within this time window are considered for correlation.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"causality_rules":       me.CausalityRules,
		"correlation_namespace": me.CorrelationNamespace,
		"display_name":          me.DisplayName,
		"enabled":               me.Enabled,
		"grouping_fields":       me.GroupingFields,
		"insert_after":          me.InsertAfter,
		"matching_condition":    me.MatchingCondition,
		"source":                me.Source,
		"source_type":           me.SourceType,
		"time_window":           me.TimeWindow,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Source != nil) && (me.SourceType == nil || !slices.Contains([]string{"EXTENSION", "BUILTIN"}, string(*me.SourceType))) {
		return fmt.Errorf("'source' must not be specified unless 'source_type' is one of ['EXTENSION', 'BUILTIN']; got 'source_type'='%v'", opt.ValOrNil(me.SourceType))
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"causality_rules":       &me.CausalityRules,
		"correlation_namespace": &me.CorrelationNamespace,
		"display_name":          &me.DisplayName,
		"enabled":               &me.Enabled,
		"grouping_fields":       &me.GroupingFields,
		"insert_after":          &me.InsertAfter,
		"matching_condition":    &me.MatchingCondition,
		"source":                &me.Source,
		"source_type":           &me.SourceType,
		"time_window":           &me.TimeWindow,
	})
}
