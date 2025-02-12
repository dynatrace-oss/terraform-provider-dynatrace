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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID               string                `json:"-" scope:"applicationId"`     // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	CacheControlHeaders         *CacheControlHeaders  `json:"cacheControlHeaders"`         // Cache control headers
	MonitoringCodeSourceSection *MonitoringCodeSource `json:"monitoringCodeSourceSection"` // Real User Monitoring code source
	SnippetFormat               *SnippetFormat        `json:"snippetFormat"`               // *Code Snippet:* OneAgent injects an inline script that initializes Dynatrace and dynamically downloads the monitoring code into your application. Use when you want to inject the monitoring code in deferred mode.\n\n *Inline Code:* OneAgent injects the configuration and the monitoring code inline into your application. Use this injection type when you need to keep the number of web requests at a minimum.\n\n *OneAgent JavaScript Tag:* OneAgent injects a JavaScript tag into your application, containing the configuration and a link to the monitoring code. This is our default injection type, since it's most versatile.\n\nCompare the different [injection formats](https://dt-url.net/vx5g0ptn).
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
		"cache_control_headers": {
			Type:        schema.TypeList,
			Description: "Cache control headers",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(CacheControlHeaders).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"monitoring_code_source_section": {
			Type:        schema.TypeList,
			Description: "Real User Monitoring code source",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(MonitoringCodeSource).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"snippet_format": {
			Type:        schema.TypeList,
			Description: "*Code Snippet:* OneAgent injects an inline script that initializes Dynatrace and dynamically downloads the monitoring code into your application. Use when you want to inject the monitoring code in deferred mode.\n\n *Inline Code:* OneAgent injects the configuration and the monitoring code inline into your application. Use this injection type when you need to keep the number of web requests at a minimum.\n\n *OneAgent JavaScript Tag:* OneAgent injects a JavaScript tag into your application, containing the configuration and a link to the monitoring code. This is our default injection type, since it's most versatile.\n\nCompare the different [injection formats](https://dt-url.net/vx5g0ptn).",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SnippetFormat).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":                 me.ApplicationID,
		"cache_control_headers":          me.CacheControlHeaders,
		"monitoring_code_source_section": me.MonitoringCodeSourceSection,
		"snippet_format":                 me.SnippetFormat,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":                 &me.ApplicationID,
		"cache_control_headers":          &me.CacheControlHeaders,
		"monitoring_code_source_section": &me.MonitoringCodeSourceSection,
		"snippet_format":                 &me.SnippetFormat,
	})
}
