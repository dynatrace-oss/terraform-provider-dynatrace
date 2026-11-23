/**
* @license
* Copyright 2026 Dynatrace LLC
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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataMaskingAttributes struct {
	DataMaskingFields []string     `json:"dataMaskingFields"` // Data masking fields
	MaskingRules      MaskingRules `json:"maskingRules"`      // Masking rules
}

func (me *DataMaskingAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data_masking_fields": {
			Type:        schema.TypeList,
			Description: "Data masking fields",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"masking_rules": {
			Type:        schema.TypeList,
			Description: "Masking rules",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(MaskingRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *DataMaskingAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"data_masking_fields": me.DataMaskingFields,
		"masking_rules":       me.MaskingRules,
	})
}

func (me *DataMaskingAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"data_masking_fields": &me.DataMaskingFields,
		"masking_rules":       &me.MaskingRules,
	})
}
