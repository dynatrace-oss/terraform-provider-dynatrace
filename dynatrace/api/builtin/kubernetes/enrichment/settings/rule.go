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

package enrichment

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rules []*Rule

func (me *Rules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Rule).Schema()},
		},
	}
}

func (me Rules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *Rules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

type Rule struct {
	Enabled         *bool        `json:"enabled,omitempty"`         // This setting is enabled (`true`) or disabled (`false`)
	Source          string       `json:"source"`                    // The source must follow the syntax of Kubernetes annotation/label keys as defined in the [Kubernetes documentation](https://dt-url.net/2c02sbn).\n\n`source := (prefix/)?name`\n\n`prefix := [a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`\n\n`name := ([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]`\n\nAdditionally, the name can have at most 63 characters, and the overall length of the source must not exceed 75 characters.
	Target          TargetOption `json:"target,omitempty"`          // Required when `primary_grail_tag` is omitted or `false`. Possible Values: `dt.cost.costcenter``, `dt.cost.product``, `dt.security_context
	Type            MetadataType `json:"type"`                      // Possible Values: `ANNOTATION`, `LABEL`
	PrimaryGrailTag *bool        `json:"primaryGrailTag,omitempty"` // Uses the key of the annotation or label as field name
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Optional:    true,
			Deprecated:  "Attribute no longer exists in the schema.",
		},
		"source": {
			Type:        schema.TypeString,
			Description: "The source must follow the syntax of Kubernetes annotation/label keys as defined in the [Kubernetes documentation](https://dt-url.net/2c02sbn).\n\n`source := (prefix/)?name`\n\n`prefix := [a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`\n\n`name := ([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]`\n\nAdditionally, the name can have at most 63 characters, and the overall length of the source must not exceed 75 characters.",
			Required:    true,
		},
		"target": {
			Type:        schema.TypeString,
			Description: "Required when `primary_grail_tag` is omitted or `false`. Possible Values: `dt.cost.costcenter``, `dt.cost.product``, `dt.security_context",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ANNOTATION`, `LABEL`",
			Required:    true,
		},
		"primary_grail_tag": {
			Type:        schema.TypeBool,
			Description: "Uses the key of the annotation or label as field name",
			Optional:    true, // nullable
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":           me.Enabled,
		"source":            me.Source,
		"target":            me.Target,
		"type":              me.Type,
		"primary_grail_tag": me.PrimaryGrailTag,
	})
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":           &me.Enabled,
		"source":            &me.Source,
		"target":            &me.Target,
		"type":              &me.Type,
		"primary_grail_tag": &me.PrimaryGrailTag,
	})
}
