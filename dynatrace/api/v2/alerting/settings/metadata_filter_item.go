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

package alerting

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MetadataFilterItems []*MetadataFilterItem

func (me *MetadataFilterItems) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(MetadataFilterItem).Schema()},
		},
	}
}

func (me MetadataFilterItems) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("filter", me)
}

func (me *MetadataFilterItems) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("filter", me)
}

type MetadataFilterItem struct {
	MetadataKey   string `json:"metadataKey"`   // Type 'dt.' for key hints.
	MetadataValue string `json:"metadataValue"` // Value
}

func (me *MetadataFilterItem) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "Type 'dt.' for key hints.",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Value",
			Required:    true,
		},
	}
}

func (me *MetadataFilterItem) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":   me.MetadataKey,
		"value": me.MetadataValue,
	})
}

func (me *MetadataFilterItem) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":   &me.MetadataKey,
		"value": &me.MetadataValue,
	})
}
