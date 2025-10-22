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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled          bool                `json:"enabled"`            // This setting is enabled (`true`) or disabled (`false`)
	Metadata         MetadataItems       `json:"metadata,omitempty"` // Set of additional key-value properties to be attached to the triggered event. You can retrieve the available property keys using the [Events API v2](https://dt-url.net/9622g1w). Additionally any Host resource attribute can be dynamically substituted (agent 1.325+).
	Name             string              `json:"name"`               // Monitored rule name
	Rules            DetectionConditions `json:"rules,omitempty"`    // Define process detection rules by selecting a process property and a condition. Each monitoring rule can have multiple detection rules associated with it.
	Scope            *string             `json:"-" scope:"scope"`    // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	MinimumProcesses int                 `json:"minimumProcesses"`   // Specify a minimum number of processes matching the monitoring rule. If it's not satisfied, an alert will open.
	OperatingSystem  []OperatingSystem   `json:"operatingSystem"`    // Select the operating systems on which the monitoring rule should be applied.
	InsertAfter      string              `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Set of additional key-value properties to be attached to the triggered event. You can retrieve the available property keys using the [Events API v2](https://dt-url.net/9622g1w). Additionally any Host resource attribute can be dynamically substituted (agent 1.325+).",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(MetadataItems).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Monitoring rule name",
			Required:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "Define process detection rules by selecting a process property and a condition. Each monitoring rule can have multiple detection rules associated with it.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DetectionConditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"minimum_processes": {
			Type:        schema.TypeInt,
			Description: "Specify a minimum number of processes matching the monitoring rule. If it's not satisfied, an alert will open.",
			Optional:    true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// minimum_processes was introduced in v286 as a required field, added code below to have successful results for old/new tenants.
				return newValue == "0"
			},
			// Default: 1,
		},
		"operating_system": {
			Type:        schema.TypeSet,
			Description: "Select the operating systems on which the monitoring rule should be applied.",
			Optional:    true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// operating_system was introduced in v286 as a required field, added code below to have successful results for old/new tenants.
				if newValue == "0" || newValue == "" {
					return true
				}
				return false
			},
			Elem: &schema.Schema{Type: schema.TypeString},
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
		"enabled":           me.Enabled,
		"metadata":          me.Metadata,
		"name":              me.Name,
		"rules":             me.Rules,
		"scope":             me.Scope,
		"minimum_processes": me.MinimumProcesses,
		"operating_system":  me.OperatingSystem,
		"insert_after":      me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"enabled":           &me.Enabled,
		"metadata":          &me.Metadata,
		"name":              &me.Name,
		"rules":             &me.Rules,
		"scope":             &me.Scope,
		"minimum_processes": &me.MinimumProcesses,
		"operating_system":  &me.OperatingSystem,
		"insert_after":      &me.InsertAfter,
	})
	// MinimumProcesses and OperatingSystem were introduced in v286 as required fields, added code below to have successful results for old/new tenants.
	if me.MinimumProcesses == 0 {
		me.MinimumProcesses = 1
	}
	if len(me.OperatingSystem) == 0 {
		me.OperatingSystem = []OperatingSystem{"AIX", "LINUX", "WINDOWS"}
	}
	return err
}
