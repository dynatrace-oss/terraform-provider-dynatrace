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

type SliReference struct {
	TemplateId string                `json:"templateId" minlength:"1" maxlength:"800"`
	Variables  SliReferenceVariables `json:"variables"`
}

func (me *SliReference) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"template_id": {
			Type:             schema.TypeString,
			Description:      "Template ID of the SLI reference",
			Required:         true,
			ValidateDiagFunc: Validate(ValidateMinLength(1), ValidateMaxLength(800)),
		},
		"variables": {
			Type:        schema.TypeList,
			Description: "Variables of the SLI reference",
			MinItems:    1,
			MaxItems:    1,
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SliReferenceVariables).Schema()},
		},
	}
}

func (me *SliReference) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"template_id": me.TemplateId,
		"variables":   me.Variables,
	})
}

func (me *SliReference) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"template_id": &me.TemplateId,
		"variables":   &me.Variables,
	})
}
