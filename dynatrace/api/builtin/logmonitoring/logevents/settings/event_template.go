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

package logevents

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventTemplate struct {
	DavisMerge  *bool         `json:"davisMerge,omitempty"` // Davis® AI will try to merge this event into existing problems, otherwise a new problem will always be created.
	Description string        `json:"description"`          // The description of the event to trigger.
	EventType   EventTypeEnum `json:"eventType"`            // Possible Values: `AVAILABILITY`, `CUSTOM_ALERT`, `CUSTOM_ANNOTATION`, `CUSTOM_CONFIGURATION`, `CUSTOM_DEPLOYMENT`, `ERROR`, `INFO`, `MARKED_FOR_TERMINATION`, `RESOURCE`, `SLOWDOWN`, `WARNING`
	Metadata    MetadataItems `json:"metadata,omitempty"`   // Set of additional key-value properties to be attached to the triggered event.
	Title       string        `json:"title"`                // The title of the event to trigger.
}

func (me *EventTemplate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"davis_merge": {
			Type:        schema.TypeBool,
			Description: "Davis® AI will try to merge this event into existing problems, otherwise a new problem will always be created.",
			Optional:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the event to trigger.",
			Required:    true,
		},
		"event_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `AVAILABILITY`, `CUSTOM_ALERT`, `CUSTOM_ANNOTATION`, `CUSTOM_CONFIGURATION`, `CUSTOM_DEPLOYMENT`, `ERROR`, `INFO`, `MARKED_FOR_TERMINATION`, `RESOURCE`, `SLOWDOWN`, `WARNING`",
			Required:    true,
		},
		"metadata": {
			Type:        schema.TypeList,
			Description: "Set of additional key-value properties to be attached to the triggered event.",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(MetadataItems).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"title": {
			Type:        schema.TypeString,
			Description: "The title of the event to trigger.",
			Required:    true,
		},
	}
}

func (me *EventTemplate) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"davis_merge": me.DavisMerge,
		"description": me.Description,
		"event_type":  me.EventType,
		"metadata":    me.Metadata,
		"title":       me.Title,
	})
}

func (me *EventTemplate) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"davis_merge": &me.DavisMerge,
		"description": &me.Description,
		"event_type":  &me.EventType,
		"metadata":    &me.Metadata,
		"title":       &me.Title,
	})
	if me.DavisMerge == nil && me.EventType != EventTypeEnums.Info {
		me.DavisMerge = opt.NewBool(false)
	}
	return err
}
