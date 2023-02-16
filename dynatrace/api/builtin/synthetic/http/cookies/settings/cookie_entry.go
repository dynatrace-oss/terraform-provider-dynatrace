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

package cookies

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CookieEntries []*CookieEntry

func (me *CookieEntries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cookie": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CookieEntry).Schema()},
		},
	}
}

func (me CookieEntries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("cookie", me)
}

func (me *CookieEntries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("cookie", me)
}

type CookieEntry struct {
	Domain string  `json:"domain"`         // Enclose placeholder values in brackets, for example \\{email\\}
	Name   string  `json:"name"`           // Enclose placeholder values in brackets, for example \\{email\\}
	Path   *string `json:"path,omitempty"` // Enclose placeholder values in brackets, for example \\{email\\}
	Value  string  `json:"value"`          // Enclose placeholder values in brackets, for example \\{email\\}
}

func (me *CookieEntry) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain": {
			Type:        schema.TypeString,
			Description: "Enclose placeholder values in brackets, for example \\{email\\}",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Enclose placeholder values in brackets, for example \\{email\\}",
			Required:    true,
		},
		"path": {
			Type:        schema.TypeString,
			Description: "Enclose placeholder values in brackets, for example \\{email\\}",
			Optional:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Enclose placeholder values in brackets, for example \\{email\\}",
			Required:    true,
		},
	}
}

func (me *CookieEntry) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"domain": me.Domain,
		"name":   me.Name,
		"path":   me.Path,
		"value":  me.Value,
	})
}

func (me *CookieEntry) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"domain": &me.Domain,
		"name":   &me.Name,
		"path":   &me.Path,
		"value":  &me.Value,
	})
}
