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

package rumcustomcrashrateincrease

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CrashRateIncreaseFixed struct {
	AbsoluteCrashRate float64 `json:"absoluteCrashRate"` // Absolute threshold
	ConcurrentUsers   int     `json:"concurrentUsers"`   // Amount of users
}

func (me *CrashRateIncreaseFixed) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"absolute_crash_rate": {
			Type:        schema.TypeFloat,
			Description: "Absolute threshold",
			Required:    true,
		},
		"concurrent_users": {
			Type:        schema.TypeInt,
			Description: "Amount of users",
			Required:    true,
		},
	}
}

func (me *CrashRateIncreaseFixed) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"absolute_crash_rate": me.AbsoluteCrashRate,
		"concurrent_users":    me.ConcurrentUsers,
	})
}

func (me *CrashRateIncreaseFixed) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"absolute_crash_rate": &me.AbsoluteCrashRate,
		"concurrent_users":    &me.ConcurrentUsers,
	})
}
