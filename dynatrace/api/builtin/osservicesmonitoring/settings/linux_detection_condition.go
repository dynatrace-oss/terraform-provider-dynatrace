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

package osservicesmonitoring

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type LinuxDetectionConditions []*LinuxDetectionCondition

func (me *LinuxDetectionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"linux_detection_condition": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(LinuxDetectionCondition).Schema()},
		},
	}
}

func (me LinuxDetectionConditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("linux_detection_condition", me)
}

func (me *LinuxDetectionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("linux_detection_condition", me)
}

type LinuxDetectionCondition struct {
	Condition             *string                `json:"condition,omitempty"`             // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n- `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n- `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n- `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\nAvailable logic operations:\n- `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n- `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n- `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.\n\nBrackets **(** and **)** that are part of the matched property **must be escaped with a tilde (~)**
	HostMetadataCondition *HostMetadataCondition `json:"hostMetadataCondition,omitempty"` // Custom metadata
	Property              *LinuxServiceProp      `json:"property,omitempty"`              // Possible Values: `ServiceName`, `StartupType`
	RuleType              *RuleType              `json:"ruleType,omitempty"`              // Possible Values: `RuleTypeHost`, `RuleTypeOsService`
	StartupCondition      *string                `json:"startupCondition,omitempty"`      // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(enabled)` – Matches services with startup type equal to enabled.\n\nAvailable logic operations:\n- `$not($eq(enabled))` – Matches services with startup type different from enabled.\n- `$or($eq(enabled),$eq(disabled))` - Matches services that are either enabled or disabled.\n\nUse one of the following values as a parameter for this condition:\n\n- `enabled`\n- `enabled-runtime`\n- `static`\n- `disabled`
}

func (me *LinuxDetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n- `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n- `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n- `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\nAvailable logic operations:\n- `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n- `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n- `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.\n\nBrackets **(** and **)** that are part of the matched property **must be escaped with a tilde (~)**",
			Optional:    true, // precondition
		},
		"host_metadata_condition": {
			Type:        schema.TypeList,
			Description: "Custom metadata",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(HostMetadataCondition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ServiceName`, `StartupType`",
			Optional:    true, // precondition
		},
		"rule_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `RuleTypeHost`, `RuleTypeOsService`",
			Optional:    true, // nullable
		},
		"startup_condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(enabled)` – Matches services with startup type equal to enabled.\n\nAvailable logic operations:\n- `$not($eq(enabled))` – Matches services with startup type different from enabled.\n- `$or($eq(enabled),$eq(disabled))` - Matches services that are either enabled or disabled.\n\nUse one of the following values as a parameter for this condition:\n\n- `enabled`\n- `enabled-runtime`\n- `static`\n- `disabled`",
			Optional:    true, // precondition
		},
	}
}

func (me *LinuxDetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":               me.Condition,
		"host_metadata_condition": me.HostMetadataCondition,
		"property":                me.Property,
		"rule_type":               me.RuleType,
		"startup_condition":       me.StartupCondition,
	})
}

func (me *LinuxDetectionCondition) HandlePreconditions() error {
	if (me.Condition == nil) && (me.Property != nil && slices.Contains([]string{"ServiceName"}, string(*me.Property))) {
		return fmt.Errorf("'condition' must be specified if 'property' is set to '%v'", me.Property)
	}
	if (me.HostMetadataCondition == nil) && (me.RuleType != nil && slices.Contains([]string{"RuleTypeHost"}, string(*me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.HostMetadataCondition != nil) && (me.RuleType != nil && !slices.Contains([]string{"RuleTypeHost"}, string(*me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.StartupCondition == nil) && (me.Property != nil && slices.Contains([]string{"StartupType"}, string(*me.Property))) {
		return fmt.Errorf("'startup_condition' must be specified if 'property' is set to '%v'", me.Property)
	}
	// ---- Property *LinuxServiceProp -> {"preconditions":[{"property":"ruleType","type":"NULL"},{"expectedValues":["RuleTypeOsService"],"property":"ruleType","type":"IN"}],"type":"OR"}
	return nil
}

func (me *LinuxDetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":               &me.Condition,
		"host_metadata_condition": &me.HostMetadataCondition,
		"property":                &me.Property,
		"rule_type":               &me.RuleType,
		"startup_condition":       &me.StartupCondition,
	})
}
