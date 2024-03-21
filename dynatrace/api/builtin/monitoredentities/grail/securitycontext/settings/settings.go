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

package securitycontext

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DestinationProperty string `json:"destinationProperty"` // The case-sensitive name of a property of the destination type.
	EntityType          string `json:"entityType"`          // Type of the entity whose security context to override.
	InsertAfter         string `json:"-"`
}

func (me *Settings) Name() string {
	return me.EntityType + "_" + me.DestinationProperty
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"destination_property": {
			Type:        schema.TypeString,
			Description: "The case-sensitive name of a property of the destination type.",
			Required:    true,
		},
		"entity_type": {
			Type:        schema.TypeString,
			Description: "Type of the entity whose security context to override.",
			Required:    true,
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
		"destination_property": me.DestinationProperty,
		"entity_type":          me.EntityType,
		"insert_after":         me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"destination_property": &me.DestinationProperty,
		"entity_type":          &me.EntityType,
		"insert_after":         &me.InsertAfter,
	})
}
