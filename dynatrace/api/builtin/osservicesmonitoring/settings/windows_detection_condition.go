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
	Condition        *string             `json:"condition,omitempty"`        // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n- `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n- `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n- `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\nAvailable logic operations:\n- `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n- `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n- `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.
	Property         WindowsServiceProps `json:"property"`                   // Possible Values: `DisplayName`, `Manufacturer`, `Path`, `ServiceName`, `StartupType`
	StartupCondition *string             `json:"startupCondition,omitempty"` // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(manual)` – Matches services that are started manually.\n\nAvailable logic operations:\n- `$not($eq(auto))` – Matches services with startup type different from Automatic.\n- `$or($eq(auto),$eq(manual))` – Matches if service's startup type is either Automatic or Manual.\n\nUse one of the following values as a parameter for this condition:\n\n- `manual` for Manual\n- `manual_trigger` for Manual (Trigger Start)\n- `auto` for Automatic\n- `auto_delay` for Automatic (Delayed Start)\n- `auto_trigger` for Automatic (Trigger Start)\n- `auto_delay_trigger` for Automatic (Delayed Start, Trigger Start)\n- `disabled` for Disabled
}

func (me *WindowsDetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n- `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n- `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n- `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\nAvailable logic operations:\n- `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n- `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n- `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.",
			Optional:    true,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DisplayName`, `Manufacturer`, `Path`, `ServiceName`, `StartupType`",
			Required:    true,
		},
		"startup_condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(manual)` – Matches services that are started manually.\n\nAvailable logic operations:\n- `$not($eq(auto))` – Matches services with startup type different from Automatic.\n- `$or($eq(auto),$eq(manual))` – Matches if service's startup type is either Automatic or Manual.\n\nUse one of the following values as a parameter for this condition:\n\n- `manual` for Manual\n- `manual_trigger` for Manual (Trigger Start)\n- `auto` for Automatic\n- `auto_delay` for Automatic (Delayed Start)\n- `auto_trigger` for Automatic (Trigger Start)\n- `auto_delay_trigger` for Automatic (Delayed Start, Trigger Start)\n- `disabled` for Disabled",
			Optional:    true,
		},
	}
}

func (me *WindowsDetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":         me.Condition,
		"property":          me.Property,
		"startup_condition": me.StartupCondition,
	})
}

func (me *WindowsDetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":         &me.Condition,
		"property":          &me.Property,
		"startup_condition": &me.StartupCondition,
	})
}
