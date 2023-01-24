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

package request

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Cookies contains the list of cookies to be created for the monitor. Every cookie must be unique within the list. However, you can use the same cookie again in other event
type Cookies []*Cookie

func (me *Cookies) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cookie": {
			Type:        schema.TypeList,
			Description: "A request cookie",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Cookie).Schema()},
		},
	}
}

func (me Cookies) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeSlice("cookie", me); err != nil {
		return err
	}
	return nil
}

func (me *Cookies) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("cookie", me)
}

// Cookie a request cookie
type Cookie struct {
	Name   string  `json:"name"`           // The name of the cookie. The following cookie names are now allowed: `dtCookie`, `dtLatC`, `dtPC`, `rxVisitor`, `rxlatency`, `rxpc`, `rxsession` and `rxvt`
	Value  string  `json:"value"`          // The value of the cookie. The following symbols are not allowed: `;`, `,`, `\` and `"`.
	Domain string  `json:"domain"`         // The domain of the cookie
	Path   *string `json:"path,omitempty"` // The path of the cookie
}

func (me *Cookie) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the cookie. The following cookie names are now allowed: `dtCookie`, `dtLatC`, `dtPC`, `rxVisitor`, `rxlatency`, `rxpc`, `rxsession` and `rxvt`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the cookie. The following symbols are not allowed: `;`, `,`, `\\` and `\"`.",
			Required:    true,
		},
		"domain": {
			Type:        schema.TypeString,
			Description: "The domain of the cookie.",
			Required:    true,
		},
		"path": {
			Type:        schema.TypeString,
			Description: "The path of the cookie.",
			Optional:    true,
		},
	}
}

func (me *Cookie) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}
	if err := properties.Encode("domain", me.Domain); err != nil {
		return err
	}
	if err := properties.Encode("path", me.Path); err != nil {
		return err
	}
	return nil
}

func (me *Cookie) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("name", &me.Name); err != nil {
		return err
	}
	if err := decoder.Decode("value", &me.Value); err != nil {
		return err
	}
	if err := decoder.Decode("domain", &me.Domain); err != nil {
		return err
	}
	if err := decoder.Decode("path", &me.Path); err != nil {
		return err
	}
	return nil
}
