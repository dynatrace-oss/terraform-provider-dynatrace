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

package remoteaccess

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The configuration of the user
type Settings struct {
	ID            *string `json:"requestId"`     // Request id
	UserId        string  `json:"userId"`        // User id, cannot be changed once created
	Reason        string  `json:"reason"`        // Request reason description, cannot be changed once created
	RequestedDays int     `json:"requestedDays"` // For how many days access is requested, cannot be changed once created
	Role          Role    `json:"role"`          // Requested role, cannot be changed once created
	State         *State  `json:"state"`         // Access request state. Automatically set as `ACCEPTED` on create, state can be changed in subsequent updates.
}

type UpdateSettings struct {
	State State `json:"state"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_id": {
			Type:             schema.TypeString,
			Description:      "User id, cannot be changed once created",
			Required:         true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return d.Id() != "" },
		},
		"reason": {
			Type:             schema.TypeString,
			Description:      "Request reason description, cannot be changed once created",
			Required:         true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return d.Id() != "" },
		},
		"requested_days": {
			Type:             schema.TypeInt,
			Description:      "For how many days access is requested, cannot be changed once created",
			Required:         true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return d.Id() != "" },
		},
		"role": {
			Type:             schema.TypeString,
			Description:      "Requested role, cannot be changed once created",
			Required:         true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return d.Id() != "" },
		},
		"state": {
			Type:             schema.TypeString,
			Description:      "Access request state. Automatically set as `ACCEPTED` on create, state can be changed in subsequent updates.",
			Optional:         true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return newValue == "" },
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"user_id":        me.UserId,
		"reason":         me.Reason,
		"requested_days": me.RequestedDays,
		"role":           me.Role,
		"state":          me.State,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"user_id":        &me.UserId,
		"reason":         &me.Reason,
		"requested_days": &me.RequestedDays,
		"role":           &me.Role,
		"state":          &me.State,
	})
}
