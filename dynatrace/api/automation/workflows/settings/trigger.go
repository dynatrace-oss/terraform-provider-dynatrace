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

type Trigger struct {
	Schedule     *Schedule     `json:"schedule,omitempty"`     // If specified the workflow is getting triggered based on a schedule
	EventTrigger *EventTrigger `json:"eventTrigger,omitempty"` // If specified the workflow is getting triggered based on events
}

func (me *Trigger) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"schedule": {
			Type:          schema.TypeList,
			Description:   "If specified the workflow is getting triggered based on a schedule",
			MinItems:      1,
			MaxItems:      1,
			Optional:      true,
			Elem:          &schema.Resource{Schema: new(Schedule).Schema(prefix + ".0.schedule")},
			ConflictsWith: []string{prefix + ".0.event"},
		},
		"event": {
			Type:          schema.TypeList,
			Description:   "If specified the workflow is getting triggered based on events",
			MinItems:      1,
			MaxItems:      1,
			Optional:      true,
			Elem:          &schema.Resource{Schema: new(EventTrigger).Schema(prefix + ".0.event")},
			ConflictsWith: []string{prefix + ".0.schedule"},
		},
	}
}

func (me *Trigger) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"schedule": me.Schedule,
		"event":    me.EventTrigger,
	})
}

func (me *Trigger) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"schedule": &me.Schedule,
		"event":    &me.EventTrigger,
	})
}
