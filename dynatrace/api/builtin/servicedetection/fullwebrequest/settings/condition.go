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

package fullwebrequest

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Conditions []*Condition

func (me *Conditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Condition).Schema()},
		},
	}
}

func (me Conditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("condition", me)
}

func (me *Conditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("condition", me)
}

// condition. Matches requests by comparing one detected attribute with one operation.
type Condition struct {
	Attribute            string          `json:"attribute"`             // The detected attribute that should be compared with the specified operation.
	CompareOperationType string          `json:"compareOperationType"`  // The type of comparison operation that should be applied to the detected attribute.. When using this field over the Settings API, it is stored as a string and must use one of the fixed compare-operation identifiers. The available subset depends on the selected `attribute`.\n\n  - `Exists`, `NotExists`\n - `BoolIsTrue`, `BoolIsFalse`\n - `TagEquals`, `TagKeyEquals`\n - `StringEquals`, `NotStringEquals`, `StringStartsWith`, `NotStringStartsWith`, `StringEndsWith`, `NotStringEndsWith`, `StringContains`, `NotStringContains`\n - `FrameworkEquals`, `NotFrameworkEquals`\n - `IpInRange`, `NotIpInRange`\n - `IntEquals`, `NotIntEquals`, `IntGreaterThan`, `IntLessThan`
	Framework            []FrameworkType `json:"framework,omitempty"`   // The technology that should be compared with the detected attribute.\n\n  Select one or more technologies. The condition matches if the detected attribute value equals (for `FrameworkEquals`) or does not equal (for `NotFrameworkEquals`) at least one of the selected technologies. Possible values: `AXIS`, `CXF`, `HESSIAN`, `JAX_WS_RI`, `JBOSS`, `JERSEY`, `PROGRESS`, `RESTEASY`, `RESTLET`, `SPRING`, `TIBCO`, `WEBLOGIC`, `WEBMETHODS`, `WEBSPHERE`, `WINK`
	IgnoreCase           *bool           `json:"ignoreCase,omitempty"`  // Ignore case sensitivity for texts.
	IntValue             *int            `json:"intValue,omitempty"`    // The integer value to compare the detected attribute with.
	IntValues            []int           `json:"intValues,omitempty"`   // If multiple values are specified, at least one of them must match for the condition to match.
	IpRangeFrom          *string         `json:"ipRangeFrom,omitempty"` // The beginning of the IP range. The condition matches if the detected attribute value is greater than or equal to this value (for `IpInRange`) or less than this value (for `NotIpInRange`).
	IpRangeTo            *string         `json:"ipRangeTo,omitempty"`   // The end of the IP range. The condition matches if the detected attribute value is less than or equal to this value (for `IpInRange`) or greater than this value (for `NotIpInRange`).
	TagValues            []string        `json:"tagValues,omitempty"`   // If multiple values are specified, at least one of them must match for the condition to match.
	TextValues           []string        `json:"textValues,omitempty"`  // If multiple values are specified, at least one of them must match for the condition to match
}

func (me *Condition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "The detected attribute that should be compared with the specified operation.",
			Required:    true,
		},
		"compare_operation_type": {
			Type:        schema.TypeString,
			Description: "The type of comparison operation that should be applied to the detected attribute.. When using this field over the Settings API, it is stored as a string and must use one of the fixed compare-operation identifiers. The available subset depends on the selected `attribute`.\n\n  - `Exists`, `NotExists`\n - `BoolIsTrue`, `BoolIsFalse`\n - `TagEquals`, `TagKeyEquals`\n - `StringEquals`, `NotStringEquals`, `StringStartsWith`, `NotStringStartsWith`, `StringEndsWith`, `NotStringEndsWith`, `StringContains`, `NotStringContains`\n - `FrameworkEquals`, `NotFrameworkEquals`\n - `IpInRange`, `NotIpInRange`\n - `IntEquals`, `NotIntEquals`, `IntGreaterThan`, `IntLessThan`",
			Required:    true,
		},
		"framework": {
			Type:        schema.TypeSet,
			Description: "The technology that should be compared with the detected attribute.\n\n  Select one or more technologies. The condition matches if the detected attribute value equals (for `FrameworkEquals`) or does not equal (for `NotFrameworkEquals`) at least one of the selected technologies. Possible values: `AXIS`, `CXF`, `HESSIAN`, `JAX_WS_RI`, `JBOSS`, `JERSEY`, `PROGRESS`, `RESTEASY`, `RESTLET`, `SPRING`, `TIBCO`, `WEBLOGIC`, `WEBMETHODS`, `WEBSPHERE`, `WINK`",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ignore_case": {
			Type:        schema.TypeBool,
			Description: "Ignore case sensitivity for texts.",
			Optional:    true, // precondition
		},
		"int_value": {
			Type:        schema.TypeInt,
			Description: "The integer value to compare the detected attribute with.",
			Optional:    true, // precondition
		},
		"int_values": {
			Type:        schema.TypeSet,
			Description: "If multiple values are specified, at least one of them must match for the condition to match.",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"ip_range_from": {
			Type:        schema.TypeString,
			Description: "The beginning of the IP range. The condition matches if the detected attribute value is greater than or equal to this value (for `IpInRange`) or less than this value (for `NotIpInRange`).",
			Optional:    true, // precondition
		},
		"ip_range_to": {
			Type:        schema.TypeString,
			Description: "The end of the IP range. The condition matches if the detected attribute value is less than or equal to this value (for `IpInRange`) or greater than this value (for `NotIpInRange`).",
			Optional:    true, // precondition
		},
		"tag_values": {
			Type:        schema.TypeSet,
			Description: "If multiple values are specified, at least one of them must match for the condition to match.",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"text_values": {
			Type:        schema.TypeSet,
			Description: "If multiple values are specified, at least one of them must match for the condition to match",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Condition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attribute":              me.Attribute,
		"compare_operation_type": me.CompareOperationType,
		"framework":              me.Framework,
		"ignore_case":            me.IgnoreCase,
		"int_value":              me.IntValue,
		"int_values":             me.IntValues,
		"ip_range_from":          me.IpRangeFrom,
		"ip_range_to":            me.IpRangeTo,
		"tag_values":             me.TagValues,
		"text_values":            me.TextValues,
	})
}

