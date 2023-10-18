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

package customlogsourcesettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name              string           `json:"config-item-title"` // Name
	Context           Contexts         `json:"context,omitempty"` // Define Custom Log Source only within context if provided
	Custom_log_source *CustomLogSource `json:"custom-log-source"`
	Enabled           bool             `json:"enabled"`         // This setting is enabled (`true`) or disabled (`false`)
	Scope             *string          `json:"-" scope:"scope"` // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"context": {
			Type:        schema.TypeList,
			Description: "Define Custom Log Source only within context if provided",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Contexts).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"custom_log_source": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(CustomLogSource).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":              me.Name,
		"context":           me.Context,
		"custom_log_source": me.Custom_log_source,
		"enabled":           me.Enabled,
		"scope":             me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":              &me.Name,
		"context":           &me.Context,
		"custom_log_source": &me.Custom_log_source,
		"enabled":           &me.Enabled,
		"scope":             &me.Scope,
	})
}
