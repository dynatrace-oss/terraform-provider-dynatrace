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

package infrastructurehosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NetworkErrorsDetectionThresholds struct {
	ErrorsPercentage int              `json:"errorsPercentage"` // Receive/transmit error packet percentage threshold
	EventThresholds  *EventThresholds `json:"eventThresholds"`
	TotalPacketsRate int              `json:"totalPacketsRate"` // Total packets rate threshold
}

func (me *NetworkErrorsDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"errors_percentage": {
			Type:        schema.TypeInt,
			Description: "Receive/transmit error packet percentage threshold",
			Required:    true,
		},
		"event_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(EventThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"total_packets_rate": {
			Type:        schema.TypeInt,
			Description: "Total packets rate threshold",
			Required:    true,
		},
	}
}

func (me *NetworkErrorsDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"errors_percentage":  me.ErrorsPercentage,
		"event_thresholds":   me.EventThresholds,
		"total_packets_rate": me.TotalPacketsRate,
	})
}

func (me *NetworkErrorsDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"errors_percentage":  &me.ErrorsPercentage,
		"event_thresholds":   &me.EventThresholds,
		"total_packets_rate": &me.TotalPacketsRate,
	})
}
