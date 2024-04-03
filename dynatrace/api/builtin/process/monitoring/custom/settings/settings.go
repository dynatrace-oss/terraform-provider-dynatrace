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

package customprocessmonitoring

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Condition   *Condition     `json:"condition"`             // Condition
	Enabled     bool           `json:"enabled"`               // This setting is enabled (`true`) or disabled (`false`)
	Mode        MonitoringMode `json:"mode"`                  // Possible Values: `MONITORING_ON`, `MONITORING_OFF`
	HostGroupID string         `json:"-" scope:"hostGroupId"` // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope
	InsertAfter string         `json:"-"`
}

func (me *Settings) Name() string {
	name := me.HostGroupID

	if me.Condition != nil {
		if string(me.Condition.Item) != "" {
			if name != "" {
				name += "_"
			}
			name += string(me.Condition.Item)
		}
		if string(me.Condition.Operator) != "" {
			if name != "" {
				name += "_"
			}
			name += string(me.Condition.Operator)
		}
		if me.Condition.Value != nil && (*me.Condition.Value) != "" {
			if name != "" {
				name += "_"
			}
			name += (*me.Condition.Value)
		}
	}

	if name != "" {
		return name
	}

	return uuid.NewString()
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeList,
			Description: "Condition",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Condition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `MONITORING_ON`, `MONITORING_OFF`",
			Required:    true,
		},
		"host_group_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
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
		"condition":     me.Condition,
		"enabled":       me.Enabled,
		"mode":          me.Mode,
		"host_group_id": me.HostGroupID,
		"insert_after":  me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":     &me.Condition,
		"enabled":       &me.Enabled,
		"mode":          &me.Mode,
		"host_group_id": &me.HostGroupID,
		"insert_after":  &me.InsertAfter,
	})
}
