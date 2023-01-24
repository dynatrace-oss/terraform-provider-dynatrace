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

// AWSSupportingServiceMetric A metric of supporting service to be monitored.
type AWSSupportingServiceMetric struct {
	Name       string                     `json:"name"`       // The name of the metric of the supporting service.
	Statistic  Statistic                  `json:"statistic"`  // The statistic (aggregation) to be used for the metric. AVG_MIN_MAX value is 3 statistics at once: AVERAGE, MINIMUM and MAXIMUM
	Dimensions []string                   `json:"dimensions"` // A list of metric's dimensions names.
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (assm *AWSSupportingServiceMetric) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the name of the metric of the supporting service",
			Optional:    true,
		},
		"statistic": {
			Type:        schema.TypeString,
			Description: "the statistic (aggregation) to be used for the metric. AVG_MIN_MAX value is 3 statistics at once: AVERAGE, MINIMUM and MAXIMUM",
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
func (assm *AWSSupportingServiceMetric) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &assm.Name); err != nil {
			return err
		}
	}
	if v, found := m["statistic"]; found {
		if err := json.Unmarshal(v, &assm.Statistic); err != nil {
			return err
		}
	}
	if v, found := m["dimensions"]; found {
		if err := json.Unmarshal(v, &assm.Dimensions); err != nil {
			return err
		}
	}
	delete(m, "name")
	delete(m, "statistic")
	delete(m, "dimensions")
	if len(m) > 0 {
		assm.Unknowns = m
	}
	return nil
}

// MarshalJSON provides custom JSON serialization
func (assm *AWSSupportingServiceMetric) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(assm.Unknowns) > 0 {
		for k, v := range assm.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(assm.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(assm.Statistic)
		if err != nil {
			return nil, err
		}
		m["statistic"] = rawMessage
	}
	if assm.Dimensions != nil {
		rawMessage, err := json.Marshal(assm.Dimensions)
		if err != nil {
			return nil, err
		}
		m["dimensions"] = rawMessage
	}
	return json.Marshal(m)
}

func (assm *AWSSupportingServiceMetric) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(assm.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", assm.Name); err != nil {
		return err
	}
	if err := properties.Encode("statistic", assm.Statistic); err != nil {
		return err
	}
	if err := properties.Encode("dimensions", assm.Dimensions); err != nil {
		return err
	}
	return nil
}

func (assm *AWSSupportingServiceMetric) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), assm); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &assm.Unknowns); err != nil {
			return err
		}
		delete(assm.Unknowns, "name")
		delete(assm.Unknowns, "statistic")
		delete(assm.Unknowns, "dimensions")
		if len(assm.Unknowns) == 0 {
			assm.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		assm.Name = value.(string)
	}
	if value, ok := decoder.GetOk("statistic"); ok {
		assm.Statistic = Statistic(value.(string))
	}
	if _, ok := decoder.GetOk("dimensions.#"); ok {
		assm.Dimensions = []string{}
		if dims, ok := decoder.GetOk("dimensions"); ok {
			for _, dim := range dims.([]any) {
				assm.Dimensions = append(assm.Dimensions, dim.(string))
			}
		}
	}
	return nil
}
