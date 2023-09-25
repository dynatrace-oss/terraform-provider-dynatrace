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

package log

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PathDefinitions []*PathDefinition

func (me *PathDefinitions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"definition": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A list of filtering criteria for log path. If several criteria are specified, the OR logic applies.",
			Elem:        &schema.Resource{Schema: new(PathDefinition).Schema()},
		},
	}
}

func (me PathDefinitions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("definition", me)
}

func (me *PathDefinitions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("definition", me)
}

// PathDefinitions Parameters of a filter of a calculated log metric.
type PathDefinition struct {
	Definition string `json:"definition"` // The path to the required log path. If the type is set to WILDCARD, it may contain wildcard characters (*).
	Type       Type   `json:"type"`       // The type of the log path definition: fixed or an expression with wildcards.
}

func (me *PathDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"definition": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The path to the required log path. If the type is set to WILDCARD, it may contain wildcard characters (*).",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the log path definition: fixed or an expression with wildcards.",
		},
	}
}

func (me *PathDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"definition": me.Definition,
		"type":       me.Type,
	})
}

func (me *PathDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"definition": &me.Definition,
		"type":       &me.Type,
	})
}

type Type string

var Types = struct {
	Fixed    Type
	Wildcard Type
}{
	"FIXED",
	"WILDCARD",
}
