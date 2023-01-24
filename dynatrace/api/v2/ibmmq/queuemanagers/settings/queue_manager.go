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
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// QueueManager TODO: documentation
type QueueManager struct {
	Name          string        `json:"name"`
	Clusters      []string      `json:"clusters"`
	AliasQueues   AliasQueues   `json:"aliasQueues,omitempty"`
	RemoteQueues  RemoteQueues  `json:"remoteQueues,omitempty"`
	ClusterQueues ClusterQueues `json:"clusterQueues,omitempty"`
}

func (me *QueueManager) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the queue manager",
		},
		"clusters": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Name of the cluster(s) this queue manager is part of",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"alias_queues": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "The alias queues in the queue manager",
			Elem: &schema.Resource{
				Schema: new(AliasQueues).Schema(),
			},
		},
		"remote_queues": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "The alias queues in the queue manager",
			Elem: &schema.Resource{
				Schema: new(RemoteQueues).Schema(),
			},
		},
		"cluster_queues": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "The alias queues in the queue manager",
			Elem: &schema.Resource{
				Schema: new(ClusterQueues).Schema(),
			},
		},
	}
}

func (me *QueueManager) EnsurePredictableOrder() {
	if len(me.AliasQueues) > 0 {
		conds := []*AliasQueue{}
		condStrings := sort.StringSlice{}
		for _, entry := range me.AliasQueues {
			condBytes, _ := json.Marshal(entry)
			condStrings = append(condStrings, string(condBytes))
		}
		condStrings.Sort()
		for _, condString := range condStrings {
			cond := AliasQueue{}
			json.Unmarshal([]byte(condString), &cond)
			conds = append(conds, &cond)
		}
		me.AliasQueues = conds
	}
	if len(me.RemoteQueues) > 0 {
		conds := []*RemoteQueue{}
		condStrings := sort.StringSlice{}
		for _, entry := range me.RemoteQueues {
			condBytes, _ := json.Marshal(entry)
			condStrings = append(condStrings, string(condBytes))
		}
		condStrings.Sort()
		for _, condString := range condStrings {
			cond := RemoteQueue{}
			json.Unmarshal([]byte(condString), &cond)
			conds = append(conds, &cond)
		}
		me.RemoteQueues = conds
	}
	if len(me.ClusterQueues) > 0 {
		conds := []*ClusterQueue{}
		condStrings := sort.StringSlice{}
		for _, entry := range me.ClusterQueues {
			condBytes, _ := json.Marshal(entry)
			condStrings = append(condStrings, string(condBytes))
		}
		condStrings.Sort()
		for _, condString := range condStrings {
			cond := ClusterQueue{}
			json.Unmarshal([]byte(condString), &cond)
			conds = append(conds, &cond)
		}
		me.ClusterQueues = conds
	}
}

func (me *QueueManager) MarshalHCL(properties hcl.Properties) error {
	me.EnsurePredictableOrder()

	return properties.EncodeAll(map[string]any{
		"name":           me.Name,
		"clusters":       me.Clusters,
		"alias_queues":   me.AliasQueues,
		"remote_queues":  me.RemoteQueues,
		"cluster_queues": me.ClusterQueues,
	})
}

func (me *QueueManager) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":           &me.Name,
		"clusters":       &me.Clusters,
		"alias_queues":   &me.AliasQueues,
		"remote_queues":  &me.RemoteQueues,
		"cluster_queues": &me.ClusterQueues,
	})
}
