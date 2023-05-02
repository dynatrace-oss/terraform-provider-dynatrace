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

// EsxiHighCpuDetectionThresholds. Alert if **all three** conditions are met in 3 out of 5 samples
type EsxiHighCpuDetectionThresholds struct {
	CpuPeakPercentage    int `json:"cpuPeakPercentage"`    // At least one peak occurred when Hypervisor CPU usage was higher than
	CpuUsagePercentage   int `json:"cpuUsagePercentage"`   // CPU usage is higher than
	VmCpuReadyPercentage int `json:"vmCpuReadyPercentage"` // VM CPU ready is higher than
}

func (me *EsxiHighCpuDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpu_peak_percentage": {
			Type:        schema.TypeInt,
			Description: "At least one peak occurred when Hypervisor CPU usage was higher than",
			Required:    true,
		},
		"cpu_usage_percentage": {
			Type:        schema.TypeInt,
			Description: "CPU usage is higher than",
			Required:    true,
		},
		"vm_cpu_ready_percentage": {
			Type:        schema.TypeInt,
			Description: "VM CPU ready is higher than",
			Required:    true,
		},
	}
}

func (me *EsxiHighCpuDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cpu_peak_percentage":     me.CpuPeakPercentage,
		"cpu_usage_percentage":    me.CpuUsagePercentage,
		"vm_cpu_ready_percentage": me.VmCpuReadyPercentage,
	})
}

func (me *EsxiHighCpuDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cpu_peak_percentage":     &me.CpuPeakPercentage,
		"cpu_usage_percentage":    &me.CpuUsagePercentage,
		"vm_cpu_ready_percentage": &me.VmCpuReadyPercentage,
	})
}
