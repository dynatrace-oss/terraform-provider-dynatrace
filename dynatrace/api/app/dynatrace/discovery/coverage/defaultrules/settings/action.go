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

package defaultrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Actions []*Action

func (me *Actions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Action).Schema()},
		},
	}
}

func (me Actions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("action", me)
}

func (me *Actions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("action", me)
}

type Action struct {
	InstantAction *bool            `json:"instantAction,omitempty"` // Instant action
	Name          string           `json:"name"`
	Parameters    ActionParameters `json:"parameters,omitempty"`
}

func (me *Action) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instant_action": {
			Type:        schema.TypeBool,
			Description: "Instant action",
			Optional:    true, // nullable
		},
		"name": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"parameters": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(ActionParameters).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Action) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"instant_action": me.InstantAction,
		"name":           me.Name,
		"parameters":     me.Parameters,
	})
}

func (me *Action) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"instant_action": &me.InstantAction,
		"name":           &me.Name,
		"parameters":     &me.Parameters,
	})
}
