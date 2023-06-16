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

type RemoteQueue struct {
	LocalQueueName     string   `json:"localQueue"`
	RemoteQueueName    string   `json:"remoteQueue"`
	RemoteQueueManager string   `json:"remoteQueueManager"`
	ClusterVisibility  []string `json:"clusterVisibility,omitempty"`
}

type RemoteQueues []*RemoteQueue

func (me *RemoteQueue) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"local_queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the local queue",
		},
		"remote_queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the remote queue",
		},
		"remote_queue_manager": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the remote queue manager",
		},
		"cluster_visibility": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Name of the cluster(s) this local definition of the remote queue should be visible in",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *RemoteQueues) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"remote_queue": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Remote queue definitions for queue manager",
			Elem:        &schema.Resource{Schema: new(RemoteQueue).Schema()},
		},
	}
}

func (me *RemoteQueue) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"local_queue_name":     me.LocalQueueName,
		"remote_queue_name":    me.RemoteQueueName,
		"remote_queue_manager": me.RemoteQueueManager,
		"cluster_visibility":   me.ClusterVisibility,
	})
}

func (me *RemoteQueue) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"local_queue_name":     &me.LocalQueueName,
		"remote_queue_name":    &me.RemoteQueueName,
		"remote_queue_manager": &me.RemoteQueueManager,
		"cluster_visibility":   &me.ClusterVisibility,
	})
}

func (me RemoteQueues) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("remote_queue", me)
}

func (me *RemoteQueues) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("remote_queue", me)
}
