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

package queuesharinggroups

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Filters TODO: documentation
type QueueSharingGroup struct {
	Name          string   `json:"name"`
	QueueManagers []string `json:"queueManagers,omitempty"`
	SharedQueues  []string `json:"sharedQueues,omitempty"`
}

func (me *QueueSharingGroup) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the queue sharing group",
		},
		"queue_managers": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Queue manager(s) that belong to the queue sharing group",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"shared_queues": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Shared queue(s) that belong to the queue sharing group",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *QueueSharingGroup) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":           me.Name,
		"queue_managers": me.QueueManagers,
		"shared_queues":  me.SharedQueues,
	})
}

func (me *QueueSharingGroup) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":           &me.Name,
		"queue_managers": &me.QueueManagers,
		"shared_queues":  &me.SharedQueues,
	})
}
