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

type MetadataFilter struct {
	MetadataFilterItems MetadataFilterItems `json:"metadataFilterItems"` // Define filters for event properties. A maximum of 20 properties is allowed.
}

func (me *MetadataFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"items": {
			Type:        schema.TypeList,
			Description: "Define filters for event properties. A maximum of 20 properties is allowed.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(MetadataFilterItems).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *MetadataFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"items": me.MetadataFilterItems,
	})
}

func (me *MetadataFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"items": &me.MetadataFilterItems,
	})
}
