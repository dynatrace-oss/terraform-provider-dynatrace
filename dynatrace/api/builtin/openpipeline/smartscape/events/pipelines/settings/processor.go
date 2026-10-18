/**
* @license
* Copyright 2025 Dynatrace LLC
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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Processors []*Processor

func (me *Processors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Processor).Schema()},
		},
	}
}

func (me Processors) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", me)
}

func (me *Processors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", me)
}

// Processor. Processor definition
type Processor struct {
	Description     string                     `json:"description"`
	Dql             *DqlAttributes             `json:"dql,omitempty"`             // DQL processor attributes
	Enabled         bool                       `json:"enabled"`                   // This setting is enabled (`true`) or disabled (`false`)
	FieldsAdd       *FieldsAddAttributes       `json:"fieldsAdd,omitempty"`       // Fields add processor attributes
	FieldsRemove    *FieldsRemoveAttributes    `json:"fieldsRemove,omitempty"`    // Fields remove processor attributes
	FieldsRename    *FieldsRenameAttributes    `json:"fieldsRename,omitempty"`    // Fields rename processor attributes
	ID              string                     `json:"id"`                        // Processor identifier
	Matcher         *string                    `json:"matcher,omitempty"`         // [See our documentation](https://dt-url.net/bp234rv)
	SampleData      *string                    `json:"sampleData,omitempty"`      // Sample data
	SecurityContext *SecurityContextAttributes `json:"securityContext,omitempty"` // Security context processor attributes
	Type            ProcessorType              `json:"type"`                      // Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`
}

func (me *Processor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Required:    true,
		},
		"dql": {
			Type:        schema.TypeList,
			Description: "DQL processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(DqlAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"fields_add": {
			Type:        schema.TypeList,
			Description: "Fields add processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FieldsAddAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fields_remove": {
			Type:        schema.TypeList,
			Description: "Fields remove processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FieldsRemoveAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fields_rename": {
			Type:        schema.TypeList,
			Description: "Fields rename processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FieldsRenameAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Processor identifier",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "[See our documentation](https://dt-url.net/bp234rv)",
			Optional:    true, // precondition
		},
		"sample_data": {
			Type:        schema.TypeString,
			Description: "Sample data",
			Optional:    true, // nullable
		},
		"security_context": {
			Type:        schema.TypeList,
			Description: "Security context processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SecurityContextAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Processor type. Possible values: `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `securityContext`",
			Required:    true,
		},
	}
}

func (me *Processor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":      me.Description,
		"dql":              me.Dql,
		"enabled":          me.Enabled,
		"fields_add":       me.FieldsAdd,
		"fields_remove":    me.FieldsRemove,
		"fields_rename":    me.FieldsRename,
		"id":               me.ID,
		"matcher":          me.Matcher,
		"sample_data":      me.SampleData,
		"security_context": me.SecurityContext,
		"type":             me.Type,
	})
}

func (me *Processor) HandlePreconditions() error {
	if (me.Dql != nil) && (string(me.Type) != "dql") {
		return fmt.Errorf("'dql' must not be specified unless 'type' is set to 'dql'; got 'type'='%v'", me.Type)
	}
	if (me.Dql == nil) && (string(me.Type) == "dql") {
		return fmt.Errorf("'dql' must be specified when 'type' is set to 'dql'; got 'type'='%v'", me.Type)
	}
	if (me.FieldsAdd != nil) && (string(me.Type) != "fieldsAdd") {
		return fmt.Errorf("'fields_add' must not be specified unless 'type' is set to 'fieldsAdd'; got 'type'='%v'", me.Type)
	}
	if (me.FieldsAdd == nil) && (string(me.Type) == "fieldsAdd") {
		return fmt.Errorf("'fields_add' must be specified when 'type' is set to 'fieldsAdd'; got 'type'='%v'", me.Type)
	}
	if (me.FieldsRemove != nil) && (string(me.Type) != "fieldsRemove") {
		return fmt.Errorf("'fields_remove' must not be specified unless 'type' is set to 'fieldsRemove'; got 'type'='%v'", me.Type)
	}
	if (me.FieldsRemove == nil) && (string(me.Type) == "fieldsRemove") {
		return fmt.Errorf("'fields_remove' must be specified when 'type' is set to 'fieldsRemove'; got 'type'='%v'", me.Type)
	}
	if (me.FieldsRename != nil) && (string(me.Type) != "fieldsRename") {
		return fmt.Errorf("'fields_rename' must not be specified unless 'type' is set to 'fieldsRename'; got 'type'='%v'", me.Type)
	}
	if (me.FieldsRename == nil) && (string(me.Type) == "fieldsRename") {
		return fmt.Errorf("'fields_rename' must be specified when 'type' is set to 'fieldsRename'; got 'type'='%v'", me.Type)
	}
	if (me.Matcher != nil) && (string(me.Type) == "technology") {
		return fmt.Errorf("'matcher' must not be specified unless 'type' is not set to 'technology'; got 'type'='%v'", me.Type)
	}
	if (me.Matcher == nil) && (string(me.Type) != "technology") {
		return fmt.Errorf("'matcher' must be specified when 'type' is not set to 'technology'; got 'type'='%v'", me.Type)
	}
	if (me.SecurityContext != nil) && (string(me.Type) != "securityContext") {
		return fmt.Errorf("'security_context' must not be specified unless 'type' is set to 'securityContext'; got 'type'='%v'", me.Type)
	}
	if (me.SecurityContext == nil) && (string(me.Type) == "securityContext") {
		return fmt.Errorf("'security_context' must be specified when 'type' is set to 'securityContext'; got 'type'='%v'", me.Type)
	}
	return nil
}

func (me *Processor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description":      &me.Description,
		"dql":              &me.Dql,
		"enabled":          &me.Enabled,
		"fields_add":       &me.FieldsAdd,
		"fields_remove":    &me.FieldsRemove,
		"fields_rename":    &me.FieldsRename,
		"id":               &me.ID,
		"matcher":          &me.Matcher,
		"sample_data":      &me.SampleData,
		"security_context": &me.SecurityContext,
		"type":             &me.Type,
	})
}
