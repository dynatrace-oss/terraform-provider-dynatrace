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

package automaticinjection

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SnippetFormat struct {
	CodeSnippetType          *CodeSnippetType          `json:"codeSnippetType,omitempty"`          // Possible Values: `DEFERRED`, `SYNCHRONOUSLY`
	ScriptExecutionAttribute *ScriptExecutionAttribute `json:"scriptExecutionAttribute,omitempty"` // Possible Values: `Async`, `Defer`, `None`
	SnippetFormat            string                    `json:"snippetFormat"`                      // Snippet format
}

func (me *SnippetFormat) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code_snippet_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DEFERRED`, `SYNCHRONOUSLY`",
			Optional:    true, // precondition
		},
		"script_execution_attribute": {
			Type:        schema.TypeString,
			Description: "Possible Values: `async`, `defer`, `none`",
			Optional:    true, // nullable & precondition
		},
		"snippet_format": {
			Type:        schema.TypeString,
			Description: "Snippet format",
			Required:    true,
		},
	}
}

func (me *SnippetFormat) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"code_snippet_type":          me.CodeSnippetType,
		"script_execution_attribute": me.ScriptExecutionAttribute,
		"snippet_format":             me.SnippetFormat,
	})
}

func (me *SnippetFormat) HandlePreconditions() error {
	if (me.CodeSnippetType == nil) && (string(me.SnippetFormat) == "Code Snippet") {
		return fmt.Errorf("'code_snippet_type' must be specified if 'snippet_format' is set to '%v'", me.SnippetFormat)
	}
	if (me.CodeSnippetType != nil) && (string(me.SnippetFormat) != "Code Snippet") {
		return fmt.Errorf("'code_snippet_type' must not be specified if 'snippet_format' is set to '%v'", me.SnippetFormat)
	}
	if (me.ScriptExecutionAttribute != nil) && (!slices.Contains([]string{"OneAgent JavaScript Tag", "OneAgent JavaScript Tag with SRI"}, string(me.SnippetFormat))) {
		return fmt.Errorf("'script_execution_attribute' must not be specified if 'snippet_format' is set to '%v'", me.SnippetFormat)
	}
	return nil
}

func (me *SnippetFormat) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"code_snippet_type":          &me.CodeSnippetType,
		"script_execution_attribute": &me.ScriptExecutionAttribute,
		"snippet_format":             &me.SnippetFormat,
	})
}
