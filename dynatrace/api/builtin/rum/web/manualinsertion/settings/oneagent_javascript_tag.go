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

package manualinsertion

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (me *OneagentJavascriptTag) IsEmpty() bool {
	return me == nil || me.ScriptExecutionAttribute == nil
}

type OneagentJavascriptTag struct {
	ScriptExecutionAttribute *ScriptExecutionAttribute `json:"scriptExecutionAttribute,omitempty"` // Possible Values: `Async`, `Defer`, `None`
}

func (me *OneagentJavascriptTag) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"script_execution_attribute": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Async`, `Defer`, `None`",
			Optional:    true, // nullable
		},
	}
}

func (me *OneagentJavascriptTag) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"script_execution_attribute": me.ScriptExecutionAttribute,
	})
}

func (me *OneagentJavascriptTag) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"script_execution_attribute": &me.ScriptExecutionAttribute,
	})
}
