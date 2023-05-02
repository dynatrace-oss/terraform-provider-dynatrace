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

// GuestCPULimitDetectionThresholds. Alert if **all three** conditions are met in 3 out of 5 samples
type GuestCPULimitDetectionThresholds struct {
	HostCpuUsagePercentage int `json:"hostCpuUsagePercentage"` // Hypervisor CPU usage is higher than
	VmCpuReadyPercentage   int `json:"vmCpuReadyPercentage"`   // VM CPU ready is higher than
	VmCpuUsagePercentage   int `json:"vmCpuUsagePercentage"`   // VM CPU usage (VM CPU Usage Mhz / VM CPU limit in Mhz) is higher than
}

func (me *GuestCPULimitDetectionThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_cpu_usage_percentage": {
			Type:        schema.TypeInt,
			Description: "Hypervisor CPU usage is higher than",
			Required:    true,
		},
		"vm_cpu_ready_percentage": {
			Type:        schema.TypeInt,
			Description: "VM CPU ready is higher than",
			Required:    true,
		},
		"vm_cpu_usage_percentage": {
			Type:        schema.TypeInt,
			Description: "VM CPU usage (VM CPU Usage Mhz / VM CPU limit in Mhz) is higher than",
			Required:    true,
		},
	}
}

func (me *GuestCPULimitDetectionThresholds) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"host_cpu_usage_percentage": me.HostCpuUsagePercentage,
		"vm_cpu_ready_percentage":   me.VmCpuReadyPercentage,
		"vm_cpu_usage_percentage":   me.VmCpuUsagePercentage,
	})
}

func (me *GuestCPULimitDetectionThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"host_cpu_usage_percentage": &me.HostCpuUsagePercentage,
		"vm_cpu_ready_percentage":   &me.VmCpuReadyPercentage,
		"vm_cpu_usage_percentage":   &me.VmCpuUsagePercentage,
	})
}
