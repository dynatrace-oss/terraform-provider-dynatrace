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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Constraints []*Constraint

func (me *Constraints) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"constraint": {
			Type:        schema.TypeSet,
			Description: "The network availability monitor constraint",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Constraint).Schema()},
		},
	}
}

func (me Constraints) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("constraint", me)
}

func (me *Constraints) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("constraint", me); err != nil {
		return err
	}
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/895
	// Only known workaround is to ignore these blocks
	newEntries := Constraints{}
	for _, entry := range *me {
		if entry.Type != "" || entry.Properties != nil {
			newEntries = append(newEntries, entry)
		}
	}
	*me = newEntries
	return nil
}

type Constraint struct {
	Type       string            `json:"type"`       // Constraint type
	Properties map[string]string `json:"properties"` // Key/value pairs of constraint properties
}

func (me *Constraint) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Constraint type",
			Required:    true,
		},
		"properties": {
			Type:        schema.TypeMap,
			Description: "Key/value pairs of constraint properties",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me Constraint) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":       me.Type,
		"properties": me.Properties,
	})
}

func (me *Constraint) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":       &me.Type,
		"properties": &me.Properties,
	})
}
