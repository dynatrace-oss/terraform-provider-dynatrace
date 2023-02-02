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

package custom

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ErrorRateIncrease  *ErrorRateIncrease  `json:"errorRateIncrease"`  // Error rate increase
	Scope              *string             `json:"-" scope:"scope"`    // The scope of this setting (DEVICE_APPLICATION_METHOD CUSTOM_APPLICATION environment)
	SlowUserActions    *SlowUserActions    `json:"slowUserActions"`    // Slow user actions
	UnexpectedHighLoad *UnexpectedHighLoad `json:"unexpectedHighLoad"` // Unexpected high load
	UnexpectedLowLoad  *UnexpectedLowLoad  `json:"unexpectedLowLoad"`  // Unexpected low load
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_rate_increase": {
			Type:        schema.TypeList,
			Description: "Error rate increase",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(ErrorRateIncrease).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (DEVICE_APPLICATION_METHOD CUSTOM_APPLICATION environment)",
			Optional:    true,
			Default:     "environment",
		},
		"slow_user_actions": {
			Type:        schema.TypeList,
			Description: "Slow user actions",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SlowUserActions).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"unexpected_high_load": {
			Type:        schema.TypeList,
			Description: "Unexpected high load",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(UnexpectedHighLoad).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"unexpected_low_load": {
			Type:        schema.TypeList,
			Description: "Unexpected low load",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(UnexpectedLowLoad).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"error_rate_increase":  me.ErrorRateIncrease,
		"scope":                me.Scope,
		"slow_user_actions":    me.SlowUserActions,
		"unexpected_high_load": me.UnexpectedHighLoad,
		"unexpected_low_load":  me.UnexpectedLowLoad,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"error_rate_increase":  &me.ErrorRateIncrease,
		"scope":                &me.Scope,
		"slow_user_actions":    &me.SlowUserActions,
		"unexpected_high_load": &me.UnexpectedHighLoad,
		"unexpected_low_load":  &me.UnexpectedLowLoad,
	})
}
