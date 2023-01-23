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

package maintenance

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Window Configuration of a maintenance window.
type Window struct {
	Name                               string      `json:"name"`                                         // The name of the maintenance window, displayed in the UI.
	Description                        string      `json:"description"`                                  // A short description of the maintenance purpose.
	Schedule                           *Schedule   `json:"schedule"`                                     // The schedule of the maintenance window.
	Scope                              *Scope      `json:"scope,omitempty"`                              // The scope of the maintenance window.   The scope restricts the alert/problem detection suppression to certain Dynatrace entities. It can contain a list of entities and/or matching rules for dynamic formation of the scope.   If no scope is specified, the alert/problem detection suppression applies to the entire environment.
	Suppression                        Suppression `json:"suppression"`                                  // The type of suppression of alerting and problem detection during the maintenance.
	Type                               WindowType  `json:"type"`                                         // The type of the maintenance: planned or unplanned
	SuppressSyntheticMonitorsExecution *bool       `json:"suppressSyntheticMonitorsExecution,omitempty"` // Suppress execution of synthetic monitors during the maintenance
	Enabled                            bool        `json:"enabled"`

	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *Window) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the maintenance window, displayed in the UI",
			Required:    true,
		},
		"schedule": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The schedule of the maintenance window",
			Elem: &schema.Resource{
				Schema: new(Schedule).Schema(),
			},
		},
		"scope": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "the tiles this Dashboard consist of",
			Elem: &schema.Resource{
				Schema: new(Scope).Schema(),
			},
		},
		"suppression": {
			Type:        schema.TypeString,
			Description: "The type of suppression of alerting and problem detection during the maintenance",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the maintenance: planned or unplanned",
			Required:    true,
		},
		"suppress_synth_mon_exec": {
			Type:        schema.TypeBool,
			Description: "Suppress execution of synthetic monitors during the maintenance",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The Maintenance Window is enabled or disabled",
			Default:     true,
			Optional:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the maintenance purpose",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Window) MarshalHCL(properties hcl.Properties) error {
	if me.Unknowns != nil {
		delete(me.Unknowns, "managementZoneId")
	}
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("description", me.Description); err != nil {
		return err
	}
	if err := properties.Encode("schedule", me.Schedule); err != nil {
		return err
	}
	if err := properties.Encode("suppress_synth_mon_exec", opt.Bool(me.SuppressSyntheticMonitorsExecution)); err != nil {
		return err
	}
	if !me.Enabled {
		if err := properties.Encode("enabled", me.Enabled); err != nil {
			return err
		}
	}
	if me.Scope != nil && !me.Scope.IsEmpty() {
		if err := properties.Encode("scope", me.Scope); err != nil {
			return err
		}
	}
	if err := properties.Encode("suppression", string(me.Suppression)); err != nil {
		return err
	}
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	return nil
}

func (me *Window) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "description")
		delete(me.Unknowns, "schedule")
		delete(me.Unknowns, "scope")
		delete(me.Unknowns, "suppression")
		delete(me.Unknowns, "type")
		delete(me.Unknowns, "metadata")
		delete(me.Unknowns, "suppress_synth_mon_exec")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "id")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("suppress_synth_mon_exec"); ok {
		me.SuppressSyntheticMonitorsExecution = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	} else {
		me.Enabled = true
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("description"); ok {
		me.Description = value.(string)
	}

	if _, ok := decoder.GetOk("schedule.#"); ok {
		me.Schedule = new(Schedule)
		if err := me.Schedule.UnmarshalHCL(hcl.NewDecoder(decoder, "schedule", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("scope.#"); ok {
		me.Scope = new(Scope)
		if err := me.Scope.UnmarshalHCL(hcl.NewDecoder(decoder, "scope", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("suppression"); ok {
		me.Suppression = Suppression(value.(string))
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = WindowType(value.(string))
	}
	return nil
}

func (me *Window) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("suppressSyntheticMonitorsExecution", me.SuppressSyntheticMonitorsExecution); err != nil {
		return nil, err
	}
	if err := m.Marshal("enabled", me.Enabled); err != nil {
		return nil, err
	}
	if err := m.Marshal("name", me.Name); err != nil {
		return nil, err
	}
	if err := m.Marshal("description", me.Description); err != nil {
		return nil, err
	}
	if err := m.Marshal("schedule", me.Schedule); err != nil {
		return nil, err
	}
	if me.Scope != nil && !me.Scope.IsEmpty() {
		if err := m.Marshal("scope", me.Scope); err != nil {
			return nil, err
		}
	}
	if err := m.Marshal("suppression", me.Suppression); err != nil {
		return nil, err
	}
	if err := m.Marshal("type", me.Type); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *Window) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "metadata")
	delete(m, "id")
	if err := m.Unmarshal("suppressSyntheticMonitorsExecution", &me.SuppressSyntheticMonitorsExecution); err != nil {
		return err
	}
	if err := m.Unmarshal("enabled", &me.Enabled); err != nil {
		return err
	}
	if err := m.Unmarshal("name", &me.Name); err != nil {
		return err
	}
	if err := m.Unmarshal("description", &me.Description); err != nil {
		return err
	}
	if err := m.Unmarshal("schedule", &me.Schedule); err != nil {
		return err
	}
	if err := m.Unmarshal("scope", &me.Scope); err != nil {
		return err
	}
	if err := m.Unmarshal("suppression", &me.Suppression); err != nil {
		return err
	}
	if err := m.Unmarshal("type", &me.Type); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// WindowType The type of the maintenance: planned or unplanned.
type WindowType string

// MaintenanceWindowTypes offers the known enum values
var MaintenanceWindowTypes = struct {
	Planned   WindowType
	Unplanned WindowType
}{
	"PLANNED",
	"UNPLANNED",
}

// Suppression The type of suppression of alerting and problem detection during the maintenance.
type Suppression string

// Suppressions offers the known enum values
var Suppressions = struct {
	DetectProblemsAndAlert  Suppression
	DetectProblemsDontAlert Suppression
	DontDetectProblems      Suppression
}{
	"DETECT_PROBLEMS_AND_ALERT",
	"DETECT_PROBLEMS_DONT_ALERT",
	"DONT_DETECT_PROBLEMS",
}
