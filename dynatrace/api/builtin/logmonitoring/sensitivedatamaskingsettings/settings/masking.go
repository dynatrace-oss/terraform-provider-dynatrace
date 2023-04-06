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

package sensitivedatamasking

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Masking struct {
	Expression  string      `json:"expression"` // Maximum one capture group is allowed. If none was given, the whole expression will be treated as a capture group.
	Replacement string      `json:"replacement"`
	Type        MaskingType `json:"type"` // Possible Values: `SHA1`, `STRING`
}

func (me *Masking) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"expression": {
			Type:        schema.TypeString,
			Description: "Maximum one capture group is allowed. If none was given, the whole expression will be treated as a capture group.",
			Required:    true,
		},
		"replacement": {
			Type:        schema.TypeString,
			Description: "The string to replace the masked expression with. Use `SHA1` if it should to automatically generated every time.",
			Required:    true,
		},
	}
}

func (me *Masking) MarshalHCL(properties hcl.Properties) error {
	if me.Type == "SHA1" {
		me.Replacement = "SHA1"
	}
	return properties.EncodeAll(map[string]any{
		"expression":  me.Expression,
		"replacement": me.Replacement,
	})
}

func (me *Masking) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"expression":  &me.Expression,
		"replacement": &me.Replacement,
	}); err != nil {
		return err
	}
	if me.Replacement == "SHA1" {
		me.Replacement = ""
		me.Type = "SHA1"
	} else {
		me.Type = "STRING"
	}
	return nil
}
