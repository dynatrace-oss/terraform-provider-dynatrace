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

package groups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Group struct {
	Name                     string      `json:"name"`
	Description              string      `json:"description"`
	FederatedAttributeValues []string    `json:"federatedAttributeValues"`
	Permissions              Permissions `json:"-"`
}

func (me *Group) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"federated_attribute_values": {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"permissions": {
			Type:       schema.TypeList,
			Deprecated: "Assigning permissions directly when creating a group is deprecated. Use the resource `dynatrace_iam_permission` instead.",
			Optional:   true,
			MinItems:   1,
			MaxItems:   1,
			Elem:       &schema.Resource{Schema: new(Permissions).Schema()},
		},
	}
}

func (me *Group) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":                       me.Name,
		"description":                me.Description,
		"federated_attribute_values": me.FederatedAttributeValues,
		"permissions":                me.Permissions,
	})
}

func (me *Group) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":                       &me.Name,
		"description":                &me.Description,
		"federated_attribute_values": &me.FederatedAttributeValues,
		"permissions":                &me.Permissions,
	})
}
