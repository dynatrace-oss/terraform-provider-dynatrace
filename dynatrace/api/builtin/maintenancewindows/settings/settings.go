/**
* @license
* Copyright 2026 Dynatrace LLC
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

package maintenancewindows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AutoDelete  bool                `json:"autoDelete"`            // When enabled, this maintenance window configuration will be automatically deleted 30 days after the end of its last execution.
	Description *string             `json:"description,omitempty"` // Description
	Enabled     bool                `json:"enabled"`               // This setting is enabled (`true`) or disabled (`false`)
	Filter      string              `json:"filter"`                // DQL Filter
	Name        string              `json:"name"`                  // Name of the maintenance window
	Schedule    *ScheduleDefinition `json:"schedule"`              // Schedule definition
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_delete": {
			Type:        schema.TypeBool,
			Description: "When enabled, this maintenance window configuration will be automatically deleted 30 days after the end of its last execution.",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"filter": {
			Type:        schema.TypeString,
			Description: "DQL Filter",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the maintenance window",
			Required:    true,
		},
		"schedule": {
			Type:        schema.TypeList,
			Description: "Schedule definition",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ScheduleDefinition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_delete": me.AutoDelete,
		"description": me.Description,
		"enabled":     me.Enabled,
		"filter":      me.Filter,
		"name":        me.Name,
		"schedule":    me.Schedule,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_delete": &me.AutoDelete,
		"description": &me.Description,
		"enabled":     &me.Enabled,
		"filter":      &me.Filter,
		"name":        &me.Name,
		"schedule":    &me.Schedule,
	})
}
