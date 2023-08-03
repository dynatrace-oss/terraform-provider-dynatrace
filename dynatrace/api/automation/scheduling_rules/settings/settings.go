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

package scheduling_rules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Title              string              `json:"title" minlength:"1" maxlength:"2"`
	Description        *string             `json:"description,omitempty"`
	RuleType           RuleType            `json:"ruleType"`
	RecurrenceRule     *RecurrenceRule     `json:"rrule,omitempty"`
	GroupingRule       *GroupingRule       `json:"groupingRule,omitempty"`
	FixedOffsetRule    *FixedOffsetRule    `json:"fixedOffsetRule,omitempty"`
	RelativeOffsetRule *RelativeOffsetRule `json:"relativeOffsetRule,omitempty"`
	BusinessCalendar   *string             `json:"businessCalendar,omitempty" format:"uuid"`
}

func (me *Settings) Name() string {
	return me.Title
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:             schema.TypeString,
			Description:      "The title / name of the scheduling rule",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(200),
		},
		"description": {
			Type:        schema.TypeString,
			Description: "An optional description for the scheduling rule",
			Optional:    true,
		},
		"business_calendar": {
			Type:             schema.TypeString,
			Description:      "",
			Optional:         true,
			ValidateDiagFunc: ValidateUUID,
		},
		"recurrence": {
			Type:         schema.TypeList,
			Description:  "",
			MinItems:     1,
			MaxItems:     1,
			Optional:     true,
			Elem:         &schema.Resource{Schema: new(RecurrenceRule).Schema()},
			ExactlyOneOf: []string{"recurrence", "grouping", "fixed_offset", "relative_offset"},
		},
		"grouping": {
			Type:         schema.TypeList,
			Description:  "",
			MinItems:     1,
			MaxItems:     1,
			Optional:     true,
			Elem:         &schema.Resource{Schema: new(GroupingRule).Schema()},
			ExactlyOneOf: []string{"recurrence", "grouping", "fixed_offset", "relative_offset"},
		},
		"fixed_offset": {
			Type:         schema.TypeList,
			Description:  "",
			MinItems:     1,
			MaxItems:     1,
			Optional:     true,
			Elem:         &schema.Resource{Schema: new(FixedOffsetRule).Schema()},
			ExactlyOneOf: []string{"recurrence", "grouping", "fixed_offset", "relative_offset"},
		},
		"relative_offset": {
			Type:         schema.TypeList,
			Description:  "",
			MinItems:     1,
			MaxItems:     1,
			Optional:     true,
			Elem:         &schema.Resource{Schema: new(RelativeOffsetRule).Schema()},
			ExactlyOneOf: []string{"recurrence", "grouping", "fixed_offset", "relative_offset"},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"title":             me.Title,
		"description":       me.Description,
		"business_calendar": me.BusinessCalendar,
		"recurrence":        me.RecurrenceRule,
		"grouping":          me.GroupingRule,
		"fixed_offset":      me.FixedOffsetRule,
		"relative_offset":   me.RelativeOffsetRule,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"title":             &me.Title,
		"description":       &me.Description,
		"business_calendar": &me.BusinessCalendar,
		"recurrence":        &me.RecurrenceRule,
		"grouping":          &me.GroupingRule,
		"fixed_offset":      &me.FixedOffsetRule,
		"relative_offset":   &me.RelativeOffsetRule,
	}); err != nil {
		return err
	}
	if me.RecurrenceRule != nil {
		me.RuleType = RuleTypes.RRule
	} else if me.RelativeOffsetRule != nil {
		me.RuleType = RuleTypes.RelativeOffset
	} else if me.FixedOffsetRule != nil {
		me.RuleType = RuleTypes.FixedOffset
	} else if me.GroupingRule != nil {
		me.RuleType = RuleTypes.Grouping
	}
	return nil
}
