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

type CodeSnippet struct {
	CodeSnippetType CodeSnippetType `json:"codeSnippetType"` // Possible Values: `DEFERRED`, `SYNCHRONOUSLY`
}

func (me *CodeSnippet) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code_snippet_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DEFERRED`, `SYNCHRONOUSLY`",
			Required:    true,
		},
	}
}

func (me *CodeSnippet) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"code_snippet_type": me.CodeSnippetType,
	})
}

func (me *CodeSnippet) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"code_snippet_type": &me.CodeSnippetType,
	})
}
