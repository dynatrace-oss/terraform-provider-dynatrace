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

package event

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Validations []*Validation

func (me *Validations) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"validation": {
			Type:        schema.TypeList,
			Description: "The element to wait for. Required for the `validation` type, not applicable otherwise.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Validation).Schema()},
		},
	}
}

func (me Validations) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeSlice("validation", me); err != nil {
		return err
	}
	return nil
}

func (me *Validations) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("validation", me)
}

type Validation struct {
	Type        ValidationType `json:"type"`              // The goal of the validation. `content_match` (check page for the specific content. Not allowed for validation inside of wait condition), `element_match` (check page for the specific element)
	Match       string         `json:"match"`             // The content to look for on the page.\nRegular expressions are allowed. In that case set `isRegex` as `true`. Required for `content_match`, optional for `element_match`.
	IsRegex     bool           `json:"isRegex,omitempty"` // Defines whether `match` is plain text (`false`) of a regular expression (`true`)
	FailIfFound bool           `json:"failIfFound"`       // The condition of the validation. `false` means the validation succeeds if the specified content/element is found. `true` means the validation fails if the specified content/element is found
	Target      *Target        `json:"target,omitempty"`  // The elemnt to look for on the page
}

func (me *Validation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The goal of the validation. `content_match` (check page for the specific content. Not allowed for validation inside of wait condition), `element_match` (check page for the specific element).",
			Required:    true,
		},
		"match": {
			Type:        schema.TypeString,
			Description: "The content to look for on the page.\nRegular expressions are allowed. In that case set `isRegex` as `true`. Required for `content_match`, optional for `element_match`.",
			Optional:    true,
		},
		"regex": {
			Type:        schema.TypeBool,
			Description: "Defines whether `match` is plain text (`false`) or a regular expression (`true`)",
			Optional:    true,
		},
		"fail_if_found": {
			Type:        schema.TypeBool,
			Description: "The condition of the validation. `false` means the validation succeeds if the specified content/element is found. `true` means the validation fails if the specified content/element is found",
			Optional:    true,
		},
		"target": {
			Type:        schema.TypeList,
			Description: "The elemnt to look for on the page",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Target).Schema()},
		},
	}
}

func (me *Validation) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("match", me.Match); err != nil {
		return err
	}
	if err := properties.Encode("regex", me.IsRegex); err != nil {
		return err
	}
	if err := properties.Encode("fail_if_found", me.FailIfFound); err != nil {
		return err
	}
	if err := properties.Encode("target", me.Target); err != nil {
		return err
	}
	return nil
}

func (me *Validation) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	if err := decoder.Decode("match", &me.Match); err != nil {
		return err
	}
	if err := decoder.Decode("regex", &me.IsRegex); err != nil {
		return err
	}
	if err := decoder.Decode("fail_if_found", &me.FailIfFound); err != nil {
		return err
	}
	if err := decoder.Decode("target", &me.Target); err != nil {
		return err
	}
	return nil
}

// ValidationType The goal of the validation. `content_match` (check page for the specific content. Not allowed for validation inside of wait condition), `element_match` (check page for the specific element)
type ValidationType string

// ValidationTypes offers the known enum values
var ValidationTypes = struct {
	ContentMatch ValidationType
	ElementMatch ValidationType
}{
	`content_match`,
	`element_match`,
}
