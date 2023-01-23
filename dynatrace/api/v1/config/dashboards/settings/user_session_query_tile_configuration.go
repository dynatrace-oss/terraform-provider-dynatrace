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

// UserSessionQueryTileConfiguration Configuration of a User session query visualization tile
type UserSessionQueryTileConfiguration struct {
	HasAxisBucketing *bool                      `json:"hasAxisBucketing,omitempty"` // The axis bucketing when enabled groups similar series in the same virtual axis
	Unknowns         map[string]json.RawMessage `json:"-"`
}

func (me *UserSessionQueryTileConfiguration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"has_axis_bucketing": {
			Type:        schema.TypeBool,
			Description: "The axis bucketing when enabled groups similar series in the same virtual axis",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *UserSessionQueryTileConfiguration) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("has_axis_bucketing", opt.Bool(me.HasAxisBucketing)); err != nil {
		return err
	}
	return nil
}

func (me *UserSessionQueryTileConfiguration) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "has_axis_bucketing")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("has_axis_bucketing"); ok {
		me.HasAxisBucketing = opt.NewBool(value.(bool))
	}
	return nil
}

func (me *UserSessionQueryTileConfiguration) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.HasAxisBucketing != nil {
		rawMessage, err := json.Marshal(me.HasAxisBucketing)
		if err != nil {
			return nil, err
		}
		m["hasAxisBucketing"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *UserSessionQueryTileConfiguration) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["hasAxisBucketing"]; found {
		if err := json.Unmarshal(v, &me.HasAxisBucketing); err != nil {
			return err
		}
	}
	delete(m, "hasAxisBucketing")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
