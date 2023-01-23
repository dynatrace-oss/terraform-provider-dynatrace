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

// AWSSupportingServiceConfig A supporting service to be monitored.
type AWSSupportingServiceConfig struct {
	Name             string                        `json:"name"`             // The name of the supporting service.
	MonitoredMetrics []*AWSSupportingServiceMetric `json:"monitoredMetrics"` // A list of metrics to be monitored for this service.
	Unknowns         map[string]json.RawMessage    `json:"-"`
}

func (assc *AWSSupportingServiceConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the name of the supporting service",
			Optional:    true,
		},
		"monitored_metrics": {
			Type:        schema.TypeList,
			Description: "a list of metrics to be monitored for this service",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(AWSSupportingServiceMetric).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

// UnmarshalJSON provides custom JSON deserialization
func (assc *AWSSupportingServiceConfig) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &assc.Name); err != nil {
			return err
		}
	}
	if v, found := m["monitoredMetrics"]; found {
		if err := json.Unmarshal(v, &assc.MonitoredMetrics); err != nil {
			return err
		}
	}
	delete(m, "name")
	delete(m, "monitoredMetrics")
	if len(m) > 0 {
		assc.Unknowns = m
	}
	return nil
}

// MarshalJSON provides custom JSON serialization
func (assc *AWSSupportingServiceConfig) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(assc.Unknowns) > 0 {
		for k, v := range assc.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(assc.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	if assc.MonitoredMetrics != nil {
		rawMessage, err := json.Marshal(assc.MonitoredMetrics)
		if err != nil {
			return nil, err
		}
		m["monitoredMetrics"] = rawMessage
	}
	return json.Marshal(m)
}

func (assc *AWSSupportingServiceConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(assc.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", assc.Name); err != nil {
		return err
	}
	if err := properties.Encode("monitored_metrics", assc.MonitoredMetrics); err != nil {
		return err
	}
	return nil
}

func (assc *AWSSupportingServiceConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), assc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &assc.Unknowns); err != nil {
			return err
		}
		delete(assc.Unknowns, "name")
		delete(assc.Unknowns, "monitored_metrics")
		if len(assc.Unknowns) == 0 {
			assc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		assc.Name = value.(string)
	}
	if result, ok := decoder.GetOk("monitored_metrics.#"); ok {
		assc.MonitoredMetrics = []*AWSSupportingServiceMetric{}
		for idx := 0; idx < result.(int); idx++ {
			monitoredMetric := new(AWSSupportingServiceMetric)
			if err := monitoredMetric.UnmarshalHCL(hcl.NewDecoder(decoder, "monitored_metrics", idx)); err != nil {
				return err
			}
			assc.MonitoredMetrics = append(assc.MonitoredMetrics, monitoredMetric)
		}
	}
	return nil
}
