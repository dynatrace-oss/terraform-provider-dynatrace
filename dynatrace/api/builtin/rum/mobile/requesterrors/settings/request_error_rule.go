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

package requesterrors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RequestErrorRules []*RequestErrorRule

func (me *RequestErrorRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_rule": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(RequestErrorRule).Schema()},
		},
	}
}

func (me RequestErrorRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("error_rule", me)
}

func (me *RequestErrorRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("error_rule", me)
}

type RequestErrorRule struct {
	ErrorCodes string `json:"errorCodes"` // Exclude response codes
}

func (me *RequestErrorRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_codes": {
			Type:        schema.TypeString,
			Description: "Exclude response codes",
			Required:    true,
		},
	}
}

func (me *RequestErrorRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"error_codes": me.ErrorCodes,
	})
}

func (me *RequestErrorRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"error_codes": &me.ErrorCodes,
	})
}
