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

package documents

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SchemaVersion = 3

type Document struct {
	Name          string `json:"name" maxlength:"200"`                         // The name of the document
	Content       string `json:"content,omitempty"`                            // The content of the document
	Type          string `json:"type"`                                         // Type of the document
	Actor         string `json:"actor,omitempty" maxlength:"36" format:"uuid"` // The user context the executions of the document will happen with
	Owner         string `json:"owner,omitempty" format:"uuid"`                // The ID of the owner of this document
	Version       int    `json:"version,omitempty"`                            // The version of the document
	SchemaVersion int    `json:"schemaVersion,omitempty"`                      //
}

func (me *Document) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Description:      "The name/name of the document",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(200),
		},
		"type": {
			Type:             schema.TypeString,
			Description:      "Type of the document",
			Required:         true,
			ValidateDiagFunc: ValidateTypePossibleValues([]string{"dashboard", "notebook"}),
		},
		"actor": {
			Type:             schema.TypeString,
			Description:      "The user context the executions of the document will happen with",
			Optional:         true,
			ValidateDiagFunc: Validate(ValidateUUID, ValidateMaxLength(36)),
		},
		"owner": {
			Type:             schema.TypeString,
			Description:      "The ID of the owner of this document",
			Optional:         true,
			ValidateDiagFunc: ValidateUUID,
		},
		"content": {
			Type:        schema.TypeString,
			Description: "Document content",
			Required:    true,
			//ValidateDiagFunc: ValidateUUID, @todo: implement a custom validation function
		},
		"version": {
			Type:        schema.TypeInt,
			Description: "The version of the document",
			Computed:    true,
		},
	}
}

func (me *Document) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"content": me.Content,
		"type":    me.Type,
		"actor":   me.Actor,
		"owner":   me.Owner,
		"version": me.Version,
	})
}

func (me *Document) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"content": &me.Content,
		"type":    &me.Type,
		"actor":   &me.Actor,
		"owner":   &me.Owner,
		"version": &me.Version,
	})
}

func (me *Document) MarshalJSON() ([]byte, error) {
	wf := struct {
		Name          string `json:"name"`
		Content       string `json:"content,omitempty"`
		Type          string `json:"type"`
		Actor         string `json:"actor,omitempty"`
		Owner         string `json:"owner,omitempty"`
		Version       int    `json:"version,omitempty"`
		SchemaVersion int    `json:"schemaVersion,omitempty"`
	}{
		SchemaVersion: SchemaVersion, // adding the Schema Version is the purpose of this custome `MarshalJSON` function
		Name:          me.Name,
		Content:       me.Content,
		Type:          me.Type,
		Actor:         me.Actor,
		Owner:         me.Owner,
		Version:       me.Version,
	}
	return json.Marshal(wf)
}
