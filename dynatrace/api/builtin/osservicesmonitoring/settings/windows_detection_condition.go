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
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WindowsDetectionConditions []*WindowsDetectionCondition

func (me *WindowsDetectionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detection_conditions_window": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(WindowsDetectionCondition).Schema()},
		},
	}
}

func (me WindowsDetectionConditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("detection_conditions_window", me)
}

func (me *WindowsDetectionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("detection_conditions_window", me)
}

type WindowsDetectionCondition struct {
	Condition             *string                `json:"condition,omitempty"`             // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n  - `$match(ip?tables*)` – Matches string with wildcards: `*` any number (including zero) of characters and `?` exactly one character.\n - `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n - `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n - `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n - `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\n  Available logic operations:\n - `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n - `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n - `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.\n\n  Brackets **(** and **)** that are part of the matched property **must be escaped with a tilde (~)**
	HostMetadataCondition *HostMetadataCondition `json:"hostMetadataCondition,omitempty"` // Resource attribute
	Property              *WindowsServiceProps   `json:"property,omitempty"`              // Service property. Possible values: `DisplayName`, `Manufacturer`, `Path`, `ServiceName`, `StartupType`
	RuleType              *RuleType              `json:"ruleType,omitempty"`              // Rule scope. Possible values: `RuleTypeHost`, `RuleTypeOsService`
	StartupCondition      *string                `json:"startupCondition,omitempty"`      // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n  - `$eq(manual)` – Matches services that are started manually.\n\n  Available logic operations:\n - `$not($eq(auto))` – Matches services with startup type different from Automatic.\n - `$or($eq(auto),$eq(manual))` – Matches if service's startup type is either Automatic or Manual.\n\n  Use one of the following values as a parameter for this condition:\n\n  - `manual` for Manual\n - `manual_trigger` for Manual (Trigger Start)\n - `auto` for Automatic\n - `auto_delay` for Automatic (Delayed Start)\n - `auto_trigger` for Automatic (Trigger Start)\n - `auto_delay_trigger` for Automatic (Delayed Start, Trigger Start)\n - `disabled` for Disabled
}

func (me *WindowsDetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n  - `$match(ip?tables*)` – Matches string with wildcards: `*` any number (including zero) of characters and `?` exactly one character.\n - `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n - `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n - `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n - `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\n  Available logic operations:\n - `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n - `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n - `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.\n\n  Brackets **(** and **)** that are part of the matched property **must be escaped with a tilde (~)**",
			Optional:    true, // precondition
		},
		"host_metadata_condition": {
			Type:        schema.TypeList,
			Description: "Resource attribute",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(HostMetadataCondition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Service property. Possible values: `DisplayName`, `Manufacturer`, `Path`, `ServiceName`, `StartupType`",
			Optional:    true, // precondition
		},
		"rule_type": {
			Type:        schema.TypeString,
			Description: "Rule scope. Possible values: `RuleTypeHost`, `RuleTypeOsService`",
			Optional:    true, // nullable
		},
		"startup_condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n  - `$eq(manual)` – Matches services that are started manually.\n\n  Available logic operations:\n - `$not($eq(auto))` – Matches services with startup type different from Automatic.\n - `$or($eq(auto),$eq(manual))` – Matches if service's startup type is either Automatic or Manual.\n\n  Use one of the following values as a parameter for this condition:\n\n  - `manual` for Manual\n - `manual_trigger` for Manual (Trigger Start)\n - `auto` for Automatic\n - `auto_delay` for Automatic (Delayed Start)\n - `auto_trigger` for Automatic (Trigger Start)\n - `auto_delay_trigger` for Automatic (Delayed Start, Trigger Start)\n - `disabled` for Disabled",
			Optional:    true, // precondition
		},
	}
}

func (me *WindowsDetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":               me.Condition,
		"host_metadata_condition": me.HostMetadataCondition,
		"property":                me.Property,
		"rule_type":               me.RuleType,
		"startup_condition":       me.StartupCondition,
	})
}

func (me *WindowsDetectionCondition) HandlePreconditions() error {
	if (me.Condition != nil) && (me.Property == nil || !slices.Contains([]string{"Manufacturer", "ServiceName", "DisplayName", "Path"}, string(*me.Property))) {
		return fmt.Errorf("'condition' must not be specified unless 'property' is one of ['Manufacturer', 'ServiceName', 'DisplayName', 'Path']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	if (me.Condition == nil) && (me.Property != nil && slices.Contains([]string{"Manufacturer", "ServiceName", "DisplayName", "Path"}, string(*me.Property))) {
		return fmt.Errorf("'condition' must be specified when 'property' is one of ['Manufacturer', 'ServiceName', 'DisplayName', 'Path']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	if (me.HostMetadataCondition != nil) && (me.RuleType == nil || !slices.Contains([]string{"RuleTypeHost"}, string(*me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must not be specified unless 'rule_type' is one of ['RuleTypeHost']; got 'rule_type'='%v'", opt.ValOrNil(me.RuleType))
	}
	if (me.HostMetadataCondition == nil) && (me.RuleType != nil && slices.Contains([]string{"RuleTypeHost"}, string(*me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must be specified when 'rule_type' is one of ['RuleTypeHost']; got 'rule_type'='%v'", opt.ValOrNil(me.RuleType))
	}
	if (me.Property != nil) && ((me.RuleType != nil) && (me.RuleType == nil || !slices.Contains([]string{"RuleTypeOsService"}, string(*me.RuleType)))) {
		return fmt.Errorf("'property' must not be specified unless ('rule_type' is not set or 'rule_type' is one of ['RuleTypeOsService']); got 'rule_type'='%v'", opt.ValOrNil(me.RuleType))
	}
	if (me.Property == nil) && ((me.RuleType == nil) || (me.RuleType != nil && slices.Contains([]string{"RuleTypeOsService"}, string(*me.RuleType)))) {
		return fmt.Errorf("'property' must be specified when ('rule_type' is not set or 'rule_type' is one of ['RuleTypeOsService']); got 'rule_type'='%v'", opt.ValOrNil(me.RuleType))
	}
	if (me.StartupCondition != nil) && (me.Property == nil || !slices.Contains([]string{"StartupType"}, string(*me.Property))) {
		return fmt.Errorf("'startup_condition' must not be specified unless 'property' is one of ['StartupType']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	if (me.StartupCondition == nil) && (me.Property != nil && slices.Contains([]string{"StartupType"}, string(*me.Property))) {
		return fmt.Errorf("'startup_condition' must be specified when 'property' is one of ['StartupType']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	return nil
}

func (me *WindowsDetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":               &me.Condition,
		"host_metadata_condition": &me.HostMetadataCondition,
		"property":                &me.Property,
		"rule_type":               &me.RuleType,
		"startup_condition":       &me.StartupCondition,
	})
}
