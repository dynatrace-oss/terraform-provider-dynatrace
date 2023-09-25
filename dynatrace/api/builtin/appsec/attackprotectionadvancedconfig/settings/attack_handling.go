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

package attackprotectionadvancedconfig

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AttackHandling struct {
	BlockingStrategy BlockingStrategy `json:"blockingStrategy"` // Possible Values: `BLOCK`, `MONITOR`, `OFF`
}

func (me *AttackHandling) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"blocking_strategy": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BLOCK`, `MONITOR`, `OFF`",
			Required:    true,
		},
	}
}

func (me *AttackHandling) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"blocking_strategy": me.BlockingStrategy,
	})
}

func (me *AttackHandling) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"blocking_strategy": &me.BlockingStrategy,
	})
}
