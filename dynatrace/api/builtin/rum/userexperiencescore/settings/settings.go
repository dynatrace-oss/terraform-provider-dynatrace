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

package userexperiencescore

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ConsiderLastAction                bool `json:"considerLastAction"`                // If last user action in a session is classified as Frustrating, classify the entire session as Frustrating
	ConsiderRageClick                 bool `json:"considerRageClick"`                 // Consider rage clicks / rage taps in score calculation
	MaxFrustratedUserActionsThreshold int  `json:"maxFrustratedUserActionsThreshold"` // User experience is considered Frustrating when the selected percentage or more of the user actions in a session are rated as Frustrating.
	MinSatisfiedUserActionsThreshold  int  `json:"minSatisfiedUserActionsThreshold"`  // User experience is considered Satisfying when at least the selected percentage of the user actions in a session are rated as Satisfying.
}

func (me *Settings) Name() string {
	return "user_experience_score"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"consider_last_action": {
			Type:        schema.TypeBool,
			Description: "If last user action in a session is classified as Frustrating, classify the entire session as Frustrating",
			Required:    true,
		},
		"consider_rage_click": {
			Type:        schema.TypeBool,
			Description: "Consider rage clicks / rage taps in score calculation",
			Required:    true,
		},
		"max_frustrated_user_actions_threshold": {
			Type:        schema.TypeInt,
			Description: "User experience is considered Frustrating when the selected percentage or more of the user actions in a session are rated as Frustrating.",
			Required:    true,
		},
		"min_satisfied_user_actions_threshold": {
			Type:        schema.TypeInt,
			Description: "User experience is considered Satisfying when at least the selected percentage of the user actions in a session are rated as Satisfying.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"consider_last_action":                  me.ConsiderLastAction,
		"consider_rage_click":                   me.ConsiderRageClick,
		"max_frustrated_user_actions_threshold": me.MaxFrustratedUserActionsThreshold,
		"min_satisfied_user_actions_threshold":  me.MinSatisfiedUserActionsThreshold,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"consider_last_action":                  &me.ConsiderLastAction,
		"consider_rage_click":                   &me.ConsiderRageClick,
		"max_frustrated_user_actions_threshold": &me.MaxFrustratedUserActionsThreshold,
		"min_satisfied_user_actions_threshold":  &me.MinSatisfiedUserActionsThreshold,
	})
}
