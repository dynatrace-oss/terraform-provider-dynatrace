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

package managementzones

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Description *string `json:"description,omitempty"` // Description
	Name        string  `json:"name"`                  // **Be careful when renaming** - if there are policies that are referencing this Management zone, they will need to be adapted to the new name!
	Rules       Rules   `json:"rules,omitempty"`       // Rules
	LegacyID    *string `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Description",
			Optional:    true, // nullable
		},
		"name": {
			Type:        schema.TypeString,
			Description: "**Be careful when renaming** - if there are policies that are referencing this Management zone, they will need to be adapted to the new name!",
			Required:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "Rules",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Rules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of this setting when referred to by the Config REST API V1",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description": me.Description,
		"rules":       me.Rules,
		"name":        me.Name,
		"legacy_id":   me.LegacyID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description": &me.Description,
		"rules":       &me.Rules,
		"name":        &me.Name,
		"legacy_id":   &me.LegacyID,
	})
}
