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

package customlogsourcesettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomLogSource struct {
	Accept_binary         *bool                          `json:"accept-binary,omitempty"`         // Accept binary content
	Type                  LogSourceType                  `json:"type"`                            // Possible Values: `LOG_PATH_PATTERN`, `WINDOWS_EVENT_LOG`
	Values                []string                       `json:"values,omitempty"`                // (Required attribute for cluster v1.291 and under) It might be either an absolute path to log(s) with optional wildcards or Windows Event Log name.
	Values_and_enrichment CustomLogSourceWithEnrichments `json:"values-and-enrichment,omitempty"` // (Required attribute for cluster v1.292+) It might be either an absolute path to log(s) with optional wildcards or Windows Event Log name.
}

func (me *CustomLogSource) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"accept_binary": {
			Type:        schema.TypeBool,
			Description: "Accept binary content",
			Optional:    true, // nullable
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `LOG_PATH_PATTERN`, `WINDOWS_EVENT_LOG`",
			Required:    true,
		},
		"values": {
			Type:        schema.TypeSet,
			Description: "(Required attribute for cluster v1.291 and under) It might be either an absolute path to log(s) with optional wildcards or Windows Event Log name.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"values_and_enrichment": {
			Type:        schema.TypeList,
			Description: "(Required attribute for cluster v1.292+) It might be either an absolute path to log(s) with optional wildcards or Windows Event Log name.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(CustomLogSourceWithEnrichments).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *CustomLogSource) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"accept_binary":         me.Accept_binary,
		"type":                  me.Type,
		"values":                me.Values,
		"values_and_enrichment": me.Values_and_enrichment,
	})
}

func (me *CustomLogSource) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"accept_binary":         &me.Accept_binary,
		"type":                  &me.Type,
		"values":                &me.Values,
		"values_and_enrichment": &me.Values_and_enrichment,
	})
}
