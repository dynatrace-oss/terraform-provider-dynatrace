/**
* @license
* Copyright 2025 Dynatrace LLC
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

package ddulimit

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (me *DDUPool) Deprecated() string {
	return "This resource API endpoint has been deprecated."
}

// DDUPool TODO: documentation
type DDUPool struct {
	MetricsPool       DDUPoolConfig `json:"metrics"`
	LogMonitoringPool DDUPoolConfig `json:"logMonitoring"`
	ServerlessPool    DDUPoolConfig `json:"serverless"`
	EventsPool        DDUPoolConfig `json:"events"`
	TracesPool        DDUPoolConfig `json:"traces"`
}

func (me *DDUPool) Name() string {
	return "ddupool"
}

func (me *DDUPool) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "DDU pool settings for Metrics",
			Elem: &schema.Resource{
				Schema: new(DDUPoolConfig).Schema(),
			},
		},
		"log_monitoring": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "DDU pool settings for Log Monitoring",
			Elem: &schema.Resource{
				Schema: new(DDUPoolConfig).Schema(),
			},
		},
		"serverless": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "DDU pool settings for Serverless",
			Elem: &schema.Resource{
				Schema: new(DDUPoolConfig).Schema(),
			},
		},
		"events": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "DDU pool settings for Events",
			Elem: &schema.Resource{
				Schema: new(DDUPoolConfig).Schema(),
			},
		},
		"traces": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "DDU pool settings for Traces",
			Elem: &schema.Resource{
				Schema: new(DDUPoolConfig).Schema(),
			},
		},
	}
}

func (me *DDUPool) MarshalHCL(properties hcl.Properties) error {

	if me.MetricsPool.LimitEnabled {
		if err := properties.Encode("metrics", &me.MetricsPool); err != nil {
			return err
		}
	}
	if me.LogMonitoringPool.LimitEnabled {
		if err := properties.Encode("log_monitoring", &me.LogMonitoringPool); err != nil {
			return err
		}
	}
	if me.ServerlessPool.LimitEnabled {
		if err := properties.Encode("serverless", &me.ServerlessPool); err != nil {
			return err
		}
	}
	if me.EventsPool.LimitEnabled {
		if err := properties.Encode("events", &me.EventsPool); err != nil {
			return err
		}
	}
	if me.TracesPool.LimitEnabled {
		if err := properties.Encode("traces", &me.TracesPool); err != nil {
			return err
		}
	}

	return nil
}

func (me *DDUPool) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]interface{}{
		"metrics":        &me.MetricsPool,
		"log_monitoring": &me.LogMonitoringPool,
		"serverless":     &me.ServerlessPool,
		"events":         &me.EventsPool,
		"traces":         &me.TracesPool,
	})
}
