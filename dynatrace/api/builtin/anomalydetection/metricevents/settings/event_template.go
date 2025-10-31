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

package metricevents

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventTemplate struct {
	DavisMerge  *bool           `json:"davisMerge,omitempty"` // Davis® AI will try to merge this event into existing problems, otherwise a new problem will always be created.
	Description string          `json:"description"`          // The description of the event to trigger.
	EventType   EventTypeEnum   `json:"eventType"`            // Possible Values: `AVAILABILITY`, `CUSTOM_ALERT`, `CUSTOM_ANNOTATION`, `CUSTOM_CONFIGURATION`, `CUSTOM_DEPLOYMENT`, `ERROR`, `INFO`, `MARKED_FOR_TERMINATION`, `RESOURCE`, `SLOWDOWN`, `WARNING`
	Metadata    []*MetadataItem `json:"metadata,omitempty"`   // Set of additional key-value properties to be attached to the triggered event.
	Title       string          `json:"title"`                // The title of the event to trigger.
}

func (me *EventTemplate) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"davis_merge": {
			Type:        schema.TypeBool,
			Description: "Davis® AI will try to merge this event into existing problems, otherwise a new problem will always be created.",
			Optional:    true, // precondition
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "The description of the event to trigger.",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressEOT,
		},
		"event_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `AVAILABILITY`, `CUSTOM_ALERT`, `CUSTOM_ANNOTATION`, `CUSTOM_CONFIGURATION`, `CUSTOM_DEPLOYMENT`, `ERROR`, `INFO`, `MARKED_FOR_TERMINATION`, `RESOURCE`, `SLOWDOWN`, `WARNING`",
			Required:    true,
		},
		"metadata": {
			Type:        schema.TypeSet,
			Description: "Set of additional key-value properties to be attached to the triggered event.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(MetadataItem).Schema()},
			MinItems:    1,
		},
		"title": {
			Type:        schema.TypeString,
			Description: "The title of the event to trigger.",
			Required:    true,
		},
	}
}

func (me *EventTemplate) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"title":       me.Title,
		"description": me.Description,
		"event_type":  me.EventType,
		"davis_merge": me.DavisMerge,
	}); err != nil {
		return err
	}
	return properties.EncodeSlice("metadata", me.Metadata)
}

func (me *EventTemplate) HandlePreconditions() error {
	if me.DavisMerge == nil && (string(me.EventType) != "INFO") {
		me.DavisMerge = opt.NewBool(false)
	}
	return nil
}

func (me *EventTemplate) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"title":       &me.Title,
		"description": &me.Description,
		"event_type":  &me.EventType,
		"davis_merge": &me.DavisMerge,
	}); err != nil {
		return err
	}
	return decoder.DecodeSlice("metadata", &me.Metadata)
}
