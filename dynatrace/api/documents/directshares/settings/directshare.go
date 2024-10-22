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

package directshares

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SchemaVersion = 3

type Recipient struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Recipients []*Recipient

func (me *Recipient) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "Identifier of the recipient",
			ValidateDiagFunc: Validate(ValidateUUID, ValidateMaxLength(200)),
		},
		"type": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Type of the recipient. Possible values are `group' and `user'",
			ValidateDiagFunc: ValidateTypePossibleValues([]string{"group", "user"}),
		},
	}
}

func (me *Recipients) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"recipient": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1000,
			Description: "Recipient of the direct share",
			Elem:        &schema.Resource{Schema: new(Recipient).Schema()},
		},
	}
}

type DirectShare struct {
	ID            string     `json:"id"`
	DocumentId    string     `json:"document_id"`
	Access        string     `json:"access"`
	Recipients    Recipients `json:"recipients"`
	SchemaVersion int        `json:"schemaVersion,omitempty"`
}

func (me *DirectShare) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"document_id": {
			Type:             schema.TypeString,
			Description:      "Document ID",
			Required:         true,
			ValidateDiagFunc: Validate(ValidateUUID, ValidateMaxLength(200)),
		},
		"access": {
			Type:             schema.TypeString,
			Description:      "Access grants. Possible values are `read` and `read-write`",
			Optional:         true,
			Default:          "read",
			ValidateDiagFunc: ValidateTypePossibleValues([]string{"read", "read-write"}),
			ForceNew:         true,
		},
		"recipients": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Recipients of the direct share",
			Elem: &schema.Resource{
				Schema: new(Recipients).Schema(),
			},
		},
	}
}

func (me *DirectShare) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"document_id": me.DocumentId,
		"access":      me.Access,
		"recipients":  me.Recipients,
	})
}

func (me *DirectShare) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"document_id": &me.DocumentId,
		"access":      &me.Access,
		"recipients":  &me.Recipients,
	})
}

func (me *Recipient) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":   me.Id,
		"type": me.Type,
	})
}

func (me *Recipient) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":   &me.Id,
		"type": &me.Type,
	})
}

func (me Recipients) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("recipient", me)
}

func (me *Recipients) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("recipient", me)
}

func (me *DirectShare) MarshalJSON() ([]byte, error) {
	ds := struct {
		DocumentId    string     `json:"documentId"`
		Access        string     `json:"access"`
		Recipients    Recipients `json:"recipients"`
		SchemaVersion int        `json:"schemaVersion,omitempty"`
	}{
		DocumentId:    me.DocumentId,
		Access:        me.Access,
		Recipients:    me.Recipients,
		SchemaVersion: SchemaVersion,
	}
	return json.Marshal(ds)
}
