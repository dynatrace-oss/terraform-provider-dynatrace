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
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IMSBridge struct {
	Name          string        `json:"name"`
	QueueManagers QueueManagers `json:"queueManagers,omitempty"`
}

func (me *IMSBridge) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the IMS bridge",
		},
		"queue_managers": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "Queue manager(s) that belong to the IMS bridge",
			Elem: &schema.Resource{
				Schema: new(QueueManagers).Schema(),
			},
		},
	}
}

func (me *IMSBridge) EnsurePredictableOrder() {
	if len(me.QueueManagers) > 0 {
		conds := []*QueueManager{}
		condStrings := sort.StringSlice{}
		for _, entry := range me.QueueManagers {
			condBytes, _ := json.Marshal(entry)
			condStrings = append(condStrings, string(condBytes))
		}
		condStrings.Sort()
		for _, condString := range condStrings {
			cond := QueueManager{}
			json.Unmarshal([]byte(condString), &cond)
			conds = append(conds, &cond)
		}
		me.QueueManagers = conds
	}
}

func (me *IMSBridge) MarshalHCL(properties hcl.Properties) error {
	me.EnsurePredictableOrder()

	return properties.EncodeAll(map[string]any{
		"name":           me.Name,
		"queue_managers": me.QueueManagers,
	})
}

func (me *IMSBridge) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":           &me.Name,
		"queue_managers": &me.QueueManagers,
	})
}
