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

package anomalydetectors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ExecutionSettings struct {
	QueryOffset *int `json:"queryOffset,omitempty"` // Minute offset of sliding evaluation window for metrics with latency
}

func (me *ExecutionSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"query_offset": {
			Type:        schema.TypeInt,
			Description: "Minute offset of sliding evaluation window for metrics with latency",
			Optional:    true, // nullable
		},
	}
}

func (me *ExecutionSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"query_offset": me.QueryOffset,
	})
}

func (me *ExecutionSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"query_offset": &me.QueryOffset,
	})
}
