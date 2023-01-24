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

package azure

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CloudTag A cloud tag.
type CloudTag struct {
	Value    *string                    `json:"value,omitempty"` // The value of the tag. If set to `null`, then resources with any value of the tag are monitored.
	Name     *string                    `json:"name,omitempty"`  // The name of the tag.
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (ct *CloudTag) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the tag.   If set to `null`, then resources with any value of the tag are monitored.",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the tag.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (ct *CloudTag) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(ct.Unknowns) > 0 {
		for k, v := range ct.Unknowns {
			m[k] = v
		}
	}
	if ct.Value != nil {
		rawMessage, err := json.Marshal(ct.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	if ct.Name != nil {
		rawMessage, err := json.Marshal(ct.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	return json.Marshal(m)
}

func (ct *CloudTag) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &ct.Value); err != nil {
			return err
		}
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &ct.Name); err != nil {
			return err
		}
	}
	delete(m, "value")
	delete(m, "name")

	if len(m) > 0 {
		ct.Unknowns = m
	}
	return nil
}

func (ct *CloudTag) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(ct.Unknowns); err != nil {
		return err
	}

	if err := properties.Encode("value", ct.Value); err != nil {
		return err
	}
	if err := properties.Encode("name", ct.Name); err != nil {
		return err
	}
	return nil
}

func (ct *CloudTag) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), ct); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &ct.Unknowns); err != nil {
			return err
		}
		delete(ct.Unknowns, "value")
		delete(ct.Unknowns, "name")
		if len(ct.Unknowns) == 0 {
			ct.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("value"); ok {
		ct.Value = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("name"); ok {
		ct.Name = opt.NewString(value.(string))
	}
	return nil
}
