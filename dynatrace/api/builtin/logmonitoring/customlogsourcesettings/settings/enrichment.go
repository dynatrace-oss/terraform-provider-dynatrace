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

type Enrichments []*Enrichment

func (me *Enrichments) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enrichment": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Enrichment).Schema()},
		},
	}
}

func (me Enrichments) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("enrichment", me)
}

func (me *Enrichments) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("enrichment", me)
}

type Enrichment struct {
	Key   *string               `json:"key,omitempty"`
	Type  WildcardExpansionType `json:"type"` // Possible Values: `Attribute`
	Value *string               `json:"value,omitempty"`
}

func (me *Enrichment) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Attribute`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
	}
}

func (me *Enrichment) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":   me.Key,
		"type":  me.Type,
		"value": me.Value,
	})
}

func (me *Enrichment) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":   &me.Key,
		"type":  &me.Type,
		"value": &me.Value,
	})
}
