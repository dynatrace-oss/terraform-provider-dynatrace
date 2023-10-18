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

package namespace

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CpuLimitsQuotaSaturation      *CpuLimitsQuotaSaturation      `json:"cpuLimitsQuotaSaturation"`      // Alerts if almost no CPU-limit quota left in namespace
	CpuRequestsQuotaSaturation    *CpuRequestsQuotaSaturation    `json:"cpuRequestsQuotaSaturation"`    // Alerts if almost no CPU-request quota left in namespace
	MemoryLimitsQuotaSaturation   *MemoryLimitsQuotaSaturation   `json:"memoryLimitsQuotaSaturation"`   // Alerts if almost no memory-limit quota left in namespace
	MemoryRequestsQuotaSaturation *MemoryRequestsQuotaSaturation `json:"memoryRequestsQuotaSaturation"` // Alerts if almost no memory-request quota left in namespace
	PodsQuotaSaturation           *PodsQuotaSaturation           `json:"podsQuotaSaturation"`           // Alerts if almost no pod quota left in namespace
	Scope                         *string                        `json:"-" scope:"scope"`               // The scope of this setting (CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpu_limits_quota_saturation": {
			Type:        schema.TypeList,
			Description: "Alerts if almost no CPU-limit quota left in namespace",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(CpuLimitsQuotaSaturation).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"cpu_requests_quota_saturation": {
			Type:        schema.TypeList,
			Description: "Alerts if almost no CPU-request quota left in namespace",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(CpuRequestsQuotaSaturation).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"memory_limits_quota_saturation": {
			Type:        schema.TypeList,
			Description: "Alerts if almost no memory-limit quota left in namespace",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(MemoryLimitsQuotaSaturation).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"memory_requests_quota_saturation": {
			Type:        schema.TypeList,
			Description: "Alerts if almost no memory-request quota left in namespace",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(MemoryRequestsQuotaSaturation).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"pods_quota_saturation": {
			Type:        schema.TypeList,
			Description: "Alerts if almost no pod quota left in namespace",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(PodsQuotaSaturation).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cpu_limits_quota_saturation":      me.CpuLimitsQuotaSaturation,
		"cpu_requests_quota_saturation":    me.CpuRequestsQuotaSaturation,
		"memory_limits_quota_saturation":   me.MemoryLimitsQuotaSaturation,
		"memory_requests_quota_saturation": me.MemoryRequestsQuotaSaturation,
		"pods_quota_saturation":            me.PodsQuotaSaturation,
		"scope":                            me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cpu_limits_quota_saturation":      &me.CpuLimitsQuotaSaturation,
		"cpu_requests_quota_saturation":    &me.CpuRequestsQuotaSaturation,
		"memory_limits_quota_saturation":   &me.MemoryLimitsQuotaSaturation,
		"memory_requests_quota_saturation": &me.MemoryRequestsQuotaSaturation,
		"pods_quota_saturation":            &me.PodsQuotaSaturation,
		"scope":                            &me.Scope,
	})
}
