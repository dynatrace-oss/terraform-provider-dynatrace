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

package rule

import (
	autotag "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AutoTagId    string        `json:"-"`
	Rules        autotag.Rules `json:"rules,omitempty"` // Rules
	CurrentState string        `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_tag_id": {
			Type:        schema.TypeString,
			Description: "Automatically applied tag ID",
			Required:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "Rules",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(autotag.Rules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"current_state": {
			Type:        schema.TypeString,
			Description: "For internal use: current state of rules in JSON format",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) Name() string {
	return "Rules for " + me.AutoTagId
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_tag_id": me.AutoTagId,
		"rules":       me.Rules,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"auto_tag_id": &me.AutoTagId,
		"rules":       &me.Rules,
	})
}
