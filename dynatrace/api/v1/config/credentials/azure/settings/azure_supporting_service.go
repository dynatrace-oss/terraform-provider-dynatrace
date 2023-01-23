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

// AzureSupportingService A supporting service to be monitored.
type AzureSupportingService struct {
	Name             *string                    `json:"name,omitempty"`             // The name of the supporting service.
	MonitoredMetrics []*AzureMonitoredMetric    `json:"monitoredMetrics,omitempty"` // A list of metrics to be monitored for this service. It must include all the recommended metrics.
	Unknowns         map[string]json.RawMessage `json:"-"`
}

func (ass *AzureSupportingService) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the supporting service.",
			Optional:    true,
		},
		"monitored_metrics": {
			Type:        schema.TypeList,
			Description: "A list of Azure tags to be monitored.  You can specify up to 10 tags. A resource tagged with *any* of the specified tags is monitored.  Only applicable when the **monitorOnlyTaggedEntities** parameter is set to `true`",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(AzureMonitoredMetric).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (ass *AzureSupportingService) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(ass.Unknowns) > 0 {
		for k, v := range ass.Unknowns {
			m[k] = v
		}
	}
	if ass.Name != nil {
		rawMessage, err := json.Marshal(ass.Name)
		if err != nil {
			return nil, err
		}
		m["name"] = rawMessage
	}
	if ass.MonitoredMetrics != nil {
		rawMessage, err := json.Marshal(ass.MonitoredMetrics)
		if err != nil {
			return nil, err
		}
		m["monitoredMetrics"] = rawMessage
	}
	return json.Marshal(m)
}

func (ass *AzureSupportingService) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["name"]; found {
		if err := json.Unmarshal(v, &ass.Name); err != nil {
			return err
		}
	}
	if v, found := m["monitoredMetrics"]; found {
		if err := json.Unmarshal(v, &ass.MonitoredMetrics); err != nil {
			return err
		}
	}
	delete(m, "name")
	delete(m, "monitoredMetrics")

	if len(m) > 0 {
		ass.Unknowns = m
	}
	return nil
}

func (ass *AzureSupportingService) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(ass.Unknowns); err != nil {
		return err
	}

	if err := properties.Encode("name", ass.Name); err != nil {
		return err
	}
	if err := properties.Encode("monitored_metrics", ass.MonitoredMetrics); err != nil {
		return err
	}
	return nil
}

func (ass *AzureSupportingService) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), ass); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &ass.Unknowns); err != nil {
			return err
		}
		delete(ass.Unknowns, "name")
		delete(ass.Unknowns, "monitored_metrics")

		if len(ass.Unknowns) == 0 {
			ass.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		ass.Name = opt.NewString(value.(string))
	}
	if result, ok := decoder.GetOk("monitored_metrics.#"); ok {
		ass.MonitoredMetrics = []*AzureMonitoredMetric{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(AzureMonitoredMetric)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "monitored_metrics", idx)); err != nil {
				return err
			}
			ass.MonitoredMetrics = append(ass.MonitoredMetrics, entry)
		}
	}
	return nil
}
