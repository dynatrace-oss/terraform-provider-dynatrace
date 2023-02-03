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

package advanceddetectionrule

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Delimiter struct {
	From      string `json:"from"`      // Delimit from
	RemoveIds bool   `json:"removeIds"` // (e.g. versions, hex, dates, and build numbers)
	To        string `json:"to"`        // Delimit to
}

func (me *Delimiter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from": {
			Type:        schema.TypeString,
			Description: "Delimit from",
			Optional:    true,
		},
		"remove_ids": {
			Type:        schema.TypeBool,
			Description: "(e.g. versions, hex, dates, and build numbers)",
			Required:    true,
		},
		"to": {
			Type:        schema.TypeString,
			Description: "Delimit to",
			Optional:    true,
		},
	}
}

func (me *Delimiter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"from":       me.From,
		"remove_ids": me.RemoveIds,
		"to":         me.To,
	})
}

func (me *Delimiter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"from":       &me.From,
		"remove_ids": &me.RemoveIds,
		"to":         &me.To,
	})
}
