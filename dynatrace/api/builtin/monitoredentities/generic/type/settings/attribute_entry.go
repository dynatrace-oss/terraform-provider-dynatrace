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

package generictype

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AttributeEntries []*AttributeEntry

func (me *AttributeEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(AttributeEntry).Schema()},
		},
	}
}

func (me AttributeEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("attribute", me)
}

func (me *AttributeEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("attribute", me)
}

// Attribute entry. Describe how an attribute is extracted from ingest data.
type AttributeEntry struct {
	DisplayName *string `json:"displayName,omitempty"` // The human readable attribute name for this extraction rule. Leave blank to use the key as the display name.
	Key         string  `json:"key"`                   // The attribute key is the unique name of the attribute.
	Pattern     string  `json:"pattern"`               // Pattern for specifying the value for the extracted attribute. Can be a static value, placeholders or a combination of both.
}

func (me *AttributeEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Description: "The human readable attribute name for this extraction rule. Leave blank to use the key as the display name.",
			Optional:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The attribute key is the unique name of the attribute.",
			Required:    true,
		},
		"pattern": {
			Type:        schema.TypeString,
			Description: "Pattern for specifying the value for the extracted attribute. Can be a static value, placeholders or a combination of both.",
			Required:    true,
		},
	}
}

func (me *AttributeEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"display_name": me.DisplayName,
		"key":          me.Key,
		"pattern":      me.Pattern,
	})
}

func (me *AttributeEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"display_name": &me.DisplayName,
		"key":          &me.Key,
		"pattern":      &me.Pattern,
	})
}
