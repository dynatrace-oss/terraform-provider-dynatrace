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
	Enabled  bool                `json:"enabled"`         // This setting is enabled (`true`) or disabled (`false`)
	Metadata MetadataItems       `json:"metadata"`        // Set of additional key-value properties to be attached to the triggered event.
	Name     string              `json:"name"`            // Monitored rule name
	Rules    DetectionConditions `json:"rules"`           // Define process detection rules by selecting a process property and a condition. Each monitoring rule can have multiple detection rules associated with it.
	Scope    *string             `json:"-" scope:"scope"` // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
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
			Description: "Set of additional key-value properties to be attached to the triggered event.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(MetadataItems).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Monitored rule name",
			Required:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "Define process detection rules by selecting a process property and a condition. Each monitoring rule can have multiple detection rules associated with it.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(DetectionConditions).Schema()},
			MinItems: 1,
			MaxItems: 1,
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
		"enabled":  me.Enabled,
		"metadata": me.Metadata,
		"name":     me.Name,
		"rules":    me.Rules,
		"scope":    me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":  &me.Enabled,
		"metadata": &me.Metadata,
		"name":     &me.Name,
		"rules":    &me.Rules,
		"scope":    &me.Scope,
	})
}
