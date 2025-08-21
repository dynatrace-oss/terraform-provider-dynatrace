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

package failuredetectionrulesets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SingleExceptions []*SingleException

func (me *SingleExceptions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ignored_exception": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SingleException).Schema()},
		},
	}
}

func (me SingleExceptions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("ignored_exception", me)
}

func (me *SingleExceptions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("ignored_exception", me)
}

type SingleException struct {
	Enabled bool    `json:"enabled"`           // This setting is enabled (`true`) or disabled (`false`)
	Message *string `json:"message,omitempty"` // Evaluated attribute: `span.events[][exception.message]`
	Type    *string `json:"type,omitempty"`    // Evaluated attribute: `span.events[][exception.type]`
}

func (me *SingleException) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"message": {
			Type:        schema.TypeString,
			Description: "Evaluated attribute: `span.events[][exception.message]`",
			Optional:    true, // nullable
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Evaluated attribute: `span.events[][exception.type]`",
			Optional:    true, // nullable
		},
	}
}

func (me *SingleException) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"message": me.Message,
		"type":    me.Type,
	})
}

func (me *SingleException) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled": &me.Enabled,
		"message": &me.Message,
		"type":    &me.Type,
	})
}
