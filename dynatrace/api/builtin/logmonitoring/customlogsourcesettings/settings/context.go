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

package customlogsourcesettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Contexts []*Context

func (me *Contexts) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"context": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Context).Schema()},
		},
	}
}

func (me Contexts) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("context", me)
}

func (me *Contexts) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("context", me)
}

type Context struct {
	Attribute ContextType `json:"attribute"` // Possible Values: `Dt_entity_process_group`
	Values    []string    `json:"values"`
}

func (me *Context) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Dt_entity_process_group`",
			Required:    true,
		},
		"values": {
			Type:        schema.TypeSet,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Context) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attribute": me.Attribute,
		"values":    me.Values,
	})
}

func (me *Context) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attribute": &me.Attribute,
		"values":    &me.Values,
	})
}
