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

package disks

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Scope string `json:"-"`
	Disk  *Disk  `json:"disk"` // Disk
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope for the disk anomaly detection",
			Required:    true,
			ForceNew:    true,
		},
		"disk": {
			Type:        schema.TypeList,
			Description: "Disk",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Disk).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"scope": me.Scope,
		"disk":  me.Disk,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"scope": &me.Scope,
		"disk":  &me.Disk,
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
