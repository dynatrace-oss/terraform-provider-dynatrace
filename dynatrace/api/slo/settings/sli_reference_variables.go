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

type SliReferenceVariables []*SliReferenceVariable

func (me *SliReferenceVariables) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sli_reference_variable": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SliReferenceVariable).Schema()},
		},
	}
}

func (me SliReferenceVariables) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("sli_reference_variable", me)
}

func (me *SliReferenceVariables) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("sli_reference_variable", me)
}

type SliReferenceVariable struct {
	Name  string `json:"name" maxlength:"60"`
	Value string `json:"value" maxlength:"1000"`
}

func (me *SliReferenceVariable) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Description:      "Name of the SLI reference variable. Example: `hostIds`",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(60),
		},
		"value": {
			Type:             schema.TypeString,
			Description:      "Value of the SLI reference variable. Example: `HOST-123456789ABCDEFG`",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(1000),
		},
	}
}

func (me *SliReferenceVariable) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"value": me.Value,
	})
}

func (me *SliReferenceVariable) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"value": &me.Value,
	})
}
