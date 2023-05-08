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

type Settings struct {
	OwnershipIdentifiers OwnershipIdentifiers `json:"ownershipIdentifiers"` // Tags and metadata are key-value pairs. Define keys for tags and metadata that are considered for ownership. If a tag or any metadata starts with a key defined below, the value of the tag or metadata is considered a team identifier.
}

func (me *Settings) Name() string {
	return "environment"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ownership_identifiers": {
			Type:        schema.TypeList,
			Description: "Tags and metadata are key-value pairs. Define keys for tags and metadata that are considered for ownership. If a tag or any metadata starts with a key defined below, the value of the tag or metadata is considered a team identifier.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(OwnershipIdentifiers).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"ownership_identifiers": me.OwnershipIdentifiers,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"ownership_identifiers": &me.OwnershipIdentifiers,
	})
}
