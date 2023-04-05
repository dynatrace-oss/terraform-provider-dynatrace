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

package opentelemetrymetrics

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DropAttributeItems []*DropAttributeItem

func (me *DropAttributeItems) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"to_drop_attribute": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DropAttributeItem).Schema()},
		},
	}
}

func (me DropAttributeItems) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("to_drop_attribute", me)
}

func (me *DropAttributeItems) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("to_drop_attribute", me)
}

type DropAttributeItem struct {
	AttributeKey string `json:"attributeKey"` // Attribute key
	Enabled      bool   `json:"enabled"`      // This setting is enabled (`true`) or disabled (`false`)
}

func (me *DropAttributeItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute_key": {
			Type:        schema.TypeString,
			Description: "Attribute key",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
	}
}

func (me *DropAttributeItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attribute_key": me.AttributeKey,
		"enabled":       me.Enabled,
	})
}

func (me *DropAttributeItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attribute_key": &me.AttributeKey,
		"enabled":       &me.Enabled,
	})
}
