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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomFilterChartSeriesDimensionConfig Configuration of the charted metric splitting
type CustomFilterChartSeriesDimensionConfig struct {
	ID              string                     `json:"id"`             // The ID of the dimension by which the metric is split
	Name            *string                    `json:"name,omitempty"` // The name of the dimension by which the metric is split
	Values          []string                   `json:"values"`         // The splitting value
	EntityDimension *bool                      `json:"entityDimension,omitempty"`
	Unknowns        map[string]json.RawMessage `json:"-"`
}

func (me *CustomFilterChartSeriesDimensionConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of the dimension by which the metric is split",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the dimension by which the metric is split",
			Optional:    true,
		},
		"values": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The splitting value",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"entity_dimension": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomFilterChartSeriesDimensionConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("id", me.ID); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("values", me.Values); err != nil {
		return err
	}
	if err := properties.Encode("entity_dimension", opt.Bool(me.EntityDimension)); err != nil {
		return err
	}
	return nil
}

func (me *CustomFilterChartSeriesDimensionConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "id")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "values")
		delete(me.Unknowns, "entity_dimension")
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
	if err := decoder.Decode("values", &me.Values); err != nil {
		return err
	}
	if value, ok := decoder.GetOk("entity_dimension"); ok {
		me.EntityDimension = opt.NewBool(value.(bool))
	}
	return nil
}

func (me *CustomFilterChartSeriesDimensionConfig) MarshalJSON() ([]byte, error) {
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
	values := me.Values
	if values == nil {
		values = []string{}
	}
	{
		rawMessage, err := json.Marshal(values)
		if err != nil {
			return nil, err
		}
		m["values"] = rawMessage
	}
	if me.EntityDimension != nil {
		rawMessage, err := json.Marshal(me.EntityDimension)
		if err != nil {
			return nil, err
		}
		m["entityDimension"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *CustomFilterChartSeriesDimensionConfig) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("id", &me.ID); err != nil {
		return err
	}
	if err := m.Unmarshal("name", &me.Name); err != nil {
		return err
	}
	if err := m.Unmarshal("values", &me.Values); err != nil {
		return err
	}
	if err := m.Unmarshal("entityDimension", &me.EntityDimension); err != nil {
		return err
	}
	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
