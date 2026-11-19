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

package logstoragesettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name            string   `json:"config-item-title"` // Name
	Enabled         bool     `json:"enabled"`           // This setting is enabled (`true`) or disabled (`false`)
	InsertAfter     *string  `json:"-"`                 // Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched
	Matchers        Matchers `json:"matchers,omitempty"`
	Scope           *string  `json:"-" scope:"scope"` // The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.
	Send_to_storage bool     `json:"send-to-storage"` // If `true` matching logs will be included in storage. If `false` matching logs will be excluded from storage.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Computed:    true,
			Optional:    true,
		},
		"matchers": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Matchers).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			ForceNew:    true,
			Optional:    true,
			Default:     "environment",
		},
		"send_to_storage": {
			Type:        schema.TypeBool,
			Description: "If `true` matching logs will be included in storage. If `false` matching logs will be excluded from storage.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":            me.Name,
		"enabled":         me.Enabled,
		"insert_after":    me.InsertAfter,
		"matchers":        me.Matchers,
		"scope":           me.Scope,
		"send_to_storage": me.Send_to_storage,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":            &me.Name,
		"enabled":         &me.Enabled,
		"insert_after":    &me.InsertAfter,
		"matchers":        &me.Matchers,
		"scope":           &me.Scope,
		"send_to_storage": &me.Send_to_storage,
	})
}
