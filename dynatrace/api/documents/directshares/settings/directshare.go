/**
* @license
* Copyright 2026 Dynatrace LLC
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DirectShare struct {
	ID         string
	DocumentId string
	Access     string
	Recipients Recipients
}

func (me *DirectShare) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"document_id": {
			Type:             schema.TypeString,
			Description:      "Document ID",
			Required:         true,
			ValidateDiagFunc: Validate(ValidateMaxLength(200)),
			ForceNew:         true,
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
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Recipients of the direct share",
			Elem: &schema.Resource{
				Schema: new(Recipients).Schema(),
			},
		},
	}
}

func (me *DirectShare) Name() string {
	return "direct_shares"
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

type Recipients []*Recipient

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

func (me Recipients) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("recipient", me)
}

func (me *Recipients) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("recipient", me)
}

type Recipient struct {
	ID   string
	Type string
}

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

func (me *Recipient) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":   me.ID,
		"type": me.Type,
	})
}

func (me *Recipient) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":   &me.ID,
		"type": &me.Type,
	})
}
