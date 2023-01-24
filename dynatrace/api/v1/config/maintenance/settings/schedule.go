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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Schedule The schedule of the maintenance window.
type Schedule struct {
	Start          string                     `json:"start"`                // The start date and time of the maintenance window validity period in yyyy-mm-dd HH:mm format.
	End            string                     `json:"end"`                  // The end date and time of the maintenance window validity period in yyyy-mm-dd HH:mm format.
	ZoneID         string                     `json:"zoneId"`               // The time zone of the start and end time. Default time zone is UTC.  You can use either UTC offset `UTC+01:00` format or the IANA Time Zone Database format (for example, `Europe/Vienna`).
	Recurrence     *Recurrence                `json:"recurrence,omitempty"` // The recurrence of the maintenance window.
	RecurrenceType RecurrenceType             `json:"recurrenceType"`       // The type of the schedule recurrence.
	Unknowns       map[string]json.RawMessage `json:"-"`
}

func (me *Schedule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start": {
			Type:        schema.TypeString,
			Description: "The start date and time of the maintenance window validity period in yyyy-mm-dd HH:mm format",
			Required:    true,
		},
		"end": {
			Type:        schema.TypeString,
			Description: "The end date and time of the maintenance window validity period in yyyy-mm-dd HH:mm format",
			Required:    true,
		},
		"zone_id": {
			Type:        schema.TypeString,
			Description: "The time zone of the start and end time. Default time zone is UTC. You can use either UTC offset `UTC+01:00` format or the IANA Time Zone Database format (for example, `Europe/Vienna`)",
			Required:    true,
		},
		"recurrence": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The recurrence of the maintenance window",
			Elem: &schema.Resource{
				Schema: new(Recurrence).Schema(),
			},
		},
		"recurrence_type": {
			Type:        schema.TypeString,
			Description: "The type of the schedule recurrence. Possible values are `DAILY`, `MONTHLY`, `ONCE` and `WEEKLY`",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Schedule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("start", me.Start); err != nil {
		return err
	}
	if err := properties.Encode("end", me.End); err != nil {
		return err
	}
	if err := properties.Encode("zone_id", me.ZoneID); err != nil {
		return err
	}
	if err := properties.Encode("recurrence", me.Recurrence); err != nil {
		return err
	}
	if err := properties.Encode("recurrence_type", string(me.RecurrenceType)); err != nil {
		return err
	}
	return nil
}

func (me *Schedule) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "start")
		delete(me.Unknowns, "end")
		delete(me.Unknowns, "zone_id")
		delete(me.Unknowns, "recurrence")
		delete(me.Unknowns, "recurrence_type")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("start"); ok {
		me.Start = value.(string)
	}
	if value, ok := decoder.GetOk("end"); ok {
		me.End = value.(string)
	}
	if value, ok := decoder.GetOk("zone_id"); ok {
		me.ZoneID = value.(string)
	}
	if _, ok := decoder.GetOk("recurrence.#"); ok {
		me.Recurrence = new(Recurrence)
		if err := me.Recurrence.UnmarshalHCL(hcl.NewDecoder(decoder, "recurrence", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("recurrence_type"); ok {
		me.RecurrenceType = RecurrenceType(value.(string))
	}
	return nil
}

func (me *Schedule) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("start", me.Start); err != nil {
		return nil, err
	}
	if err := m.Marshal("end", me.End); err != nil {
		return nil, err
	}
	if err := m.Marshal("zoneId", me.ZoneID); err != nil {
		return nil, err
	}
	if err := m.Marshal("recurrence", me.Recurrence); err != nil {
		return nil, err
	}
	if err := m.Marshal("recurrenceType", me.RecurrenceType); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *Schedule) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("start", &me.Start); err != nil {
		return err
	}
	if err := m.Unmarshal("end", &me.End); err != nil {
		return err
	}
	if err := m.Unmarshal("zoneId", &me.ZoneID); err != nil {
		return err
	}
	if err := m.Unmarshal("recurrence", &me.Recurrence); err != nil {
		return err
	}
	if err := m.Unmarshal("recurrenceType", &me.RecurrenceType); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// RecurrenceType The type of the schedule recurrence.
type RecurrenceType string

// RecurrenceTypes offers the known enum values
var RecurrenceTypes = struct {
	Daily   RecurrenceType
	Monthly RecurrenceType
	Once    RecurrenceType
	Weekly  RecurrenceType
}{
	"DAILY",
	"MONTHLY",
	"ONCE",
	"WEEKLY",
}
