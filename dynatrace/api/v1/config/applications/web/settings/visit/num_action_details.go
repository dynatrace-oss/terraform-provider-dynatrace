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

package visit

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NumActionDetails Configuration of a number of user actions-based conversion goal
type NumActionDetails struct {
	NumUserActions *int32 `json:"numUserActions,omitempty"` // The number of user actions to hit the conversion goal
}

func (me *NumActionDetails) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"num_user_actions": {
			Type:        schema.TypeInt,
			Description: "The number of user actions to hit the conversion goal",
			Optional:    true,
		},
	}
}

func (me *NumActionDetails) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"num_user_actions": me.NumUserActions,
	})
}

func (me *NumActionDetails) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"num_user_actions": &me.NumUserActions,
	})
}
