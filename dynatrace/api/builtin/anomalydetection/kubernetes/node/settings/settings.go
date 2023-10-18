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

package node

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CpuRequestsSaturation    *CpuRequestsSaturation    `json:"cpuRequestsSaturation"`
	MemoryRequestsSaturation *MemoryRequestsSaturation `json:"memoryRequestsSaturation"`
	NodeProblematicCondition *NodeProblematicCondition `json:"nodeProblematicCondition"`
	PodsSaturation           *PodsSaturation           `json:"podsSaturation"`
	ReadinessIssues          *ReadinessIssues          `json:"readinessIssues"` // Alerts if node has not been available for a given amount of time
	Scope                    *string                   `json:"-" scope:"scope"` // The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cpu_requests_saturation": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(CpuRequestsSaturation).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"memory_requests_saturation": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(MemoryRequestsSaturation).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"node_problematic_condition": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(NodeProblematicCondition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"pods_saturation": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(PodsSaturation).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"readiness_issues": {
			Type:        schema.TypeList,
			Description: "Alerts if node has not been available for a given amount of time",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ReadinessIssues).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cpu_requests_saturation":    me.CpuRequestsSaturation,
		"memory_requests_saturation": me.MemoryRequestsSaturation,
		"node_problematic_condition": me.NodeProblematicCondition,
		"pods_saturation":            me.PodsSaturation,
		"readiness_issues":           me.ReadinessIssues,
		"scope":                      me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cpu_requests_saturation":    &me.CpuRequestsSaturation,
		"memory_requests_saturation": &me.MemoryRequestsSaturation,
		"node_problematic_condition": &me.NodeProblematicCondition,
		"pods_saturation":            &me.PodsSaturation,
		"readiness_issues":           &me.ReadinessIssues,
		"scope":                      &me.Scope,
	})
}
