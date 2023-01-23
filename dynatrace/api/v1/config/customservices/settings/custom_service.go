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

package customservices

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomService has no documentation
type CustomService struct {
	Name                string                     `json:"name"`                          // The name of the custom service, displayed in the UI
	Technology          Technology                 `json:"-"`                             // The technology of the custom service
	Order               *string                    `json:"order,omitempty"`               // The order string. Sorting custom services alphabetically by their order string determines their relative ordering. Typically this is managed by Dynatrace internally and will not be present in GET responses
	Enabled             bool                       `json:"enabled"`                       // Custom service enabled/disabled
	Rules               []*DetectionRule           `json:"rules,omitempty"`               // The list of rules defining the custom service
	QueueEntryPoint     *bool                      `json:"queueEntryPoint"`               // The queue entry point flag. Set to `true` for custom messaging services
	QueueEntryPointType *QueueEntryPointType       `json:"queueEntryPointType,omitempty"` // The queue entry point type
	ProcessGroups       []string                   `json:"processGroups,omitempty"`       // The list of process groups the custom service should belong to
	Unknowns            map[string]json.RawMessage `json:"-"`
}

func (me *CustomService) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the custom service, displayed in the UI",
			Required:    true,
		},
		// "order": {
		// 	Type:        schema.TypeString,
		// 	Description: "The order string. Sorting custom services alphabetically by their order string determines their relative ordering. Typically this is managed by Dynatrace internally and will not be present in GET responses",
		// 	Optional:    true,
		// },
		"technology": {
			Type:        schema.TypeString,
			Description: "Matcher applying to the file name (ENDS_WITH, EQUALS or STARTS_WITH). Default value is ENDS_WITH (if applicable)",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Custom service enabled/disabled",
			Required:    true,
		},
		"queue_entry_point": {
			Type:        schema.TypeBool,
			Description: "The queue entry point flag. Set to `true` for custom messaging services",
			Optional:    true,
		},
		"queue_entry_point_type": {
			Type:        schema.TypeString,
			Description: "The queue entry point type (IBM_MQ, JMS, KAFKA, MSMQ or RABBIT_MQ)",
			Optional:    true,
		},
		"process_groups": {
			Type:        schema.TypeSet,
			Description: "The list of process groups the custom service should belong to",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"rule": {
			Type:        schema.TypeList,
			Description: "The list of rules defining the custom service",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(DetectionRule).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomService) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	// if me.Order != nil {
	// 	if err := properties.Encode("order", me.Order); err != nil { return err }
	// }
	if err := properties.Encode("technology", string(me.Technology)); err != nil {
		return err
	}
	if err := properties.Encode("enabled", me.Enabled); err != nil {
		return err
	}
	if err := properties.Encode("rule", me.Rules); err != nil {
		return err
	}
	if err := properties.Encode("queue_entry_point", opt.Bool(me.QueueEntryPoint)); err != nil {
		return err
	}
	if err := properties.Encode("queue_entry_point_type", me.QueueEntryPointType); err != nil {
		return err
	}
	if err := properties.Encode("process_groups", me.ProcessGroups); err != nil {
		return err
	}
	return nil
}

func (me *CustomService) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "name")
		// delete(me.Unknowns, "order")
		delete(me.Unknowns, "technology")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "queue_entry_point")
		delete(me.Unknowns, "queue_entry_point_type")
		delete(me.Unknowns, "process_groups")
		delete(me.Unknowns, "rule")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if result, ok := decoder.GetOk("rule.#"); ok {
		me.Rules = []*DetectionRule{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(DetectionRule)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "rule", idx)); err != nil {
				return err
			}
			me.Rules = append(me.Rules, entry)
		}
	}
	if value, ok := decoder.GetOk("queue_entry_point"); ok {
		me.QueueEntryPoint = opt.NewBool(value.(bool))
	}
	if me.QueueEntryPoint == nil {
		me.QueueEntryPoint = opt.NewBool(false)
	}
	if value, ok := decoder.GetOk("queue_entry_point_type"); ok && len(value.(string)) > 0 {
		me.QueueEntryPointType = QueueEntryPointType(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("technology"); ok && len(value.(string)) > 0 {
		me.Technology = Technology(value.(string))
	}
	if err := decoder.Decode("process_groups", &me.ProcessGroups); err != nil {
		return err
	}
	return nil
}

func (me *CustomService) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	delete(m, "metadata")
	delete(m, "id")
	delete(m, "technology")
	if err := m.Marshal("name", me.Name); err != nil {
		return nil, err
	}
	if err := m.Marshal("order", me.Order); err != nil {
		return nil, err
	}
	if err := m.Marshal("enabled", me.Enabled); err != nil {
		return nil, err
	}
	if err := m.Marshal("rules", me.Rules); err != nil {
		return nil, err
	}
	if err := m.Marshal("queueEntryPoint", me.QueueEntryPoint != nil && *me.QueueEntryPoint); err != nil {
		return nil, err
	}
	if err := m.Marshal("queueEntryPointType", me.QueueEntryPointType); err != nil {
		return nil, err
	}
	if err := m.Marshal("processGroups", me.ProcessGroups); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *CustomService) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "metadata")
	delete(m, "id")
	if err := m.Unmarshal("name", &me.Name); err != nil {
		return err
	}
	if err := m.Unmarshal("order", &me.Order); err != nil {
		return err
	}
	if err := m.Unmarshal("enabled", &me.Enabled); err != nil {
		return err
	}
	if err := m.Unmarshal("rules", &me.Rules); err != nil {
		return err
	}
	if err := m.Unmarshal("queueEntryPoint", &me.QueueEntryPoint); err != nil {
		return err
	}
	if err := m.Unmarshal("queueEntryPointType", &me.QueueEntryPointType); err != nil {
		return err
	}
	if err := m.Unmarshal("processGroups", &me.ProcessGroups); err != nil {
		return err
	}
	if err := m.Unmarshal("technology", &me.Technology); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

func (me *CustomService) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(me.Technology); err != nil {
		return nil, err
	}
	m["technology"] = data
	return json.MarshalIndent(m, "", "  ")
}

func (me *CustomService) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		Technology string `json:"technology"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.Technology = Technology(c.Technology)

	return nil
}
