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
	Name          string `json:"name" maxlength:"200"`
	Content       string `json:"content,omitempty"`
	IsPrivate     bool   `json:"isPrivate,omitempty"`
	Type          string `json:"type"`
	Owner         string `json:"owner,omitempty" format:"uuid"`
	Version       int    `json:"version,omitempty"`
	SchemaVersion int    `json:"schemaVersion,omitempty"`
}

func (me *Document) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Description:      "The name/name of the document",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(200),
		},
		"private": {
			Type:        schema.TypeBool,
			Description: "Specifies whether the document is private or readable by everybody",
			Optional:    true,
			Default:     false,
		},
		"type": {
			Type:             schema.TypeString,
			Description:      "Type of the document. Possible Values are `dashboard`, `launchpad` and `notebook`",
			Required:         true,
			ValidateDiagFunc: ValidateTypePossibleValues([]string{"dashboard", "notebook", "launchpad"}),
		},
		"owner": {
			Type:        schema.TypeString,
			Description: "The ID of the owner of this document",
			Computed:    true,
		},
		"content": {
			Type:        schema.TypeString,
			Description: "Document content as JSON",
			Required:    true,
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
		"private": me.IsPrivate,
		"type":    me.Type,
		"owner":   me.Owner,
		"version": me.Version,
	})
}

func (me *Document) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"content": &me.Content,
		"private": &me.IsPrivate,
		"type":    &me.Type,
		"owner":   &me.Owner,
		"version": &me.Version,
	})
}

func (me *Document) MarshalJSON() ([]byte, error) {
	d := struct {
		Name          string `json:"name"`
		Content       string `json:"content,omitempty"`
		Private       bool   `json:"isPrivate,omitempty"`
		Type          string `json:"type"`
		Actor         string `json:"actor,omitempty"`
		Owner         string `json:"owner,omitempty"`
		Version       int    `json:"version,omitempty"`
		SchemaVersion int    `json:"schemaVersion,omitempty"`
	}{
		SchemaVersion: SchemaVersion,
		Name:          me.Name,
		Private:       me.IsPrivate,
		Content:       me.Content,
		Type:          me.Type,
		Owner:         me.Owner,
		Version:       me.Version,
	}
	return json.Marshal(d)
}
