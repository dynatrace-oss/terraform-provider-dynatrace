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

// SlowPhysicalStorageDetectionThresholds. Alert if **any** condition is met in 4 out of 5 samples
type SlowPhysicalStorageDetectionThresholds struct {
	AvgReadWriteLatency  int `json:"avgReadWriteLatency"`  // Read/write latency is higher than
	PeakReadWriteLatency int `json:"peakReadWriteLatency"` // Peak value for read/write latency is higher than
}

func (me *SlowPhysicalStorageDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"avg_read_write_latency": {
			Type:        schema.TypeInt,
			Description: "Read/write latency is higher than",
			Required:    true,
		},
		"peak_read_write_latency": {
			Type:        schema.TypeInt,
			Description: "Peak value for read/write latency is higher than",
			Required:    true,
		},
	}
}

func (me *SlowPhysicalStorageDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"avg_read_write_latency":  me.AvgReadWriteLatency,
		"peak_read_write_latency": me.PeakReadWriteLatency,
	})
}

func (me *SlowPhysicalStorageDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"avg_read_write_latency":  &me.AvgReadWriteLatency,
		"peak_read_write_latency": &me.PeakReadWriteLatency,
	})
}
