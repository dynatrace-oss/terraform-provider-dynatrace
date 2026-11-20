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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JsonExtractionAttributes struct {
	MaxNestingDepth    int `json:"maxNestingDepth"`    // Maximum depth of nested JSON that is unwrapped.
	MaxObjectsToUnwrap int `json:"maxObjectsToUnwrap"` // Maximum number of JSON objects that are unwrapped.
}

func (me *JsonExtractionAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"max_nesting_depth": {
			Type:        schema.TypeInt,
			Description: "Maximum depth of nested JSON that is unwrapped.",
			Required:    true,
		},
		"max_objects_to_unwrap": {
			Type:        schema.TypeInt,
			Description: "Maximum number of JSON objects that are unwrapped.",
			Required:    true,
		},
	}
}

func (me *JsonExtractionAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"max_nesting_depth":     me.MaxNestingDepth,
		"max_objects_to_unwrap": me.MaxObjectsToUnwrap,
	})
}

func (me *JsonExtractionAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"max_nesting_depth":     &me.MaxNestingDepth,
		"max_objects_to_unwrap": &me.MaxObjectsToUnwrap,
	})
}
