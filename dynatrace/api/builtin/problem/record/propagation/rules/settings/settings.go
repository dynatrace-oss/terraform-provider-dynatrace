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

package rules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled            bool   `json:"enabled"`            // This setting is enabled (`true`) or disabled (`false`)
	SourceAttributeKey string `json:"sourceAttributeKey"` // Attribute key from the event that will be propagated.
	TargetAttributeKey string `json:"targetAttributeKey"` // Attribute key under which the propagated event data will be stored on the problem.
}

func (me *Settings) Name() string {
	return "environment"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"source_attribute_key": {
			Type:        schema.TypeString,
			Description: "Attribute key from the event that will be propagated.",
			Required:    true,
		},
		"target_attribute_key": {
			Type:        schema.TypeString,
			Description: "Attribute key under which the propagated event data will be stored on the problem.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":              me.Enabled,
		"source_attribute_key": me.SourceAttributeKey,
		"target_attribute_key": me.TargetAttributeKey,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":              &me.Enabled,
		"source_attribute_key": &me.SourceAttributeKey,
		"target_attribute_key": &me.TargetAttributeKey,
	})
}
