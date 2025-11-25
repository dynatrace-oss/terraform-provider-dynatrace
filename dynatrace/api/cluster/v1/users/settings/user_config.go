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

package users

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The configuration of the user
type UserConfig struct {
	UserName  string   `json:"id"`               // User ID
	Email     string   `json:"email"`            // User's email address
	FirstName string   `json:"firstName"`        // User's first name
	LastName  string   `json:"lastName"`         // User's last name
	Groups    []string `json:"groups,omitempty"` // List of user's user group IDs
}

func (me *UserConfig) Name() string {
	return me.UserName
}

func (me *UserConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The User Name",
		},
		"email": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User's email address",
		},
		"first_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User's first name",
		},
		"last_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User's last name",
		},
		"groups": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "List of user's user group IDs",
		},
	}
}

func (me *UserConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"user_name":  me.UserName,
		"email":      me.Email,
		"first_name": me.FirstName,
		"last_name":  me.LastName,
		"groups":     me.Groups,
	})
}

func (me *UserConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"user_name":  &me.UserName,
		"email":      &me.Email,
		"first_name": &me.FirstName,
		"last_name":  &me.LastName,
		"groups":     &me.Groups,
	})
}
