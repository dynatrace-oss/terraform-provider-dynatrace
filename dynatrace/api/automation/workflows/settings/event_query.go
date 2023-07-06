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

package workflows

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventQuery struct {
	Query     string    `json:"query" minlength:"1" maxlength:"800"` //
	EventType EventType `json:"eventType"`                           //
}

func (me *EventQuery) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"query": {
			Type:             schema.TypeString,
			Description:      "A query based on DQL for events that trigger executions",
			Required:         true,
			ValidateDiagFunc: Validate(ValidateMinLength(1), ValidateMaxLength(800)),
		},
		"event_type": {
			Type:        schema.TypeString,
			Description: "Possible values: `events` or `bizevents`. Default: `events`",
			Optional:    true,
			Default:     "events",
		},
	}
}

func (me *EventQuery) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"query":      me.Query,
		"event_type": me.EventType,
	})
}

func (me *EventQuery) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"query":      &me.Query,
		"event_type": &me.EventType,
	})
}
