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

type EventTrigger struct {
	Active               bool                `json:"isActive"`             // The trigger is active (`true`) or not (`false`)
	TriggerConfiguration *EventTriggerConfig `json:"triggerConfiguration"` //

	// FilterQuery          string              `json:"filterQuery" flags:"readonly"`      //
	// UniqueExpression     string              `json:"uniqueExpression" flags:"readonly"` //
}

func (me *EventTrigger) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"active": {
			Type:        schema.TypeBool,
			Description: "If specified the workflow is getting triggered based on a schedule",
			Optional:    true,
		},
		"config": {
			Type:        schema.TypeList,
			Description: "If specified the workflow is getting triggered based on events",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(EventTriggerConfig).Schema(prefix + ".0.config")},
		},
	}
}

func (me *EventTrigger) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"active": me.Active,
		"config": me.TriggerConfiguration,
	})
}

func (me *EventTrigger) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"active": &me.Active,
		"config": &me.TriggerConfiguration,
	})
}
