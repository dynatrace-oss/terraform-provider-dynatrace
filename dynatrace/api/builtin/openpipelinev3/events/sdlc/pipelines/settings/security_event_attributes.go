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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecurityEventAttributes struct {
	FieldExtraction *FieldExtraction `json:"fieldExtraction"` // Field Extraction
}

func (me *SecurityEventAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_extraction": {
			Type:        schema.TypeList,
			Description: "Field Extraction",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FieldExtraction).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *SecurityEventAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"field_extraction": me.FieldExtraction,
	})
}

func (me *SecurityEventAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"field_extraction": &me.FieldExtraction,
	})
}
