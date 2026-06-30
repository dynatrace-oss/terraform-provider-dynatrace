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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
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
	PrimaryGrailTag *bool         `json:"primaryGrailTag,omitempty"` // Uses the key of the annotation or label as field name directly
	Source          string        `json:"source"`                    // The source must follow the syntax of Kubernetes annotation/label keys as defined in the [Kubernetes documentation](https://dt-url.net/2c02sbn).\n\n  `source := (prefix/)?name`\n\n  `prefix := [a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`\n\n  `name := ([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]`\n\n  Additionally, the name can have at most 63 characters, and the overall length of the source must not exceed 75 characters.
	Target          *TargetOption `json:"target,omitempty"`          // Possible values: `dt.cost.costcenter`, `dt.cost.product`, `dt.security_context`
	Type            MetadataType  `json:"type"`                      // Metadata type. Possible values: `ANNOTATION`, `LABEL`
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"primary_grail_tag": {
			Type:        schema.TypeBool,
			Description: "Uses the key of the annotation or label as field name directly",
			Optional:    true, // nullable
		},
		"source": {
			Type:        schema.TypeString,
			Description: "The source must follow the syntax of Kubernetes annotation/label keys as defined in the [Kubernetes documentation](https://dt-url.net/2c02sbn).\n\n  `source := (prefix/)?name`\n\n  `prefix := [a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`\n\n  `name := ([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]`\n\n  Additionally, the name can have at most 63 characters, and the overall length of the source must not exceed 75 characters.",
			Required:    true,
		},
		"target": {
			Type:        schema.TypeString,
			Description: "Possible values: `dt.cost.costcenter`, `dt.cost.product`, `dt.security_context`",
			Optional:    true, // precondition
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Metadata type. Possible values: `ANNOTATION`, `LABEL`",
			Required:    true,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"primary_grail_tag": me.PrimaryGrailTag,
		"source":            me.Source,
		"target":            me.Target,
		"type":              me.Type,
	})
}

func (me *Rule) HandlePreconditions() error {
	if (me.Target != nil) && ((me.PrimaryGrailTag != nil) && (me.PrimaryGrailTag == nil || *me.PrimaryGrailTag)) {
		return fmt.Errorf("'target' must not be specified unless ('primary_grail_tag' is not set or 'primary_grail_tag' is set to 'false'); got 'primary_grail_tag'='%v'", opt.ValOrNil(me.PrimaryGrailTag))
	}
	if (me.Target == nil) && ((me.PrimaryGrailTag == nil) || (me.PrimaryGrailTag != nil && !*me.PrimaryGrailTag)) {
		return fmt.Errorf("'target' must be specified when ('primary_grail_tag' is not set or 'primary_grail_tag' is set to 'false'); got 'primary_grail_tag'='%v'", opt.ValOrNil(me.PrimaryGrailTag))
	}
	return nil
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"primary_grail_tag": &me.PrimaryGrailTag,
		"source":            &me.Source,
		"target":            &me.Target,
		"type":              &me.Type,
	})
}
