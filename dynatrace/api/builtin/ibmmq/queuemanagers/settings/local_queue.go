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

type LocalQueue struct {
	LocalQueueName string `json:"localQueue"`
}

func (me *LocalQueue) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"local_queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the local queue",
		},
	}
}

func (me *LocalQueue) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"local_queue_name": me.LocalQueueName,
	})
}

func (me *LocalQueue) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"local_queue_name": &me.LocalQueueName,
	})
}
