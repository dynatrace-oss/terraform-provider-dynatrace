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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RdsHighWriteReadLatencyDetectionThresholds. Alert if the condition is met in 3 out of 5 samples
type RdsHighWriteReadLatencyDetectionThresholds struct {
	ReadWriteLatency int `json:"readWriteLatency"` // Read/write latency is higher than
}

func (me *RdsHighWriteReadLatencyDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"read_write_latency": {
			Type:        schema.TypeInt,
			Description: "Read/write latency is higher than",
			Required:    true,
		},
	}
}

func (me *RdsHighWriteReadLatencyDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"read_write_latency": me.ReadWriteLatency,
	})
}

func (me *RdsHighWriteReadLatencyDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"read_write_latency": &me.ReadWriteLatency,
	})
}
