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

package rpcbasedsampling

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type QueryParameter struct {
	Name             string  `json:"name"`                       // Query parameter name
	Value            *string `json:"value,omitempty"`            // The value must be equal for a match.
	ValueIsUndefined *bool   `json:"valueIsUndefined,omitempty"` // If enabled, the value is treated as undefined
}

func (me *QueryParameter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Query parameter name",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value must be equal for a match.",
			Optional:    true, // nullable
		},
		"value_is_undefined": {
			Type:        schema.TypeBool,
			Description: "If enabled, the value is treated as undefined",
			Optional:    true, // precondition
		},
	}
}

func (me *QueryParameter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":               me.Name,
		"value":              me.Value,
		"value_is_undefined": me.ValueIsUndefined,
	})
}

func (me *QueryParameter) HandlePreconditions() error {
	// ---- ValueIsUndefined *bool -> {"property":"value","type":"NULL"}
	return nil
}

func (me *QueryParameter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":               &me.Name,
		"value":              &me.Value,
		"value_is_undefined": &me.ValueIsUndefined,
	})
}