func (me *Condition) HandlePreconditions() error {
	if (me.IgnoreCase == nil) && (slices.Contains([]string{"TagEquals", "TagKeyEquals", "StringEndsWith", "NotStringEndsWith", "StringStartsWith", "NotStringStartsWith", "StringContains", "NotStringContains", "StringEquals", "NotStringEquals"}, string(me.CompareOperationType))) {
		me.IgnoreCase = new(false)
	}
	if (me.IntValue == nil) && (slices.Contains([]string{"IntGreaterThan", "IntLessThan"}, string(me.CompareOperationType))) {
		me.IntValue = new(0)
	}
	if (me.IgnoreCase != nil) && (!slices.Contains([]string{"TagEquals", "TagKeyEquals", "StringEndsWith", "NotStringEndsWith", "StringStartsWith", "NotStringStartsWith", "StringContains", "NotStringContains", "StringEquals", "NotStringEquals"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'ignore_case' must not be specified unless 'compare_operation_type' is one of ['TagEquals', 'TagKeyEquals', 'StringEndsWith', 'NotStringEndsWith', 'StringStartsWith', 'NotStringStartsWith', 'StringContains', 'NotStringContains', 'StringEquals', 'NotStringEquals']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.IntValue != nil) && (!slices.Contains([]string{"IntGreaterThan", "IntLessThan"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'int_value' must not be specified unless 'compare_operation_type' is one of ['IntGreaterThan', 'IntLessThan']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.IpRangeFrom != nil) && (!slices.Contains([]string{"IpInRange", "NotIpInRange"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'ip_range_from' must not be specified unless 'compare_operation_type' is one of ['IpInRange', 'NotIpInRange']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.IpRangeFrom == nil) && (slices.Contains([]string{"IpInRange", "NotIpInRange"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'ip_range_from' must be specified when 'compare_operation_type' is one of ['IpInRange', 'NotIpInRange']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.IpRangeTo != nil) && (!slices.Contains([]string{"IpInRange", "NotIpInRange"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'ip_range_to' must not be specified unless 'compare_operation_type' is one of ['IpInRange', 'NotIpInRange']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	if (me.IpRangeTo == nil) && (slices.Contains([]string{"IpInRange", "NotIpInRange"}, string(me.CompareOperationType))) {
		return fmt.Errorf("'ip_range_to' must be specified when 'compare_operation_type' is one of ['IpInRange', 'NotIpInRange']; got 'compare_operation_type'='%v'", me.CompareOperationType)
	}
	// ---- Framework []FrameworkType -> {"expectedValues":["FrameworkEquals","NotFrameworkEquals"],"property":"compareOperationType","type":"IN"}
	// ---- IntValues []int -> {"expectedValues":["IntEquals","NotIntEquals"],"property":"compareOperationType","type":"IN"}
	// ---- TagValues []string -> {"expectedValue":"TagEquals","property":"compareOperationType","type":"EQUALS"}
	// ---- TextValues []string -> {"expectedValues":["TagKeyEquals","StringEndsWith","NotStringEndsWith","StringStartsWith","NotStringStartsWith","StringContains","NotStringContains","StringEquals","NotStringEquals"],"property":"compareOperationType","type":"IN"}
	return nil
}

func (me *Condition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attribute":              &me.Attribute,
		"compare_operation_type": &me.CompareOperationType,
		"framework":              &me.Framework,
		"ignore_case":            &me.IgnoreCase,
		"int_value":              &me.IntValue,
		"int_values":             &me.IntValues,
		"ip_range_from":          &me.IpRangeFrom,
		"ip_range_to":            &me.IpRangeTo,
		"tag_values":             &me.TagValues,
		"text_values":            &me.TextValues,
	})
}
