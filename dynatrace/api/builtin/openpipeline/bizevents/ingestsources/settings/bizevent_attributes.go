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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type BizeventAttributes struct {
	EventProvider   *GenericValueAssignment `json:"eventProvider"`       // Event provider
	EventType       *GenericValueAssignment `json:"eventType,omitempty"` // Event type
	FieldExtraction *FieldExtraction        `json:"fieldExtraction"`     // Field extraction
}

func (me *BizeventAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_provider": {
			Type:        schema.TypeList,
			Description: "Event provider",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(GenericValueAssignment).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"event_type": {
			Type:        schema.TypeList,
			Description: "Event type",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(GenericValueAssignment).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"field_extraction": {
			Type:        schema.TypeList,
			Description: "Field extraction",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FieldExtraction).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *BizeventAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_provider":   me.EventProvider,
		"event_type":       me.EventType,
		"field_extraction": me.FieldExtraction,
	})
}

func (me *BizeventAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_provider":   &me.EventProvider,
		"event_type":       &me.EventType,
		"field_extraction": &me.FieldExtraction,
	})
}
