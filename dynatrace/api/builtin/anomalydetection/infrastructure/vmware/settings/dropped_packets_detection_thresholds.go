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

package vmware

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DroppedPacketsDetectionThresholds. Alert if the condition is met in 3 out of 5 samples
type DroppedPacketsDetectionThresholds struct {
	DroppedPacketsPerSecond int `json:"droppedPacketsPerSecond"` // Receive/transmit dropped packets rate on NIC is higher than
}

func (me *DroppedPacketsDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dropped_packets_per_second": {
			Type:        schema.TypeInt,
			Description: "Receive/transmit dropped packets rate on NIC is higher than",
			Required:    true,
		},
	}
}

func (me *DroppedPacketsDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dropped_packets_per_second": me.DroppedPacketsPerSecond,
	})
}

func (me *DroppedPacketsDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dropped_packets_per_second": &me.DroppedPacketsPerSecond,
	})
}
