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

package resourcetypes

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	PrimaryResourceType   PrimaryResourceType `json:"primaryResourceType"`             // Possible Values: `CSS`, `IMAGE`, `OTHER`, `SCRIPT`
	RegularExpression     string              `json:"regularExpression"`               // The regular expression to detect the resource.
	SecondaryResourceType *string             `json:"secondaryResourceType,omitempty"` // The secondary type of the resource.
	InsertAfter           string              `json:"-"`
}

func (me *Settings) Name() string {
	return "resource_types"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"primary_resource_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CSS`, `IMAGE`, `OTHER`, `SCRIPT`",
			Required:    true,
		},
		"regular_expression": {
			Type:        schema.TypeString,
			Description: "The regular expression to detect the resource.",
			Required:    true,
		},
		"secondary_resource_type": {
			Type:        schema.TypeString,
			Description: "The secondary type of the resource.",
			Optional:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"primary_resource_type":   me.PrimaryResourceType,
		"regular_expression":      me.RegularExpression,
		"secondary_resource_type": me.SecondaryResourceType,
		"insert_after":            me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"primary_resource_type":   &me.PrimaryResourceType,
		"regular_expression":      &me.RegularExpression,
		"secondary_resource_type": &me.SecondaryResourceType,
		"insert_after":            &me.InsertAfter,
	})
}
