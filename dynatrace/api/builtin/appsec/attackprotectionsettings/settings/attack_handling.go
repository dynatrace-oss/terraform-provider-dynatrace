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

package attackprotectionsettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Default attack handling. Default settings for handling attacks.
type AttackHandling struct {
	BlockingStrategyDotNet BlockingStrategy  `json:"blockingStrategyDotNet"`       // (v1.290) Possible Values: `BLOCK`, `MONITOR`, `OFF`
	BlockingStrategyGo     *BlockingStrategy `json:"blockingStrategyGo,omitempty"` // Possible Values: `BLOCK`, `MONITOR`, `OFF`
	BlockingStrategyJava   BlockingStrategy  `json:"blockingStrategyJava"`         // Possible Values: `BLOCK`, `MONITOR`, `OFF`
}

func (me *AttackHandling) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"blocking_strategy_dot_net": {
			Type:        schema.TypeString,
			Description: "(v1.290) Possible Values: `BLOCK`, `MONITOR`, `OFF`",
			Optional:    true, // nullable
			Default:     "OFF",
		},
		"blocking_strategy_go": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BLOCK`, `MONITOR`, `OFF`",
			Optional:    true, // nullable
			Default:     "OFF",
		},
		"blocking_strategy_java": {
			Type:        schema.TypeString,
			Description: "Possible Values: `BLOCK`, `MONITOR`, `OFF`",
			Required:    true,
		},
	}
}

func (me *AttackHandling) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"blocking_strategy_dot_net": me.BlockingStrategyDotNet,
		"blocking_strategy_go":      me.BlockingStrategyGo,
		"blocking_strategy_java":    me.BlockingStrategyJava,
	})
}

func (me *AttackHandling) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"blocking_strategy_dot_net": &me.BlockingStrategyDotNet,
		"blocking_strategy_go":      &me.BlockingStrategyGo,
		"blocking_strategy_java":    &me.BlockingStrategyJava,
	})
}
