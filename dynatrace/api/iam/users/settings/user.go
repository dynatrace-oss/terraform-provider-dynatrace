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

type User struct {
	Email  string   `json:"email"`
	UID    string   `json:"uid"`
	Groups []string `json:"-"`
}

func (me *User) Name() string {
	return me.Email
}

func (me *User) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"groups": {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"uid": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (me *User) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"email":  me.Email,
		"groups": me.Groups,
		"uid":    me.UID,
	})
}

func (me *User) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"email":  &me.Email,
		"groups": &me.Groups,
		"uid":    &me.UID,
	})
}
