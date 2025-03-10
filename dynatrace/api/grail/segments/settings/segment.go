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

package segments

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SchemaVersion = 3

type Segment struct {
	Name        string     `json:"name" maxlength:"100"`
	Description string     `json:"description,omitempty"`
	Variables   *Variables `json:"variables,omitempty"`
	IsPublic    bool       `json:"isPublic"`
	Includes    Includes   `json:"includes,omitempty"`
}

func (me *Segment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Description:      "Name of the filter-segment",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(200),
		},
		"description": {
			Type:             schema.TypeString,
			Description:      "Description of the filter-segment",
			Optional:         true,
			ValidateDiagFunc: ValidateMaxLength(200),
		},
		"variables": {
			Type:        schema.TypeList,
			Description: "Variables of the filter-segment",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Variables).Schema()},
		},
		"is_public": {
			Type:        schema.TypeBool,
			Description: "Indicates if the filter-segment is publicly accessible within the tenant",
			Required:    true,
		},
		"includes": {
			Type:        schema.TypeList,
			Description: "List of includes of the filter-segment",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Includes).Schema()},
		},
	}
}

func (me *Segment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":        me.Name,
		"description": me.Description,
		"variables":   me.Variables,
		"is_public":   me.IsPublic,
		"includes":    me.Includes,
	})
}

func (me *Segment) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":        &me.Name,
		"description": &me.Description,
		"variables":   &me.Variables,
		"is_public":   &me.IsPublic,
		"includes":    &me.Includes,
	})
}
