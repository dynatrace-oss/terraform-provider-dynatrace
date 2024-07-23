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

package vault

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AllowedEntities []*CredentialAccessData

func (me *AllowedEntities) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "	The set of entities allowed to use the credential.",
			Elem:        &schema.Resource{Schema: new(CredentialAccessData).Schema()},
		},
	}
}

func (me AllowedEntities) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("entity", me)
}

func (me *AllowedEntities) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("entity", me)
}

type CredentialAccessData struct {
	Type EntityType `json:"type,omitempty"` // Type of entity. Possible values: `USER`, `APPLICATION`, `UNKNOWN`
	Id   *string    `json:"id,omitempty"`   // ID of the entity
}

func (me *CredentialAccessData) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Type of entity. Possible values: `USER`, `APPLICATION`, `UNKNOWN`",
			Optional:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "ID of the entity",
			Optional:    true,
		},
	}
}

func (me *CredentialAccessData) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type": me.Type,
		"id":   me.Id,
	})
}

func (me *CredentialAccessData) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type": &me.Type,
		"id":   &me.Id,
	})
}
