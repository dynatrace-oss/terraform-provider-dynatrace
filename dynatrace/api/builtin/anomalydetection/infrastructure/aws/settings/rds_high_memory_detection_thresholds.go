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

// RdsHighMemoryDetectionThresholds. Alert if **both** conditions is met in 3 out of 5 samples
type RdsHighMemoryDetectionThresholds struct {
	FreeMemory float64 `json:"freeMemory"` // Freeable memory is lower than
	SwapUsage  float64 `json:"swapUsage"`  // Swap usage is higher than
}

func (me *RdsHighMemoryDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"free_memory": {
			Type:        schema.TypeFloat,
			Description: "Freeable memory is lower than",
			Required:    true,
		},
		"swap_usage": {
			Type:        schema.TypeFloat,
			Description: "Swap usage is higher than",
			Required:    true,
		},
	}
}

func (me *RdsHighMemoryDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"free_memory": me.FreeMemory,
		"swap_usage":  me.SwapUsage,
	})
}

func (me *RdsHighMemoryDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"free_memory": &me.FreeMemory,
		"swap_usage":  &me.SwapUsage,
	})
}
