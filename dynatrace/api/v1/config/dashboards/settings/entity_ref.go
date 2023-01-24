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

package dashboards

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// EntityRef is the short representation of a Dynatrace entity
type EntityRef struct {
	ID          string                     `json:"id"`                    // the ID of the Dynatrace entity
	Name        *string                    `json:"name,omitempty"`        // the name of the Dynatrace entity
	Description *string                    `json:"description,omitempty"` // a short description of the Dynatrace entity
	Unknowns    map[string]json.RawMessage `json:"-"`
}

func (me *EntityRef) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "the ID of the Dynatrace entity",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "the name of the Dynatrace entity",
			Optional:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "a short description of the Dynatrace entity",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *EntityRef) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("id", me.ID); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("description", me.Description); err != nil {
		return err
	}
	return nil
}

func (me *EntityRef) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "id")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "description")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}

	if value, ok := decoder.GetOk("id"); ok {
		me.ID = value.(string)
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("description"); ok {
		me.Description = opt.NewString(value.(string))
	}
	return nil
}

func (me *EntityRef) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.ID)
		if err != nil {
			return nil, err
		}
		m["id"] = rawMessage
	}
	if me.Name != nil {
		rawMessage, err := json.Marshal(me.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	if me.Description != nil {
		rawMessage, err := json.Marshal(me.Description)
		if err != nil {
			return nil, err
		}
		m["description"] = rawMessage
	}
	return json.Marshal(m)
}
