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
	Condition        *string          `json:"condition,omitempty"`        // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n- `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n- `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n- `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\nAvailable logic operations:\n- `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n- `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n- `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.
	Property         LinuxServiceProp `json:"property"`                   // Possible Values: `ServiceName`, `StartupType`
	StartupCondition *string          `json:"startupCondition,omitempty"` // This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(enabled)` – Matches services with startup type equal to enabled.\n\nAvailable logic operations:\n- `$not($eq(enabled))` – Matches services with startup type different from enabled.\n- `$or($eq(enabled),$eq(disabled))` - Matches services that are either enabled or disabled.\n\nUse one of the following values as a parameter for this condition:\n\n- `enabled`\n- `enabled-runtime`\n- `static`\n- `disabled`
}

func (me *LinuxDetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$contains(ssh)` – Matches if `ssh` appears anywhere in the service's property value.\n- `$eq(sshd)` – Matches if `sshd` matches the service's property value exactly.\n- `$prefix(ss)` – Matches if `ss` matches the prefix of the service's property value.\n- `$suffix(hd)` – Matches if `hd` matches the suffix of the service's property value.\n\nAvailable logic operations:\n- `$not($eq(sshd))` – Matches if the service's property value is different from `sshd`.\n- `$and($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` and ends with `hd`.\n- `$or($prefix(ss),$suffix(hd))` – Matches if service's property value starts with `ss` or ends with `hd`.",
			Optional:    true,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ServiceName`, `StartupType`",
			Required:    true,
		},
		"startup_condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format. See [OS services monitoring](https://dt-url.net/vl03xzk).\n\n- `$eq(enabled)` – Matches services with startup type equal to enabled.\n\nAvailable logic operations:\n- `$not($eq(enabled))` – Matches services with startup type different from enabled.\n- `$or($eq(enabled),$eq(disabled))` - Matches services that are either enabled or disabled.\n\nUse one of the following values as a parameter for this condition:\n\n- `enabled`\n- `enabled-runtime`\n- `static`\n- `disabled`",
			Optional:    true,
		},
	}
}

func (me *LinuxDetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":         me.Condition,
		"property":          me.Property,
		"startup_condition": me.StartupCondition,
	})
}

func (me *LinuxDetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":         &me.Condition,
		"property":          &me.Property,
		"startup_condition": &me.StartupCondition,
	})
}
