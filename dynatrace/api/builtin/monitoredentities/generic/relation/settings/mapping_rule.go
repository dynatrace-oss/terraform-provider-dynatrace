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

package relation

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MappingRules []*MappingRule

func (me *MappingRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mapping_rule": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(MappingRule).Schema()},
		},
	}
}

func (me MappingRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("mapping_rule", me)
}

func (me *MappingRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("mapping_rule", me)
}

type MappingRule struct {
	DestinationProperty       string        `json:"destinationProperty"`       // The case-sensitive name of a property of the destination type.
	DestinationTransformation Normalization `json:"destinationTransformation"` // Possible Values: `Leavetextas_is`, `Tolowercase`, `Touppercase`
	SourceProperty            string        `json:"sourceProperty"`            // The case-sensitive name of a property of the source type.
	SourceTransformation      Normalization `json:"sourceTransformation"`      // Possible Values: `Leavetextas_is`, `Tolowercase`, `Touppercase`
}

func (me *MappingRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"destination_property": {
			Type:        schema.TypeString,
			Description: "The case-sensitive name of a property of the destination type.",
			Required:    true,
		},
		"destination_transformation": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Leavetextas_is`, `Tolowercase`, `Touppercase`",
			Required:    true,
		},
		"source_property": {
			Type:        schema.TypeString,
			Description: "The case-sensitive name of a property of the source type.",
			Required:    true,
		},
		"source_transformation": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Leavetextas_is`, `Tolowercase`, `Touppercase`",
			Required:    true,
		},
	}
}

func (me *MappingRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"destination_property":       me.DestinationProperty,
		"destination_transformation": me.DestinationTransformation,
		"source_property":            me.SourceProperty,
		"source_transformation":      me.SourceTransformation,
	})
}

func (me *MappingRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"destination_property":       &me.DestinationProperty,
		"destination_transformation": &me.DestinationTransformation,
		"source_property":            &me.SourceProperty,
		"source_transformation":      &me.SourceTransformation,
	})
}
