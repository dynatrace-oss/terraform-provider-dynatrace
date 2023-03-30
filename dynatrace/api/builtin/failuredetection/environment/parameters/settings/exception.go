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

package parameters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Exceptions []*Exception

func (me *Exceptions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_handled_exception": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Exception).Schema()},
		},
	}
}

func (me Exceptions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("custom_handled_exception", me)
}

func (me *Exceptions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("custom_handled_exception", me)
}

type Exception struct {
	ClassPattern   *string `json:"classPattern,omitempty"`   // The pattern will match if it is contained within the actual class name.
	MessagePattern *string `json:"messagePattern,omitempty"` // Optionally, define an exception message pattern. The pattern will match if the actual exception message contains the pattern.
}

func (me *Exception) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"class_pattern": {
			Type:        schema.TypeString,
			Description: "The pattern will match if it is contained within the actual class name.",
			Optional:    true, // nullable
		},
		"message_pattern": {
			Type:        schema.TypeString,
			Description: "Optionally, define an exception message pattern. The pattern will match if the actual exception message contains the pattern.",
			Optional:    true, // nullable
		},
	}
}

func (me *Exception) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"class_pattern":   me.ClassPattern,
		"message_pattern": me.MessagePattern,
	})
}

func (me *Exception) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"class_pattern":   &me.ClassPattern,
		"message_pattern": &me.MessagePattern,
	})
}
