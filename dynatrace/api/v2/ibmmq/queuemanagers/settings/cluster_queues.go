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

package queuemanagers

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ClusterQueue struct {
	LocalQueueName    string   `json:"localQueue"`
	ClusterVisibility []string `json:"clusterVisibility,omitempty"`
}

type ClusterQueues []*ClusterQueue

func (me *ClusterQueue) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"local_queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the local queue",
		},
		"cluster_visibility": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Name of the cluster(s) this local queue should be visible in",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *ClusterQueues) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_queue": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Cluster queue definitions for queue manager",
			Elem:        &schema.Resource{Schema: new(ClusterQueue).Schema()},
		},
	}
}

func (me *ClusterQueue) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"local_queue_name":   me.LocalQueueName,
		"cluster_visibility": me.ClusterVisibility,
	})
}

func (me *ClusterQueue) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"local_queue_name":   &me.LocalQueueName,
		"cluster_visibility": &me.ClusterVisibility,
	})
}

func (me ClusterQueues) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("cluster_queue", me)
}

func (me *ClusterQueues) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("cluster_queue", me)
}
