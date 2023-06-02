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

type Settings struct {
	CredentialsID         string                  `json:"-"`
	Name                  string                  `json:"name"`                       // The name of the supporting service.
	MonitoredMetrics      []*AzureMonitoredMetric `json:"monitoredMetrics,omitempty"` // A list of metrics to be monitored for this service. It must include all the recommended metrics.
	BuiltIn               bool                    `json:"-"`
	RequiredMetrics       string                  `json:"-"`
	UseRecommendedMetrics bool                    `json:"-"`
}

func (me *Settings) IsComputer() bool {
	return true
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"credentials_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "the ID of the azure credentials this supported service belongs to",
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the supporting service.",
			Optional:    true,
			ForceNew:    true,
		},
		"use_recommended_metrics": {
			Type:        schema.TypeBool,
			Description: "If `true` Terraform will negotiate with the Dynatrace API about the recommended/enforced metrics to be applied. Any `metric` specified will be therefore ignored.",
			Optional:    true,
		},
		"metric": {
			Type:        schema.TypeSet,
			Description: "A list of metrics to be monitored for this service. Depending on the service Dynatrace insists on a set of recommended metrics to be configured for that service. If any of these recommended metrics is missing here, the Terraform Provider will automatically add them during `terraform apply`. This usually results in a non-empty plan, until all of the recommended metrics are present within your configuration. For services considered `built-in` by Dynatrace any metrics specified here will be ignored - Dynatrace enforces a fixed set of metrics for these services.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(AzureMonitoredMetric).Schema()},
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				built_in := d.Get("built_in")
				use_recommended_metrics := d.Get("use_recommended_metrics")
				return (built_in != nil && built_in.(bool)) || (use_recommended_metrics != nil && use_recommended_metrics.(bool))
			},
		},
		"required_metrics": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Used internally by the Terraform Provider in order to remember the metrics enforced by Dynatrace",
		},
		"built_in": {
			Type:        schema.TypeBool,
			Description: "This attribute is automatically set to `true` if Dynatrace considers the supporting service with the given name to be a built-in service",
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("credentials_id", me.CredentialsID); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("built_in", me.BuiltIn); err != nil {
		return err
	}
	if me.UseRecommendedMetrics {
		properties["use_recommended_metrics"] = true
	}
	if len(me.RequiredMetrics) > 0 {
		if err := properties.Encode("required_metrics", me.RequiredMetrics); err != nil {
			return err
		}
	}
	if err := properties.EncodeSlice("metric", me.MonitoredMetrics); err != nil {
		return err
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("credentials_id", &me.CredentialsID); err != nil {
		return err
	}
	if err := decoder.Decode("name", &me.Name); err != nil {
		return err
	}
	if err := decoder.Decode("built_in", &me.BuiltIn); err != nil {
		return err
	}
	if err := decoder.Decode("use_recommended_metrics", &me.UseRecommendedMetrics); err != nil {
		return err
	}
	if err := decoder.Decode("required_metrics", &me.RequiredMetrics); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("metric", &me.MonitoredMetrics); err != nil {
		return err
	}
	return nil
}
