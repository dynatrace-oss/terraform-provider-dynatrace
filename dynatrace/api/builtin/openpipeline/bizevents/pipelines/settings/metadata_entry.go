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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetadataEntries []*MetadataEntry

func (me *MetadataEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(MetadataEntry).Schema()},
		},
	}
}

func (me MetadataEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("metadata", me)
}

func (me *MetadataEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("metadata", me)
}

type MetadataEntry struct {
	EntryKey   string  `json:"entryKey"`             // Metadata entry key
	EntryValue *string `json:"entryValue,omitempty"` // Metadata entry value
}

func (me *MetadataEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entry_key": {
			Type:        schema.TypeString,
			Description: "Metadata entry key",
			Required:    true,
		},
		"entry_value": {
			Type:        schema.TypeString,
			Description: "Metadata entry value",
			Optional:    true, // nullable
		},
	}
}

func (me *MetadataEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"entry_key":   me.EntryKey,
		"entry_value": me.EntryValue,
	})
}

func (me *MetadataEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"entry_key":   &me.EntryKey,
		"entry_value": &me.EntryValue,
	})
}
