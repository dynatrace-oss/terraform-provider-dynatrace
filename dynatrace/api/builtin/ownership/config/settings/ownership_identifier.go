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

package config

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type OwnershipIdentifiers []*OwnershipIdentifier

func (me *OwnershipIdentifiers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ownership_identifier": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(OwnershipIdentifier).Schema()},
		},
	}
}

func (me OwnershipIdentifiers) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("ownership_identifier", me)
}

func (me *OwnershipIdentifiers) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("ownership_identifier", me)
}

type OwnershipIdentifier struct {
	Enabled bool   `json:"enabled"` // This setting is enabled (`true`) or disabled (`false`)
	Key     string `json:"key"`     // Key for ownership metadata and tags
}

func (me *OwnershipIdentifier) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Key for ownership metadata and tags",
			Required:    true,
		},
	}
}

func (me *OwnershipIdentifier) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"key":     me.Key,
	})
}

func (me *OwnershipIdentifier) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled": &me.Enabled,
		"key":     &me.Key,
	})
}
