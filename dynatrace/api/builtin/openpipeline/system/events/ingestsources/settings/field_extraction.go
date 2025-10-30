/**
* @license
* Copyright 2025 Dynatrace LLC
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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FieldExtraction struct {
	Exclude []string               `json:"exclude,omitempty"` // Fields
	Include FieldExtractionEntries `json:"include,omitempty"` // Fields
	Type    FieldExtractionType    `json:"type"`              // Fields Extraction type. Possible Values: `exclude`, `include`, `includeAll`.
}

func (me *FieldExtraction) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"exclude": {
			Type:        schema.TypeSet,
			Description: "Fields",
			Optional:    true, // precondition
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"include": {
			Type:        schema.TypeList,
			Description: "Fields",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FieldExtractionEntries).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Fields Extraction type. Possible Values: `exclude`, `include`, `includeAll`.",
			Required:    true,
		},
	}
}

func (me *FieldExtraction) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"exclude": me.Exclude,
		"include": me.Include,
		"type":    me.Type,
	})
}

func (me *FieldExtraction) HandlePreconditions() error {
	// ---- Exclude []string -> {"expectedValue":"exclude","property":"type","type":"EQUALS"}
	// ---- Include FieldExtractionEntries -> {"expectedValue":"include","property":"type","type":"EQUALS"}
	return nil
}

func (me *FieldExtraction) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"exclude": &me.Exclude,
		"include": &me.Include,
		"type":    &me.Type,
	})
}
