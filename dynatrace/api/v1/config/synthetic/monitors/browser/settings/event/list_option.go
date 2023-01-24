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

type ListOptions []*ListOption

func (me *ListOptions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"option": {
			Type:        schema.TypeList,
			Description: "The option to be selected",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ListOption).Schema()},
		},
	}
}

func (me ListOptions) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeSlice("option", me); err != nil {
		return err
	}

	return nil
}

func (me *ListOptions) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("option", &me); err != nil {
		return err
	}
	return nil
}

type ListOption struct {
	Index int    `json:"index"` // The index of the option to be selected
	Value string `json:"value"` // The value of the option to be selected
}

func (me *ListOption) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"index": {
			Type:        schema.TypeInt,
			Description: "The index of the option to be selected",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the option to be selected",
			Required:    true,
		},
	}
}

func (me *ListOption) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("index", me.Index); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}

	return nil
}

func (me *ListOption) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("index", &me.Index); err != nil {
		return err
	}
	if err := decoder.Decode("value", &me.Value); err != nil {
		return err
	}
	return nil
}
