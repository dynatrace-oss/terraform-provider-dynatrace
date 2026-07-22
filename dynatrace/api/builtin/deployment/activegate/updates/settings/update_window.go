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

package updates

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type UpdateWindows []*UpdateWindow

func (me *UpdateWindows) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"update_window": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(UpdateWindow).Schema()},
		},
	}
}

func (me UpdateWindows) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("update_window", me)
}

func (me *UpdateWindows) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("update_window", me)
}

type UpdateWindow struct {
	UpdateWindow string `json:"updateWindow"` // Select an [update window for ActiveGate updates](/ui/settings/builtin:deployment.management.update-windows)
}

func (me *UpdateWindow) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"update_window": {
			Type:        schema.TypeString,
			Description: "Select an [update window for ActiveGate updates](/ui/settings/builtin:deployment.management.update-windows)",
			Required:    true,
		},
	}
}

func (me *UpdateWindow) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"update_window": me.UpdateWindow,
	})
}

func (me *UpdateWindow) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"update_window": &me.UpdateWindow,
	})
}
