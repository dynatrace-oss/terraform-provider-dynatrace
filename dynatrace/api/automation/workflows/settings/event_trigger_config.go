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
	"encoding/json"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EventTriggerConfig struct {
	Type EventTriggerType `json:"type"` // `event`, `davis-problem` or `davis-event`

	DavisEventConfig   *DavisEventConfig   `json:"-"` // type = `event`
	DavisProblemConfig *DavisProblemConfig `json:"-"` // type = `davis-problem`
	EventQuery         *EventQuery         `json:"-"` // type = `davis-event`
	GenericConfig      *string             `json:"-"` // type = yet unknown
}

func (me *EventTriggerConfig) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:          schema.TypeString,
			Description:   "The type of the trigger configuration to expect within attribute `value`. Only required if `config` is set. Must not be set if `davis_event`, `davis_problem` or `event` are present",
			Optional:      true,
			RequiredWith:  []string{prefix + ".0.value"},
			ConflictsWith: []string{prefix + ".0.davis_event", prefix + ".0.davis_problem", prefix + ".0.event"},
		},
		"davis_event": {
			Type:          schema.TypeList,
			Description:   "Contains trigger configuration based on Davis Events. Either `davis_event`, `davis_problem`, `davis_event` or `config` need to set",
			MinItems:      1,
			MaxItems:      1,
			Optional:      true,
			Elem:          &schema.Resource{Schema: new(DavisEventConfig).Schema(prefix + ".0.davis_event")},
			ConflictsWith: []string{prefix + ".0.value", prefix + ".0.type"},
			ExactlyOneOf:  []string{prefix + ".0.value", prefix + ".0.davis_problem", prefix + ".0.davis_event", prefix + ".0.event"},
		},
		"davis_problem": {
			Type:          schema.TypeList,
			Description:   "Contains trigger configuration based on Davis Problems. Either `davis_event`, `davis_problem`, `davis_event` or `config` need to set",
			MinItems:      1,
			MaxItems:      1,
			Optional:      true,
			Elem:          &schema.Resource{Schema: new(DavisProblemConfig).Schema(prefix + ".0.davis_problem")},
			ConflictsWith: []string{prefix + ".0.value", prefix + ".0.type"},
			ExactlyOneOf:  []string{prefix + ".0.value", prefix + ".0.davis_problem", prefix + ".0.davis_event", prefix + ".0.event"},
		},
		"event": {
			Type:          schema.TypeList,
			Description:   "Contains trigger configuration based on Davis Problems. Either `davis_event`, `davis_problem`, `davis_event` or `config` need to set",
			MinItems:      1,
			MaxItems:      1,
			Optional:      true,
			Elem:          &schema.Resource{Schema: new(EventQuery).Schema(prefix + ".0.event")},
			ConflictsWith: []string{prefix + ".0.value", prefix + ".0.type"},
			ExactlyOneOf:  []string{prefix + ".0.value", prefix + ".0.davis_problem", prefix + ".0.davis_event", prefix + ".0.event"},
		},
		"value": {
			Type:          schema.TypeString,
			Description:   "Contains JSON encoded trigger configuration if the trigger type is neither `davis_event`, `davis_problem` or `event`. It requires the attribute `type` to be set in combination",
			Optional:      true,
			RequiredWith:  []string{prefix + ".0.type"},
			ConflictsWith: []string{prefix + ".0.davis_event", prefix + ".0.davis_problem", prefix + ".0.event"},
			ExactlyOneOf:  []string{prefix + ".0.value", prefix + ".0.davis_problem", prefix + ".0.davis_event", prefix + ".0.event"},
		},
	}
}

func (me *EventTriggerConfig) MarshalJSON() ([]byte, error) {
	var err error
	m := struct {
		Type   EventTriggerType `json:"type"`
		Config json.RawMessage  `json:"value"`
	}{}
	if me.DavisEventConfig != nil {
		m.Type = EventTriggerTypes.DavisEvent
		if m.Config, err = json.Marshal(me.DavisEventConfig); err != nil {
			return nil, err
		}
	} else if me.DavisProblemConfig != nil {
		m.Type = EventTriggerTypes.DavisProblem
		if m.Config, err = json.Marshal(me.DavisProblemConfig); err != nil {
			return nil, err
		}
	} else if me.EventQuery != nil {
		m.Type = EventTriggerTypes.Event
		if m.Config, err = json.Marshal(me.EventQuery); err != nil {
			return nil, err
		}
	} else {
		m.Type = me.Type
		m.Config = []byte(*me.GenericConfig)
	}
	return json.Marshal(m)
}

func (me *EventTriggerConfig) UnmarshalJSON(data []byte) error {
	m := struct {
		Type   EventTriggerType `json:"type"`
		Config json.RawMessage  `json:"value"`
	}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if len(m.Type) == 0 {
		return errors.New(string(data) + " doesn't contain a property `type`")
	}

	switch m.Type {
	case EventTriggerTypes.Event:
		me.EventQuery = new(EventQuery)
		if err := json.Unmarshal(m.Config, &me.EventQuery); err != nil {
			return err
		}
	case EventTriggerTypes.DavisProblem:
		me.DavisProblemConfig = new(DavisProblemConfig)
		if err := json.Unmarshal(m.Config, &me.DavisProblemConfig); err != nil {
			return err
		}
	case EventTriggerTypes.DavisEvent:
		me.DavisEventConfig = new(DavisEventConfig)
		if err := json.Unmarshal(m.Config, &me.DavisEventConfig); err != nil {
			return err
		}
	default:
		me.GenericConfig = opt.NewString(string(m.Config))
	}

	return nil
}

func (me *EventTriggerConfig) GetType() *EventTriggerType {
	if me.Type.IsKnown() {
		return &me.Type
	}
	return nil
}

func (me *EventTriggerConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":          me.GetType(),
		"davis_event":   me.DavisEventConfig,
		"davis_problem": me.DavisProblemConfig,
		"event":         me.EventQuery,
		"value":         me.GenericConfig,
	})
}

func (me *EventTriggerConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":          &me.Type,
		"davis_event":   &me.DavisEventConfig,
		"davis_problem": &me.DavisProblemConfig,
		"event":         &me.EventQuery,
		"value":         &me.GenericConfig,
	})
}
