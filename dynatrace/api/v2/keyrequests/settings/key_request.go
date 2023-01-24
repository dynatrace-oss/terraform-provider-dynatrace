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

package keyrequests

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// KeyRequest has no documentation
type KeyRequest struct {
	Names     []string `json:"keyRequestNames"`
	ServiceID string   `json:"-"`
}

func (me *KeyRequest) Name() string {
	return "Key Requests for " + me.ServiceID
}

func (me *KeyRequest) SetScope(scope string) {
	me.ServiceID = scope
}

func (me *KeyRequest) GetScope() string {
	return me.ServiceID
}

func (me *KeyRequest) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ID of Dynatrace Service, eg. SERVICE-123ABC45678EFGH",
		},
		"names": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The names of the key requests",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *KeyRequest) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"names":   me.Names,
		"service": me.ServiceID,
	})
}

func (me *KeyRequest) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"names":   &me.Names,
		"service": &me.ServiceID,
	})
}

func (me *KeyRequest) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(me.ServiceID); err != nil {
		return nil, err
	}
	m["serviceId"] = data
	return json.MarshalIndent(m, "", "  ")
}

func (me *KeyRequest) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		ServiceID string `json:"serviceId"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.ServiceID = c.ServiceID

	return nil
}
