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

package processgroupingrules

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DetectionConditions []*DetectionCondition

func (me *DetectionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detection_condition": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DetectionCondition).Schema()},
		},
	}
}

func (me DetectionConditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("detection_condition", me)
}

func (me *DetectionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("detection_condition", me)
}

type DetectionCondition struct {
	CaseSensitive bool    `json:"caseSensitive"`       // When enabled, matching conditions are case sensitive. When disabled, matching conditions are case insensitive
	Condition     *string `json:"condition,omitempty"` // - $contains(svc) – Matches if svc appears anywhere in the process property value.\n - $eq(svc.exe) – Matches if svc.exe matches the process property value exactly.\n - $prefix(svc) – Matches if app matches the prefix of the process property value.\n - $suffix(svc.py) – Matches if svc.py matches the suffix of the process property value.\n\n  For example, $suffix(svc.py) would detect processes named loyaltysvc.py and paymentssvc.py.\n\n  For more details, see [documentation](https://dt-url.net/j142w57).
	Name          *string `json:"name,omitempty"`      // If Dynatrace detects this property at startup of a process, it will be matched to this grouping rule.
	Property      string  `json:"property"`            // 2.1. Property
}

func (me *DetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "When enabled, matching conditions are case sensitive. When disabled, matching conditions are case insensitive",
			Required:    true,
		},
		"condition": {
			Type:        schema.TypeString,
			Description: "- $contains(svc) – Matches if svc appears anywhere in the process property value.\n - $eq(svc.exe) – Matches if svc.exe matches the process property value exactly.\n - $prefix(svc) – Matches if app matches the prefix of the process property value.\n - $suffix(svc.py) – Matches if svc.py matches the suffix of the process property value.\n\n  For example, $suffix(svc.py) would detect processes named loyaltysvc.py and paymentssvc.py.\n\n  For more details, see [documentation](https://dt-url.net/j142w57).",
			Optional:    true, // precondition
		},
		"name": {
			Type:        schema.TypeString,
			Description: "If Dynatrace detects this property at startup of a process, it will be matched to this grouping rule.",
			Optional:    true, // precondition
		},
		"property": {
			Type:        schema.TypeString,
			Description: "2.1. Property",
			Required:    true,
		},
	}
}

func (me *DetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"case_sensitive": me.CaseSensitive,
		"condition":      me.Condition,
		"name":           me.Name,
		"property":       me.Property,
	})
}

func (me *DetectionCondition) HandlePreconditions() error {
	if (me.Condition != nil) && (slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(me.Property))) {
		return fmt.Errorf("'condition' must not be specified unless 'property' is not one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", me.Property)
	}
	if (me.Condition == nil) && (!slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(me.Property))) {
		return fmt.Errorf("'condition' must be specified when 'property' is not one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", me.Property)
	}
	if (me.Name != nil) && (!slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(me.Property))) {
		return fmt.Errorf("'name' must not be specified unless 'property' is one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", me.Property)
	}
	if (me.Name == nil) && (slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(me.Property))) {
		return fmt.Errorf("'name' must be specified when 'property' is one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", me.Property)
	}
	return nil
}

func (me *DetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"case_sensitive": &me.CaseSensitive,
		"condition":      &me.Condition,
		"name":           &me.Name,
		"property":       &me.Property,
	})
}
