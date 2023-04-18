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

package workload

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ContainerRestarts        *ContainerRestarts        `json:"containerRestarts"`
	DeploymentStuck          *DeploymentStuck          `json:"deploymentStuck"`
	HighCpuThrottling        *HighCpuThrottling        `json:"highCpuThrottling"`
	HighCpuUsage             *HighCpuUsage             `json:"highCpuUsage"`
	HighMemoryUsage          *HighMemoryUsage          `json:"highMemoryUsage"`
	NotAllPodsReady          *NotAllPodsReady          `json:"notAllPodsReady"`
	PendingPods              *PendingPods              `json:"pendingPods"`
	PodStuckInTerminating    *PodStuckInTerminating    `json:"podStuckInTerminating"`
	Scope                    *string                   `json:"-" scope:"scope"` // The scope of this setting (CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.
	WorkloadWithoutReadyPods *WorkloadWithoutReadyPods `json:"workloadWithoutReadyPods"`
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"container_restarts": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ContainerRestarts).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"deployment_stuck": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DeploymentStuck).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"high_cpu_throttling": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HighCpuThrottling).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"high_cpu_usage": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HighCpuUsage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"high_memory_usage": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HighMemoryUsage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"not_all_pods_ready": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(NotAllPodsReady).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"pending_pods": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(PendingPods).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"pod_stuck_in_terminating": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(PodStuckInTerminating).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"workload_without_ready_pods": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(WorkloadWithoutReadyPods).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"container_restarts":          me.ContainerRestarts,
		"deployment_stuck":            me.DeploymentStuck,
		"high_cpu_throttling":         me.HighCpuThrottling,
		"high_cpu_usage":              me.HighCpuUsage,
		"high_memory_usage":           me.HighMemoryUsage,
		"not_all_pods_ready":          me.NotAllPodsReady,
		"pending_pods":                me.PendingPods,
		"pod_stuck_in_terminating":    me.PodStuckInTerminating,
		"scope":                       me.Scope,
		"workload_without_ready_pods": me.WorkloadWithoutReadyPods,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"container_restarts":          &me.ContainerRestarts,
		"deployment_stuck":            &me.DeploymentStuck,
		"high_cpu_throttling":         &me.HighCpuThrottling,
		"high_cpu_usage":              &me.HighCpuUsage,
		"high_memory_usage":           &me.HighMemoryUsage,
		"not_all_pods_ready":          &me.NotAllPodsReady,
		"pending_pods":                &me.PendingPods,
		"pod_stuck_in_terminating":    &me.PodStuckInTerminating,
		"scope":                       &me.Scope,
		"workload_without_ready_pods": &me.WorkloadWithoutReadyPods,
	})
}
