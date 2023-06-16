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

package services

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AWSMonitoredMetric A metric of supporting service to be monitored.
type AWSMonitoredMetric struct {
	Name       string   `json:"name,omitempty"` // The name of the metric of the supporting service.
	Statistic  string   `json:"statistic"`
	Dimensions []string `json:"dimensions"` // A list of metric's dimensions names. It must include all the recommended dimensions.
}

func (amm *AWSMonitoredMetric) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the name of the metric of the supporting service",
			Required:    true,
		},
		"statistic": {
			Type:        schema.TypeString,
			Description: "Possible values are `AVERAGE`, `AVG_MIN_MAX`, `MAXIMUM`, `MINIMUM`, `SAMPLE_COUNT` and `SUM`",
			Optional:    true,
			Default:     "AVERAGE",
		},
		"dimensions": {
			Type:        schema.TypeSet,
			Description: "a list of metric's dimensions names",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (amm *AWSMonitoredMetric) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", amm.Name); err != nil {
		return err
	}
	if err := properties.Encode("dimensions", amm.Dimensions); err != nil {
		return err
	}
	if err := properties.Encode("statistic", amm.Statistic); err != nil {
		return err
	}
	return nil
}

func (amm *AWSMonitoredMetric) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		amm.Name = value.(string)
	}
	if err := decoder.Decode("dimensions", &amm.Dimensions); err != nil {
		return err
	}
	if err := decoder.Decode("statistic", &amm.Statistic); err != nil {
		return err
	}
	return nil
}
