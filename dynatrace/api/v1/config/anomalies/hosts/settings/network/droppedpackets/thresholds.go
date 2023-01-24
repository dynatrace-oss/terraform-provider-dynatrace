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

package droppedpackets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Thresholds Custom thresholds for dropped packets. If not set, automatic mode is used.
//
//	**All** of these conditions must be met to trigger an alert.
type Thresholds struct {
	DroppedPacketsPercentage int32 `json:"droppedPacketsPercentage"` // Receive/transmit dropped packet percentage is higher than *X*% in 3 out of 5 samples.
	TotalPacketsRate         int32 `json:"totalPacketsRate"`         // Total receive/transmit packets rate is higher than *X* packets per second in 3 out of 5 samples.
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dropped_packets": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Receive/transmit dropped packet percentage is higher than *X*% in 3 out of 5 samples",
		},
		"total_packets_rate": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Total receive/transmit packets rate is higher than *X* packets per second in 3 out of 5 samples",
		},
	}
}

func (me *Thresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dropped_packets":    me.DroppedPacketsPercentage,
		"total_packets_rate": me.TotalPacketsRate,
	})
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("dropped_packets"); ok {
		me.DroppedPacketsPercentage = int32(value.(int))
	}
	if value, ok := decoder.GetOk("total_packets_rate"); ok {
		me.TotalPacketsRate = int32(value.(int))
	}
	return nil
}
