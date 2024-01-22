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

package preferences

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	PersistenceMode Options `json:"persistenceMode"` // Possible Values: `ALLOW_ALL_ATTRIBUTES`, `BLOCK_ALL_ATTRIBUTES`
}

func (me *Settings) Name() string {
	return "attributes_preferences"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"persistence_mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ALLOW_ALL_ATTRIBUTES`, `BLOCK_ALL_ATTRIBUTES`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"persistence_mode": me.PersistenceMode,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"persistence_mode": &me.PersistenceMode,
	})
}
