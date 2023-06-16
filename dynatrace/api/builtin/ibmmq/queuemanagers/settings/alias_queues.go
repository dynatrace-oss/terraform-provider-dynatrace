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

type AliasQueue struct {
	AliasQueueName    string   `json:"aliasQueue"`
	BaseQueueName     string   `json:"baseQueue"`
	ClusterVisibility []string `json:"clusterVisibility,omitempty"`
}

type AliasQueues []*AliasQueue

func (me *AliasQueue) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alias_queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the alias queue",
		},
		"base_queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the base queue",
		},
		"cluster_visibility": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Name of the cluster(s) this alias should be visible in",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *AliasQueues) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alias_queue": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Alias queue definitions for queue manager",
			Elem:        &schema.Resource{Schema: new(AliasQueue).Schema()},
		},
	}
}

func (me *AliasQueue) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"alias_queue_name":   me.AliasQueueName,
		"base_queue_name":    me.BaseQueueName,
		"cluster_visibility": me.ClusterVisibility,
	})
}

func (me *AliasQueue) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"alias_queue_name":   &me.AliasQueueName,
		"base_queue_name":    &me.BaseQueueName,
		"cluster_visibility": &me.ClusterVisibility,
	})
}

func (me AliasQueues) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("alias_queue", me)
}

func (me *AliasQueues) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("alias_queue", me)
}
