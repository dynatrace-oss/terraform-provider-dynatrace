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

package serviceusers

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceUser struct {
	Email       string   `json:"email,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Groups      []string `json:"-"`
}

func (me *ServiceUser) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the service user",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the service user",
		},
		"groups": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The UUIDs of the groups the service user belongs to",
		},
		"email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The email of the service user",
		},
	}
}

func (me *ServiceUser) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":        me.Name,
		"description": me.Description,
		"groups":      me.Groups,
		"email":       me.Email,
	})
}

func (me *ServiceUser) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":        &me.Name,
		"description": &me.Description,
		"groups":      &me.Groups,
		"email":       &me.Email,
	})
}
