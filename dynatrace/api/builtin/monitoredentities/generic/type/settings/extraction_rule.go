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

package generictype

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ExtractionRules []*ExtractionRule

func (me *ExtractionRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ExtractionRule).Schema()},
		},
	}
}

func (me ExtractionRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *ExtractionRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

// Entity extraction rule. An extraction rule defines which sources are evaluated for extracting an entity. If a source matches the specified filters, an entity is extracted.
type ExtractionRule struct {
	Attributes          AttributeEntries `json:"attributes"`                    // All attribute extraction rules will be applied and found attributes will be added to the extracted type.
	IconPattern         *string          `json:"iconPattern,omitempty"`         // Define a pattern which is used to set the icon attribute of the entity. The extracted values must reference barista icon ids. You may define placeholders referencing data source dimensions.
	IdPattern           string           `json:"idPattern"`                     // ID patterns are comprised of static text and placeholders referring to dimensions in the ingest data. An ID pattern **must** contain at least one placeholder to ensure that different entities will be created.. Take care that the pattern results in the same ID for the same entity. For example, using timestamp or counter-like dimensions as part of the ID would lead to the creation of new entities for each ingest data and is strongly discouraged!\n\nEach dimension key referred to by an identifier placeholder must be present in order to extract an entity. If any dimension key referred to in the identifier is missing, the rule will not be considered for evaluation. If you have cases where you still want to extract the same entity type but have differently named keys, consider creating multiple rules extracting the same entity type. In this case take care that each ID pattern evaluates to the same value if the same entity should be extracted.
	InstanceNamePattern *string          `json:"instanceNamePattern,omitempty"` // Define a pattern which is used to set the name attribute of the entity. You may define placeholders referencing data source dimensions.
	RequiredDimensions  DimensionFilters `json:"requiredDimensions"`            // In addition to the dimensions already referred to in the ID pattern, you may specify additional dimensions which must be present in order to evaluate this rule.
	Role                *string          `json:"role,omitempty"`                // If you want to extract multiple entities of the same type from a single ingest line you need to define multiple rules with different roles.
	Sources             SourceFilters    `json:"sources"`                       // Specify all sources which should be evaluated for this rule. A rule is evaluated if any of the specified source filters match.
}

func (me *ExtractionRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attributes": {
			Type:        schema.TypeList,
			Description: "All attribute extraction rules will be applied and found attributes will be added to the extracted type.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(AttributeEntries).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"icon_pattern": {
			Type:        schema.TypeString,
			Description: "Define a pattern which is used to set the icon attribute of the entity. The extracted values must reference barista icon ids. You may define placeholders referencing data source dimensions.",
			Optional:    true,
		},
		"id_pattern": {
			Type:        schema.TypeString,
			Description: "ID patterns are comprised of static text and placeholders referring to dimensions in the ingest data. An ID pattern **must** contain at least one placeholder to ensure that different entities will be created.. Take care that the pattern results in the same ID for the same entity. For example, using timestamp or counter-like dimensions as part of the ID would lead to the creation of new entities for each ingest data and is strongly discouraged!\n\nEach dimension key referred to by an identifier placeholder must be present in order to extract an entity. If any dimension key referred to in the identifier is missing, the rule will not be considered for evaluation. If you have cases where you still want to extract the same entity type but have differently named keys, consider creating multiple rules extracting the same entity type. In this case take care that each ID pattern evaluates to the same value if the same entity should be extracted.",
			Required:    true,
		},
		"instance_name_pattern": {
			Type:        schema.TypeString,
			Description: "Define a pattern which is used to set the name attribute of the entity. You may define placeholders referencing data source dimensions.",
			Optional:    true,
		},
		"required_dimensions": {
			Type:        schema.TypeList,
			Description: "In addition to the dimensions already referred to in the ID pattern, you may specify additional dimensions which must be present in order to evaluate this rule.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(DimensionFilters).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"role": {
			Type:        schema.TypeString,
			Description: "If you want to extract multiple entities of the same type from a single ingest line you need to define multiple rules with different roles.",
			Optional:    true,
		},
		"sources": {
			Type:        schema.TypeList,
			Description: "Specify all sources which should be evaluated for this rule. A rule is evaluated if any of the specified source filters match.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(SourceFilters).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *ExtractionRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attributes":            me.Attributes,
		"icon_pattern":          me.IconPattern,
		"id_pattern":            me.IdPattern,
		"instance_name_pattern": me.InstanceNamePattern,
		"required_dimensions":   me.RequiredDimensions,
		"role":                  me.Role,
		"sources":               me.Sources,
	})
}

func (me *ExtractionRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attributes":            &me.Attributes,
		"icon_pattern":          &me.IconPattern,
		"id_pattern":            &me.IdPattern,
		"instance_name_pattern": &me.InstanceNamePattern,
		"required_dimensions":   &me.RequiredDimensions,
		"role":                  &me.Role,
		"sources":               &me.Sources,
	})
}
