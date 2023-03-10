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

package osservicesmonitoring

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetadataItems []*MetadataItem

func (me *MetadataItems) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"item": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(MetadataItem).Schema()},
		},
	}
}

func (me MetadataItems) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("item", me)
}

func (me *MetadataItems) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("item", me)
}

type MetadataItem struct {
	MetadataKey   string `json:"metadataKey"` // Type 'dt.' for key hints.
	MetadataValue string `json:"metadataValue"`
}

func (me *MetadataItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata_key": {
			Type:        schema.TypeString,
			Description: "Type 'dt.' for key hints.",
			Required:    true,
		},
		"metadata_value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *MetadataItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"metadata_key":   me.MetadataKey,
		"metadata_value": me.MetadataValue,
	})
}

func (me *MetadataItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"metadata_key":   &me.MetadataKey,
		"metadata_value": &me.MetadataValue,
	})
}
