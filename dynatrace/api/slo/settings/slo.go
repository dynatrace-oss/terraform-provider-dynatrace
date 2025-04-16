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

package slo

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SchemaVersion = "1.0.4"

type SLO struct {
	Name         string        `json:"name" minlength:"1" maxlength:"200"`
	Description  *string       `json:"description,omitempty" minlength:"1" maxlength:"250"`
	SliReference *SliReference `json:"sliReference,omitempty"`
	CustomSli    *CustomSli    `json:"customSli,omitempty"`
	Criteria     Criteria      `json:"criteria"`
	Tags         []string      `json:"tags,omitempty"`
}

func (me *SLO) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Description:      "Name of the SLO",
			Required:         true,
			ValidateDiagFunc: Validate(ValidateMinLength(1), ValidateMaxLength(200)),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the SLO",
			Optional:         true,
			ValidateDiagFunc: Validate(ValidateMinLength(1), ValidateMaxLength(250)),
		},
		"sli_reference": {
			Type:        schema.TypeList,
			Description: "SLI reference of the SLO",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(SliReference).Schema()},
		},
		"custom_sli": {
			Type:        schema.TypeList,
			Description: "Custom SLI of the SLO",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(CustomSli).Schema()},
		},
		"criteria": {
			Type:        schema.TypeList,
			Description: "Criteria of the SLO",
			MinItems:    1,
			MaxItems:    1,
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Criteria).Schema()},
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "Tags of the SLO. Example: `Stage:DEV`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *SLO) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":          me.Name,
		"description":   me.Description,
		"sli_reference": me.SliReference,
		"custom_sli":    me.CustomSli,
		"criteria":      me.Criteria,
		"tags":          me.Tags,
	})
}

func (me *SLO) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":          &me.Name,
		"description":   &me.Description,
		"sli_reference": &me.SliReference,
		"custom_sli":    &me.CustomSli,
		"criteria":      &me.Criteria,
		"tags":          &me.Tags,
	})
}
