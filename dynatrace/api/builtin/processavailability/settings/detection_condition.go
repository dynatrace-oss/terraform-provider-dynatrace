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

package processavailability

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type DetectionConditions []*DetectionCondition

func (me *DetectionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DetectionCondition).Schema()},
		},
	}
}

func (me DetectionConditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *DetectionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

type DetectionCondition struct {
	Condition             *string                `json:"condition,omitempty"`             // - $contains(svc) – Matches if svc appears anywhere in the process property value.\n- $eq(svc.exe) – Matches if svc.exe matches the process property value exactly.\n- $prefix(svc) – Matches if app matches the prefix of the process property value.\n- $suffix(svc.py) – Matches if svc.py matches the suffix of the process property value.\n\nFor example, $suffix(svc.py) would detect processes named loyaltysvc.py and paymentssvc.py.\n\nFor more details, see [Process availability](https://dt-url.net/v923x37).
	HostMetadataCondition *HostMetadataCondition `json:"hostMetadataCondition,omitempty"` // Host custom metadata refers to user-defined key-value pairs that you can assign to hosts monitored by Dynatrace.\n\nBy defining custom metadata, you can enrich the monitoring data with context specific to your organization's needs, such as environment names, team ownership, application versions, or any other relevant details.\n\nSee [Define tags and metadata for hosts](https://dt-url.net/w3hv0kbw).
	Property              *ProcessItem           `json:"property,omitempty"`              // Possible Values: `CommandLine`, `Executable`, `ExecutablePath`, `User`
	RuleType              RuleType               `json:"ruleType"`                        // Possible Values: `RuleTypeHost`, `RuleTypeProcess`
}

func (me *DetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "- $contains(svc) – Matches if svc appears anywhere in the process property value.\n- $eq(svc.exe) – Matches if svc.exe matches the process property value exactly.\n- $prefix(svc) – Matches if app matches the prefix of the process property value.\n- $suffix(svc.py) – Matches if svc.py matches the suffix of the process property value.\n\nFor example, $suffix(svc.py) would detect processes named loyaltysvc.py and paymentssvc.py.\n\nFor more details, see [Process availability](https://dt-url.net/v923x37).",
			Optional:    true, // precondition
		},
		"host_metadata_condition": {
			Type:        schema.TypeList,
			Description: "Host custom metadata refers to user-defined key-value pairs that you can assign to hosts monitored by Dynatrace.\n\nBy defining custom metadata, you can enrich the monitoring data with context specific to your organization's needs, such as environment names, team ownership, application versions, or any other relevant details.\n\nSee [Define tags and metadata for hosts](https://dt-url.net/w3hv0kbw).",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(HostMetadataCondition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CommandLine`, `Executable`, `ExecutablePath`, `User`",
			Optional:    true, // precondition
		},
		"rule_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `RuleTypeHost`, `RuleTypeProcess`",
			Optional:    true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// rule_type was introduced in v286 as a required field, added code below to have successful results for old/new tenants.
				return newValue == ""
			},
			// Default: "RuleTypeProcess",
		},
	}
}

func (me *DetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":               me.Condition,
		"host_metadata_condition": me.HostMetadataCondition,
		"property":                me.Property,
		"rule_type":               me.RuleType,
	})
}

func (me *DetectionCondition) HandlePreconditions() error {
	if (me.HostMetadataCondition == nil) && (slices.Contains([]string{"RuleTypeHost"}, string(me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.HostMetadataCondition != nil) && (!slices.Contains([]string{"RuleTypeHost"}, string(me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	// ---- Condition *string -> {"preconditions":[{"expectedValues":["RuleTypeProcess"],"property":"ruleType","type":"IN"}],"type":"OR"}
	// ---- Property *ProcessItem -> {"preconditions":[{"expectedValues":["RuleTypeProcess"],"property":"ruleType","type":"IN"}],"type":"OR"}
	return nil
}

func (me *DetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"condition":               &me.Condition,
		"host_metadata_condition": &me.HostMetadataCondition,
		"property":                &me.Property,
		"rule_type":               &me.RuleType,
	})
	// RuleType was introduced in v286 as a required field, added code below to have successful results for old/new tenants.
	if me.RuleType == "" {
		me.RuleType = RuleTypes.Ruletypeprocess
	}
	return err
}
