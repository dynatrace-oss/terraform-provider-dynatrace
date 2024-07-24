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

package processavailability

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HostMetadataCondition struct {
	KeyMustExist      bool   `json:"keyMustExist"`      // When enabled, the condition requires a metadata key to exist and match the constraints; when disabled, the key is optional but must still match the constrains if it is present.
	MetadataCondition string `json:"metadataCondition"` // This string has to match a required format.\n\n- `$contains(production)` – Matches if `production` appears anywhere in the host metadata value.\n- `$eq(production)` – Matches if `production` matches the host metadata value exactly.\n- `$prefix(production)` – Matches if `production` matches the prefix of the host metadata value.\n- `$suffix(production)` – Matches if `production` matches the suffix of the host metadata value.\n\nAvailable logic operations:\n- `$not($eq(production))` – Matches if the host metadata value is different from `production`.\n- `$and($prefix(production),$suffix(main))` – Matches if host metadata value starts with `production` and ends with `main`.\n- `$or($prefix(production),$suffix(main))` – Matches if host metadata value starts with `production` or ends with `main`.\n\nBrackets **(** and **)** that are part of the matched property **must be escaped with a tilde (~)**
	MetadataKey       string `json:"metadataKey"`       // Key
}

func (me *HostMetadataCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key_must_exist": {
			Type:        schema.TypeBool,
			Description: "When enabled, the condition requires a metadata key to exist and match the constraints; when disabled, the key is optional but must still match the constrains if it is present.",
			Optional:    true,
			Default:     true,
		},
		"metadata_condition": {
			Type:        schema.TypeString,
			Description: "This string has to match a required format.\n\n- `$contains(production)` – Matches if `production` appears anywhere in the host metadata value.\n- `$eq(production)` – Matches if `production` matches the host metadata value exactly.\n- `$prefix(production)` – Matches if `production` matches the prefix of the host metadata value.\n- `$suffix(production)` – Matches if `production` matches the suffix of the host metadata value.\n\nAvailable logic operations:\n- `$not($eq(production))` – Matches if the host metadata value is different from `production`.\n- `$and($prefix(production),$suffix(main))` – Matches if host metadata value starts with `production` and ends with `main`.\n- `$or($prefix(production),$suffix(main))` – Matches if host metadata value starts with `production` or ends with `main`.\n\nBrackets **(** and **)** that are part of the matched property **must be escaped with a tilde (~)**",
			Required:    true,
		},
		"metadata_key": {
			Type:        schema.TypeString,
			Description: "Key",
			Required:    true,
		},
	}
}

func (me *HostMetadataCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key_must_exist":     me.KeyMustExist,
		"metadata_condition": me.MetadataCondition,
		"metadata_key":       me.MetadataKey,
	})
}

func (me *HostMetadataCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key_must_exist":     &me.KeyMustExist,
		"metadata_condition": &me.MetadataCondition,
		"metadata_key":       &me.MetadataKey,
	})
}
