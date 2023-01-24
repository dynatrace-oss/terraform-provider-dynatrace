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

package alerting

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomEventFilter struct {
	Title       *TextFilter `json:"titleFilter,omitempty"`       // Title filter
	Description *TextFilter `json:"descriptionFilter,omitempty"` // Description filter
}

func (me *CustomEventFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeList,
			Description: "Configuration of a matching filter",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(TextFilter).Schema()},
		},
		"title": {
			Type:        schema.TypeList,
			Description: "Configuration of a matching filter",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(TextFilter).Schema()},
		},
	}
}

func (me *CustomEventFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("description", me.Description); err != nil {
		return err
	}
	if err := properties.Encode("title", me.Title); err != nil {
		return err
	}
	return nil
}

func (me *CustomEventFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("description.#"); ok {
		me.Description = new(TextFilter)
		if err := me.Description.UnmarshalHCL(hcl.NewDecoder(decoder, "description", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("title.#"); ok {
		me.Title = new(TextFilter)
		if err := me.Title.UnmarshalHCL(hcl.NewDecoder(decoder, "title", 0)); err != nil {
			return err
		}
	}
	return nil
}
