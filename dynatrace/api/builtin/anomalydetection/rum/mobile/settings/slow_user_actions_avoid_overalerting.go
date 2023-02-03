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

package rummobile

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SlowUserActionsAvoidOveralerting struct {
	MinActionRate int `json:"minActionRate"`
}

func (me *SlowUserActionsAvoidOveralerting) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"min_action_rate": {
			Type:        schema.TypeInt,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *SlowUserActionsAvoidOveralerting) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"min_action_rate": me.MinActionRate,
	})
}

func (me *SlowUserActionsAvoidOveralerting) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"min_action_rate": &me.MinActionRate,
	})
}
