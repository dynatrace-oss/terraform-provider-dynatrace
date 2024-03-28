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

package simpledetectionrule

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled            bool              `json:"enabled"`               // This setting is enabled (`true`) or disabled (`false`)
	GroupIdentifier    string            `json:"groupIdentifier"`       // If Dynatrace detects this property at startup of a process, it will use its value to identify process groups more granular.
	InstanceIdentifier string            `json:"instanceIdentifier"`    // Use a variable to identify instances within a process group.\n\nThe type of variable is the same as selected in 'Property source'.
	ProcessType        *string           `json:"processType,omitempty"` // Note: Not all types can be detected at startup.
	RuleType           DetectionRuleType `json:"ruleType"`              // Possible Values: `Prop`, `Env`
	InsertAfter        string            `json:"-"`
}

func (me *Settings) Name() string {
	name := string(me.RuleType)

	if me.GroupIdentifier != "" {
		name += "_"
		name += me.GroupIdentifier
	}
	if me.InstanceIdentifier != "" {
		name += "_"
		name += me.InstanceIdentifier
	}
	if me.ProcessType != nil && (*me.ProcessType) != "" {
		name += "_"
		name += (*me.ProcessType)
	}

	return name
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"group_identifier": {
			Type:        schema.TypeString,
			Description: "If Dynatrace detects this property at startup of a process, it will use its value to identify process groups more granular.",
			Required:    true,
		},
		"instance_identifier": {
			Type:        schema.TypeString,
			Description: "Use a variable to identify instances within a process group.\n\nThe type of variable is the same as selected in 'Property source'.",
			Required:    true,
		},
		"process_type": {
			Type:        schema.TypeString,
			Description: "Note: Not all types can be detected at startup.",
			Optional:    true,
		},
		"rule_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Prop`, `Env`",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":             me.Enabled,
		"group_identifier":    me.GroupIdentifier,
		"instance_identifier": me.InstanceIdentifier,
		"process_type":        me.ProcessType,
		"rule_type":           me.RuleType,
		"insert_after":        me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":             &me.Enabled,
		"group_identifier":    &me.GroupIdentifier,
		"instance_identifier": &me.InstanceIdentifier,
		"process_type":        &me.ProcessType,
		"rule_type":           &me.RuleType,
		"insert_after":        &me.InsertAfter,
	})
}
