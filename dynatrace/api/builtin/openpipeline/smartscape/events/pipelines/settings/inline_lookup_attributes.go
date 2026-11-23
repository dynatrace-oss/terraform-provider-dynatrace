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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type InlineLookupAttributes struct {
	DefaultValue      *string `json:"defaultValue,omitempty"` // The value to write to the destination field when no lookup key matches. If absent, the destination field is left unchanged when no key matches.
	DestinationField  string  `json:"destinationField"`       // The field key to write the matched lookup value to.
	InlineLookupTable string  `json:"inlineLookupTable"`      // The key-value pairs of the inline lookup table, encoded as a compact JSON string: [[[\"key1\",\"key2\"],\"value1\"],[[\"key3\"],\"value2\"]].
	SourceField       string  `json:"sourceField"`            // The field key whose value is looked up in the lookup table.
}

func (me *InlineLookupAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_value": {
			Type:        schema.TypeString,
			Description: "The value to write to the destination field when no lookup key matches. If absent, the destination field is left unchanged when no key matches.",
			Optional:    true, // nullable
		},
		"destination_field": {
			Type:        schema.TypeString,
			Description: "The field key to write the matched lookup value to.",
			Required:    true,
		},
		"inline_lookup_table": {
			Type:        schema.TypeString,
			Description: "The key-value pairs of the inline lookup table, encoded as a compact JSON string: [[[\"key1\",\"key2\"],\"value1\"],[[\"key3\"],\"value2\"]].",
			Required:    true,
		},
		"source_field": {
			Type:        schema.TypeString,
			Description: "The field key whose value is looked up in the lookup table.",
			Required:    true,
		},
	}
}

func (me *InlineLookupAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"default_value":       me.DefaultValue,
		"destination_field":   me.DestinationField,
		"inline_lookup_table": me.InlineLookupTable,
		"source_field":        me.SourceField,
	})
}

func (me *InlineLookupAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"default_value":       &me.DefaultValue,
		"destination_field":   &me.DestinationField,
		"inline_lookup_table": &me.InlineLookupTable,
		"source_field":        &me.SourceField,
	})
}
