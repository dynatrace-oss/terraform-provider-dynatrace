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

// Scope The scope of the maintenance window.
//
//	The scope restricts the alert/problem detection suppression to certain Dynatrace entities. It can contain a list of entities and/or matching rules for dynamic formation of the scope.
//	If no scope is specified, the alert/problem detection suppression applies to the entire environment.
type Scope struct {
	Entities []string                   `json:"entities"` // A list of Dynatrace entities (for example, hosts or services) to be included in the scope.  Allowed values are Dynatrace entity IDs.
	Matches  []*Filter                  `json:"matches"`  // A list of matching rules for dynamic scope formation.  If several rules are set, the OR logic applies.
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *Scope) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entities": {
			Type:        schema.TypeSet,
			Description: "A list of Dynatrace entities (for example, hosts or services) to be included in the scope.  Allowed values are Dynatrace entity IDs",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"matches": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A list of matching rules for dynamic scope formation.  If several rules are set, the OR logic applies",
			Elem: &schema.Resource{
				Schema: new(Filter).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Scope) IsEmpty() bool {
	return len(me.Entities) == 0 && len(me.Matches) == 0
}

func (me *Scope) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("entities", me.Entities); err != nil {
		return err
	}
	if err := properties.Encode("matches", me.Matches); err != nil {
		return err
	}
	return nil
}

func (me *Scope) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "entities")
		delete(me.Unknowns, "matches")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if err := decoder.Decode("entities", &me.Entities); err != nil {
		return err
	}
	if result, ok := decoder.GetOk("matches.#"); ok {
		me.Matches = []*Filter{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Filter)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "matches", idx)); err != nil {
				return err
			}
			me.Matches = append(me.Matches, entry)
		}
	}
	return nil
}

func (me *Scope) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if len(me.Entities) == 0 {
		data, err := json.Marshal([]string{})
		if err != nil {
			return nil, err
		}
		m["entities"] = data
	} else {
		if err := m.Marshal("entities", me.Entities); err != nil {
			return nil, err
		}
	}
	if len(me.Matches) == 0 {
		data, err := json.Marshal([]string{})
		if err != nil {
			return nil, err
		}
		m["matches"] = data
	} else {
		if err := m.Marshal("matches", me.Matches); err != nil {
			return nil, err
		}
	}
	return json.Marshal(m)
}

func (me *Scope) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("entities", &me.Entities); err != nil {
		return err
	}
	if err := m.Unmarshal("matches", &me.Matches); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
