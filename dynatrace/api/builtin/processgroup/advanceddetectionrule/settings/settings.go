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

package advanceddetectionrule

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled            bool                      `json:"enabled"`            // This setting is enabled (`true`) or disabled (`false`)
	GroupExtraction    *ProcessGroupExtraction   `json:"groupExtraction"`    // You can define the properties that should be used to identify your process groups.
	InstanceExtraction ProcessInstanceExtraction `json:"instanceExtraction"` // You can define the properties that should be used to identify your process instances.
	ProcessDetection   *ProcessGroupDetection    `json:"processDetection"`   // Apply this rule to processes where the selected property contains the specified string.
	InsertAfter        string                    `json:"-"`
}

func (me *Settings) Name() string {
	name := ""
	if me.ProcessDetection != nil {
		if me.ProcessDetection.Property != "" {
			name = me.ProcessDetection.Property
		}
		if me.ProcessDetection.ContainedString != "" {
			if name != "" {
				name += "_"
			}
			name += me.ProcessDetection.ContainedString
		}
	}

	if name == "" {
		return uuid.NewString()
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
		"group_extraction": {
			Type:        schema.TypeList,
			Description: "You can define the properties that should be used to identify your process groups.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ProcessGroupExtraction).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"instance_extraction": {
			Type:        schema.TypeList,
			Description: "You can define the properties that should be used to identify your process instances.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(ProcessInstanceExtraction).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"process_detection": {
			Type:        schema.TypeList,
			Description: "Apply this rule to processes where the selected property contains the specified string.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ProcessGroupDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
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
		"group_extraction":    me.GroupExtraction,
		"instance_extraction": me.InstanceExtraction.AddrIfNotEmpty(),
		"process_detection":   me.ProcessDetection,
		"insert_after":        me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":             &me.Enabled,
		"group_extraction":    &me.GroupExtraction,
		"instance_extraction": &me.InstanceExtraction,
		"process_detection":   &me.ProcessDetection,
		"insert_after":        &me.InsertAfter,
	})
}
