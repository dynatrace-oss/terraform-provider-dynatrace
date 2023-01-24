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

package aws

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AWSConfigTag An AWS tag of the resource to be monitored.
type AWSConfigTag struct {
	Name     string                     `json:"name"`  // The key of the AWS tag.
	Value    string                     `json:"value"` // The value of the AWS tag.
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (act *AWSConfigTag) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the key of the AWS tag.",
			Optional:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "the value of the AWS tag",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

// UnmarshalJSON provides custom JSON deserialization
func (act *AWSConfigTag) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &act.Name); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &act.Value); err != nil {
			return err
		}
	}
	delete(m, "name")
	delete(m, "value")
	if len(m) > 0 {
		act.Unknowns = m
	}
	return nil
}

// MarshalJSON provides custom JSON serialization
func (act *AWSConfigTag) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(act.Unknowns) > 0 {
		for k, v := range act.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(act.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(act.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	return json.Marshal(m)
}

func (act *AWSConfigTag) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(act.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", act.Name); err != nil {
		return err
	}
	if err := properties.Encode("value", act.Value); err != nil {
		return err
	}
	return nil
}

func (act *AWSConfigTag) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), act); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &act.Unknowns); err != nil {
			return err
		}
		delete(act.Unknowns, "name")
		delete(act.Unknowns, "value")
		if len(act.Unknowns) == 0 {
			act.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		act.Name = value.(string)
	}
	if value, ok := decoder.GetOk("value"); ok {
		act.Value = value.(string)
	}
	return nil
}
