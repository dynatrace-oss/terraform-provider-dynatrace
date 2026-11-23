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

package primarygrailtag

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PrimaryGrailTagSourceFields []*PrimaryGrailTagSourceField

func (me *PrimaryGrailTagSourceFields) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source_field": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(PrimaryGrailTagSourceField).Schema()},
		},
	}
}

func (me PrimaryGrailTagSourceFields) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("source_field", me)
}

func (me *PrimaryGrailTagSourceFields) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("source_field", me)
}

type PrimaryGrailTagSourceField struct {
	FieldExpression string `json:"fieldExpression"` // A DQL field access expression identifying the source field
}

func (me *PrimaryGrailTagSourceField) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_expression": {
			Type:        schema.TypeString,
			Description: "A DQL field access expression identifying the source field",
			Required:    true,
		},
	}
}

func (me *PrimaryGrailTagSourceField) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"field_expression": me.FieldExpression,
	})
}

func (me *PrimaryGrailTagSourceField) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"field_expression": &me.FieldExpression,
	})
}
