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

package segments

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Includes []*Items

func (me *Includes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"items": {
			Type:        schema.TypeSet,
			Description: "TODO: No documentation available",
			MaxItems:    20,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Items).Schema()},
		},
	}

}

func (me Includes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("items", me)
}

func (me *Includes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("items", me)
}

type Items struct {
	Filter       string        `json:"filter"`
	DataObject   string        `json:"dataObject"`
	ApplyTo      []string      `json:"applyTo,omitempty"`
	Relationship *Relationship `json:"relationship,omitempty"`
}

func (me *Items) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:             schema.TypeString,
			Description:      "Data will be filtered by this value",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(10000),
		},
		"data_object": {
			Type:             schema.TypeString,
			Description:      "The data object that the filter will be applied to. Use '_all_data_object' to apply it to all dataObjects",
			Required:         true,
			ValidateDiagFunc: ValidateMaxLength(500),
		},
		"apply_to": {
			Type:        schema.TypeSet,
			Description: "[Experimental] The tables that the entity-filter will be applied to`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"relationship": {
			Type:        schema.TypeList,
			Description: "[Experimental] The relationship of an include which has to be be specified when the data object is an entity view",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Relationship).Schema()},
		},
	}
}

func (me *Items) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"filter":       me.Filter,
		"data_object":  me.DataObject,
		"apply_to":     me.ApplyTo,
		"relationship": me.Relationship,
	})
}

func (me *Items) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"filter":       &me.Filter,
		"data_object":  &me.DataObject,
		"apply_to":     &me.ApplyTo,
		"relationship": &me.Relationship,
	})
}
