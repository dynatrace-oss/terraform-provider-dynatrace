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

package imsbridges

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type QueueManager struct {
	Name              string   `json:"name"`
	QueueManagerQueue []string `json:"queueManagerQueue,omitempty"`
}

type QueueManagers []*QueueManager

func (me *QueueManager) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the queue manager",
		},
		"queue_manager_queue": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Queue(s) that belong to the queue manager",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *QueueManagers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"queue_manager": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Queue manager definition for IMS bridge",
			Elem:        &schema.Resource{Schema: new(QueueManager).Schema()},
		},
	}
}

func (me *QueueManager) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":                me.Name,
		"queue_manager_queue": me.QueueManagerQueue,
	})
}

func (me *QueueManager) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":                &me.Name,
		"queue_manager_queue": &me.QueueManagerQueue,
	})
}

func (me QueueManagers) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("queue_manager", me)
}

func (me *QueueManagers) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("queue_manager", me)
}
