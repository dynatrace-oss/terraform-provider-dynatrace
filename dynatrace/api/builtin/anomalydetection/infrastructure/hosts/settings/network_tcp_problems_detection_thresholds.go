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

package hosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NetworkTcpProblemsDetectionThresholds struct {
	EventThresholds                  *EventThresholds `json:"eventThresholds"`
	FailedConnectionsNumberPerMinute int              `json:"failedConnectionsNumberPerMinute"` // Number of failed connections threshold
	NewConnectionFailuresPercentage  int              `json:"newConnectionFailuresPercentage"`  // New connection failure threshold
}

func (me *NetworkTcpProblemsDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"failed_connections_number_per_minute": {
			Type:        schema.TypeInt,
			Description: "Number of failed connections threshold",
			Required:    true,
		},
		"new_connection_failures_percentage": {
			Type:        schema.TypeInt,
			Description: "New connection failure threshold",
			Required:    true,
		},
	}
}

func (me *NetworkTcpProblemsDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"event_thresholds":                     me.EventThresholds,
		"failed_connections_number_per_minute": me.FailedConnectionsNumberPerMinute,
		"new_connection_failures_percentage":   me.NewConnectionFailuresPercentage,
	})
}

func (me *NetworkTcpProblemsDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"event_thresholds":                     &me.EventThresholds,
		"failed_connections_number_per_minute": &me.FailedConnectionsNumberPerMinute,
		"new_connection_failures_percentage":   &me.NewConnectionFailuresPercentage,
	})
}
