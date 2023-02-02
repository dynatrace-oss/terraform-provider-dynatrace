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

package databases

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Scope               string               `json:"-"`
	DatabaseConnections *DatabaseConnections `json:"databaseConnections"` // Alert if the number of failed database connects within the specified time exceeds the specified absolute threshold:
	FailureRate         *FailureRate         `json:"failureRate"`         // Failure rate
	LoadDrops           *LoadDrops           `json:"loadDrops"`           // Alert if the observed load is lower than the expected load by a specified margin for a specified amount of time.
	LoadSpikes          *LoadSpikes          `json:"loadSpikes"`          // Alert if the observed load exceeds the expected load by a specified margin for a specified amount of time.
	ResponseTime        *ResponseTime        `json:"responseTime"`        // Response time
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope for the database anomaly detection",
			Required:    true,
		},
		"database_connections": {
			Type:        schema.TypeList,
			Description: "Alert if the number of failed database connects within the specified time exceeds the specified absolute threshold:",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DatabaseConnections).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"failure_rate": {
			Type:        schema.TypeList,
			Description: "Failure rate",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FailureRate).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"load_drops": {
			Type:        schema.TypeList,
			Description: "Alert if the observed load is lower than the expected load by a specified margin for a specified amount of time.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(LoadDrops).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"load_spikes": {
			Type:        schema.TypeList,
			Description: "Alert if the observed load exceeds the expected load by a specified margin for a specified amount of time.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(LoadSpikes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"response_time": {
			Type:        schema.TypeList,
			Description: "Response time",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ResponseTime).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"scope":                me.Scope,
		"database_connections": me.DatabaseConnections,
		"failure_rate":         me.FailureRate,
		"load_drops":           me.LoadDrops,
		"load_spikes":          me.LoadSpikes,
		"response_time":        me.ResponseTime,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"scope":                &me.Scope,
		"database_connections": &me.DatabaseConnections,
		"failure_rate":         &me.FailureRate,
		"load_drops":           &me.LoadDrops,
		"load_spikes":          &me.LoadSpikes,
		"response_time":        &me.ResponseTime,
	})
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) SetScope(scope string) {
	me.Scope = scope
}

func (me *Settings) GetScope() string {
	return me.Scope
}

func (me *Settings) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(me.Scope); err != nil {
		return nil, err
	}
	m["scope"] = data
	return json.MarshalIndent(m, "", "  ")
}

func (me *Settings) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		Scope string `json:"scope"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.Scope = c.Scope

	return nil
}
