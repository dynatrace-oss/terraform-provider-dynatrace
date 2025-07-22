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

type Settings struct {
	ApplicationID            string                    `json:"-" scope:"applicationId"`            // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	CodeSnippet              *CodeSnippet              `json:"codeSnippet"`                        // Code snippet is a piece of inline code that implements basic functionality and loads the full functionality either synchronously or deferred. Even though it implements an update mechanism, regular updates are still required to guarantee compatibility.
	JavascriptTag            *JavascriptTag            `json:"javascriptTag"`                      // JavaScript tag references an external file containing monitoring code and configuration. Due to its dynamic update mechanism, it is recommended for most use cases.
	OneagentJavascriptTag    *OneagentJavascriptTag    `json:"oneagentJavascriptTag,omitempty"`    // OneAgent JavaScript tag includes configuration and a reference to an external file containing the monitoring code. It needs to be updated after configuration changes and monitoring code updates.
	OneagentJavascriptTagSRI *OneagentJavascriptTagSRI `json:"oneagentJavascriptTagSRI,omitempty"` // OneAgent JavaScript tag with SRI includes configuration, a reference to an external file containing the monitoring code, and a hash that allows the browser to verify the integrity of the monitoring code before executing it. It needs to be updated after configuration changes and monitoring code updates.
}

func (me *Settings) Name() string {
	return me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
		},
		"code_snippet": {
			Type:        schema.TypeList,
			Description: "Code snippet is a piece of inline code that implements basic functionality and loads the full functionality either synchronously or deferred. Even though it implements an update mechanism, regular updates are still required to guarantee compatibility.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(CodeSnippet).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"javascript_tag": {
			Type:        schema.TypeList,
			Description: "JavaScript tag references an external file containing monitoring code and configuration. Due to its dynamic update mechanism, it is recommended for most use cases.",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(JavascriptTag).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"oneagent_javascript_tag": {
			Type:        schema.TypeList,
			Description: "OneAgent JavaScript tag includes configuration and a reference to an external file containing the monitoring code. It needs to be updated after configuration changes and monitoring code updates.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(OneagentJavascriptTag).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"oneagent_javascript_tag_sri": {
			Type:        schema.TypeList,
			Description: "OneAgent JavaScript tag with SRI includes configuration, a reference to an external file containing the monitoring code, and a hash that allows the browser to verify the integrity of the monitoring code before executing it. It needs to be updated after configuration changes and monitoring code updates.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(OneagentJavascriptTagSRI).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	m := map[string]any{
		"application_id": me.ApplicationID,
		"code_snippet":   me.CodeSnippet,
		"javascript_tag": me.JavascriptTag,
	}
	if !me.OneagentJavascriptTag.IsEmpty() {
		m["oneagent_javascript_tag"] = me.OneagentJavascriptTag
	}
	if !me.OneagentJavascriptTagSRI.IsEmpty() {
		m["oneagent_javascript_tag_sri"] = me.OneagentJavascriptTagSRI
	}
	return properties.EncodeAll(m)
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":              &me.ApplicationID,
		"code_snippet":                &me.CodeSnippet,
		"javascript_tag":              &me.JavascriptTag,
		"oneagent_javascript_tag":     &me.OneagentJavascriptTag,
		"oneagent_javascript_tag_sri": &me.OneagentJavascriptTagSRI,
	})
}
