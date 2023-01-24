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

package naming

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Placeholders []*Placeholder

func (me *Placeholders) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"placeholder": {
			Type:        schema.TypeList,
			Description: "User action placeholders",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Placeholder).Schema()},
		},
	}
}

func (me Placeholders) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("placeholder", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *Placeholders) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("placeholder", me)
}

// Placeholder The placeholder settings
type Placeholder struct {
	Name                        string          `json:"name"`                        // Placeholder name. Valid length needs to be between 1 and 50 characters.
	Input                       Input           `json:"input"`                       // The input for the place holder. Possible values are `ELEMENT_IDENTIFIER`, `INPUT_TYPE`, `METADATA`, `PAGE_TITLE`, `PAGE_URL`, `SOURCE_URL`, `TOP_XHR_URL` and `XHR_URL`.
	ProcessingPart              ProcessingPart  `json:"processingPart"`              // The part to process. Possible values are `ALL`, `ANCHOR` and `PATH`.
	ProcessingSteps             ProcessingSteps `json:"processingSteps,omitempty"`   // The processing step settings
	MetaDataID                  *int32          `json:"metadataId,omitempty"`        // The ID of the metadata
	UseGuessedElementIdentifier bool            `json:"useGuessedElementIdentifier"` // Use the element identifier that was selected by Dynatrace
}

func (me *Placeholder) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Placeholder name. Valid length needs to be between 1 and 50 characters",
			Required:    true,
		},
		"input": {
			Type:        schema.TypeString,
			Description: "The input for the place holder. Possible values are `ELEMENT_IDENTIFIER`, `INPUT_TYPE`, `METADATA`, `PAGE_TITLE`, `PAGE_URL`, `SOURCE_URL`, `TOP_XHR_URL` and `XHR_URL`",
			Required:    true,
		},
		"processing_part": {
			Type:        schema.TypeString,
			Description: "The part to process. Possible values are `ALL`, `ANCHOR` and `PATH`",
			Required:    true,
		},
		"processing_steps": {
			Type:        schema.TypeList,
			Description: "The processing step settings",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ProcessingSteps).Schema()},
		},
		"metadata_id": {
			Type:        schema.TypeInt,
			Description: "The ID of the metadata",
			Optional:    true,
		},
		"use_guessed_element_identifier": {
			Type:        schema.TypeBool,
			Description: "Use the element identifier that was selected by Dynatrace",
			Optional:    true,
		},
	}
}

func (me *Placeholder) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":                           me.Name,
		"input":                          me.Input,
		"processing_part":                me.ProcessingPart,
		"processing_steps":               me.ProcessingSteps,
		"metadata_id":                    me.MetaDataID,
		"use_guessed_element_identifier": me.UseGuessedElementIdentifier,
	})
}

func (me *Placeholder) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":                           &me.Name,
		"input":                          &me.Input,
		"processing_part":                &me.ProcessingPart,
		"processing_steps":               &me.ProcessingSteps,
		"metadata_id":                    &me.MetaDataID,
		"use_guessed_element_identifier": &me.UseGuessedElementIdentifier,
	})
}
