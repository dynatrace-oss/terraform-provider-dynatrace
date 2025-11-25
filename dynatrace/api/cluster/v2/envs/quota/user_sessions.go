/**
* @license
* Copyright 2025 Dynatrace LLC
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

package quota

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UserSessions represents user sessions consumption and quota information on environment level. If skipped when editing via PUT method then already set quotas will remain
type UserSessions struct {
	TotalAnnualLimit  *int64 `json:"totalAnnualLimit"`  // Annual total User sessions environment quota. Not set if unlimited. When updating via PUT method, skipping this field will set quota unlimited
	TotalMonthlyLimit *int64 `json:"totalMonthlyLimit"` // Monthly total User sessions environment quota. Not set if unlimited. When updating via PUT method, skipping this field will set quota unlimited
}

func (me *UserSessions) IsEmpty() bool {
	return me == nil || (me.TotalAnnualLimit == nil && me.TotalMonthlyLimit == nil)
}

func (me *UserSessions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"annual": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Annual total User sessions environment quota. Not set if unlimited",
		},
		"monthly": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Monthly total User sessions environment quota. Not set if unlimited",
		},
	}
}

func (me *UserSessions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"monthly": me.TotalMonthlyLimit,
		"annual":  me.TotalAnnualLimit,
	})
}

func (me *UserSessions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"monthly": &me.TotalMonthlyLimit,
		"annual":  &me.TotalAnnualLimit,
	})
}
