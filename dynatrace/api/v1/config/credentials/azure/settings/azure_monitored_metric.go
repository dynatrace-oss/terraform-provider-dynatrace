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

// AzureMonitoredMetric A metric of supporting service to be monitored.
type AzureMonitoredMetric struct {
	Name       *string                    `json:"name,omitempty"`       // The name of the metric of the supporting service.
	Dimensions []string                   `json:"dimensions,omitempty"` // A list of metric's dimensions names. It must include all the recommended dimensions.
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (amm *AzureMonitoredMetric) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the name of the metric of the supporting service",
			Optional:    true,
		},
		"dimensions": {
			Type:        schema.TypeList,
			Description: "a list of metric's dimensions names",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

// UnmarshalJSON provides custom JSON deserialization
func (amm *AzureMonitoredMetric) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &amm.Name); err != nil {
			return err
		}
	}
	if v, found := m["dimensions"]; found {
		if err := json.Unmarshal(v, &amm.Dimensions); err != nil {
			return err
		}
	}
	delete(m, "name")
	delete(m, "dimensions")
	if len(m) > 0 {
		amm.Unknowns = m
	}
	return nil
}

// MarshalJSON provides custom JSON serialization
func (amm *AzureMonitoredMetric) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(amm.Unknowns) > 0 {
		for k, v := range amm.Unknowns {
			m[k] = v
		}
	}
	if amm.Name != nil {
		rawMessage, err := json.Marshal(amm.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	if amm.Dimensions != nil {
		rawMessage, err := json.Marshal(amm.Dimensions)
		if err != nil {
			return nil, err
		}
		m["dimensions"] = rawMessage
	}
	return json.Marshal(m)
}

func (amm *AzureMonitoredMetric) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(amm.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", amm.Name); err != nil {
		return err
	}
	if err := properties.Encode("dimensions", amm.Dimensions); err != nil {
		return err
	}
	return nil
}

func (amm *AzureMonitoredMetric) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), amm); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &amm.Unknowns); err != nil {
			return err
		}
		delete(amm.Unknowns, "name")
		delete(amm.Unknowns, "dimensions")
		if len(amm.Unknowns) == 0 {
			amm.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		amm.Name = opt.NewString(value.(string))
	}
	if _, ok := decoder.GetOk("dimensions.#"); ok {
		amm.Dimensions = []string{}
		if dims, ok := decoder.GetOk("dimensions"); ok {
			for _, dim := range dims.([]any) {
				amm.Dimensions = append(amm.Dimensions, dim.(string))
			}
		}
	}
	return nil
}
