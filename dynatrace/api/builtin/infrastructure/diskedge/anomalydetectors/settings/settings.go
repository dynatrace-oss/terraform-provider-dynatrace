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

package anomalydetectors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Alerts                 Alerts                     `json:"alerts,omitempty"`                 // Alerts
	DiskNameFilters        []string                   `json:"diskNameFilters,omitempty"`        // Disk will be included in this policy if **any** of the filters match
	Enabled                bool                       `json:"enabled"`                          // This setting is enabled (`true`) or disabled (`false`)
	EventProperties        MetadataItems              `json:"eventProperties,omitempty"`        // Set of additional key-value properties to be attached to the triggered event.
	HostMetadataConditions HostMetadataConditionTypes `json:"hostMetadataConditions,omitempty"` // The policy will be enabled if **all** conditions are met
	OperatingSystem        []EoperatingSystem         `json:"operatingSystem,omitempty"`        // Select the operating systems on which policy should be applied
	PolicyName             string                     `json:"policyName"`                       // Policy name
	Scope                  *string                    `json:"-" scope:"scope"`                  // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	if me.Scope == nil {
		return "environment" + "_" + me.PolicyName
	}
	return *me.Scope + "_" + me.PolicyName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alerts": {
			Type:        schema.TypeList,
			Description: "Alerts",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Alerts).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"disk_name_filters": {
			Type:        schema.TypeSet,
			Description: "Disk will be included in this policy if **any** of the filters match",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"event_properties": {
			Type:        schema.TypeList,
			Description: "Set of additional key-value properties to be attached to the triggered event.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(MetadataItems).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"host_metadata_conditions": {
			Type:        schema.TypeList,
			Description: "The policy will be enabled if **all** conditions are met",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(HostMetadataConditionTypes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"operating_system": {
			Type:        schema.TypeSet,
			Description: "Select the operating systems on which policy should be applied",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"policy_name": {
			Type:        schema.TypeString,
			Description: "Policy name",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"alerts":                   me.Alerts,
		"disk_name_filters":        me.DiskNameFilters,
		"enabled":                  me.Enabled,
		"event_properties":         me.EventProperties,
		"host_metadata_conditions": me.HostMetadataConditions,
		"operating_system":         me.OperatingSystem,
		"policy_name":              me.PolicyName,
		"scope":                    me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"alerts":                   &me.Alerts,
		"disk_name_filters":        &me.DiskNameFilters,
		"enabled":                  &me.Enabled,
		"event_properties":         &me.EventProperties,
		"host_metadata_conditions": &me.HostMetadataConditions,
		"operating_system":         &me.OperatingSystem,
		"policy_name":              &me.PolicyName,
		"scope":                    &me.Scope,
	})
}
