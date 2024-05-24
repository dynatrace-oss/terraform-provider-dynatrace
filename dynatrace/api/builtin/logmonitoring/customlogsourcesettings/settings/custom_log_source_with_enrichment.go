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

type CustomLogSourceWithEnrichments []*CustomLogSourceWithEnrichment

func (me *CustomLogSourceWithEnrichments) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_log_source_with_enrichment": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CustomLogSourceWithEnrichment).Schema()},
		},
	}
}

func (me CustomLogSourceWithEnrichments) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("custom_log_source_with_enrichment", me)
}

func (me *CustomLogSourceWithEnrichments) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("custom_log_source_with_enrichment", me)
}

type CustomLogSourceWithEnrichment struct {
	Enrichment Enrichments `json:"enrichment,omitempty"` // Optional field that allows to define attributes that will enrich logs. ${N} can be used in attribute value to expand the value matched by wildcards where N denotes the number of the wildcard the expand
	Path       string      `json:"path"`                 // Values
}

func (me *CustomLogSourceWithEnrichment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enrichment": {
			Type:        schema.TypeList,
			Description: "Optional field that allows to define attributes that will enrich logs. ${N} can be used in attribute value to expand the value matched by wildcards where N denotes the number of the wildcard the expand",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Enrichments).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"path": {
			Type:        schema.TypeString,
			Description: "Values",
			Required:    true,
		},
	}
}

func (me *CustomLogSourceWithEnrichment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enrichment": me.Enrichment,
		"path":       me.Path,
	})
}

func (me *CustomLogSourceWithEnrichment) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enrichment": &me.Enrichment,
		"path":       &me.Path,
	})
}
