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

package teams

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Links []*Link

func (me *Links) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"link": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Link).Schema()},
		},
	}
}

func (me Links) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("link", me)
}

func (me *Links) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("link", me)
}

type Link struct {
	LinkType LinkType `json:"linkType"` // Possible Values: `DASHBOARD`, `DOCUMENTATION`, `HEALTH_APP`, `REPOSITORY`, `RUNBOOK`, `URL`, `WIKI`
	Url      string   `json:"url"`
}

func (me *Link) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"link_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DASHBOARD`, `DOCUMENTATION`, `HEALTH_APP`, `REPOSITORY`, `RUNBOOK`, `URL`, `WIKI`",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *Link) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"link_type": me.LinkType,
		"url":       me.Url,
	})
}

func (me *Link) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"link_type": &me.LinkType,
		"url":       &me.Url,
	})
}
