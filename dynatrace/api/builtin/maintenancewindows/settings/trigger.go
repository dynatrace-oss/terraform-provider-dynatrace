/**
* @license
* Copyright 2026 Dynatrace LLC
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

package maintenancewindows

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Trigger. Base trigger definition.
type Trigger struct {
	Once      *OnceTrigger `json:"once,omitempty"`      // Once Trigger
	Recurring *TimeTrigger `json:"recurring,omitempty"` // Time trigger
	Type      TriggerType  `json:"type"`                // Type of trigger. Possible values: `once`, `time`
}

func (me *Trigger) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"once": {
			Type:        schema.TypeList,
			Description: "Once Trigger",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(OnceTrigger).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"recurring": {
			Type:        schema.TypeList,
			Description: "Time trigger",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(TimeTrigger).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Type of trigger. Possible values: `once`, `time`",
			Required:    true,
		},
	}
}

func (me *Trigger) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"once":      me.Once,
		"recurring": me.Recurring,
		"type":      me.Type,
	})
}

func (me *Trigger) HandlePreconditions() error {
	if (me.Once != nil) && (string(me.Type) != "once") {
		return fmt.Errorf("'once' must not be specified unless 'type' is set to 'once'; got 'type'='%v'", me.Type)
	}
	if (me.Once == nil) && (string(me.Type) == "once") {
		return fmt.Errorf("'once' must be specified when 'type' is set to 'once'; got 'type'='%v'", me.Type)
	}
	if (me.Recurring != nil) && (string(me.Type) != "time") {
		return fmt.Errorf("'recurring' must not be specified unless 'type' is set to 'time'; got 'type'='%v'", me.Type)
	}
	if (me.Recurring == nil) && (string(me.Type) == "time") {
		return fmt.Errorf("'recurring' must be specified when 'type' is set to 'time'; got 'type'='%v'", me.Type)
	}
	return nil
}

func (me *Trigger) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"once":      &me.Once,
		"recurring": &me.Recurring,
		"type":      &me.Type,
	})
}
