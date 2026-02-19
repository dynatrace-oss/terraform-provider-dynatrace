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

package sitereliabilityguardian

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ObjectiveLinks []*ObjectiveLink

func (me *ObjectiveLinks) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"link": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ObjectiveLink).Schema()},
		},
	}
}

func (me ObjectiveLinks) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("link", me)
}

func (me *ObjectiveLinks) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("link", me)
}

type ObjectiveLink struct {
	Label *string `json:"label,omitempty"` // Short description for the link.
	Url   string  `json:"url"`             // HTTPS link associated with this objective.
}

func (me *ObjectiveLink) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Description: "Short description for the link.",
			Optional:    true, // nullable
		},
		"url": {
			Type:        schema.TypeString,
			Description: "HTTPS link associated with this objective.",
			Required:    true,
		},
	}
}

func (me *ObjectiveLink) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"label": me.Label,
		"url":   me.Url,
	})
}

func (me *ObjectiveLink) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"label": &me.Label,
		"url":   &me.Url,
	})
}
