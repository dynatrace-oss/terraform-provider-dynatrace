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

type AdditionalInformations []*AdditionalInformation

func (me *AdditionalInformations) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"additional_information": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(AdditionalInformation).Schema()},
		},
	}
}

func (me AdditionalInformations) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("additional_information", me)
}

func (me *AdditionalInformations) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("additional_information", me)
}

type AdditionalInformation struct {
	Key   string  `json:"key"` // Name
	Url   *string `json:"url,omitempty"`
	Value string  `json:"value"`
}

func (me *AdditionalInformation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"value": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *AdditionalInformation) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":   me.Key,
		"url":   me.Url,
		"value": me.Value,
	})
}

func (me *AdditionalInformation) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"key":   &me.Key,
		"url":   &me.Url,
		"value": &me.Value,
	})
}
