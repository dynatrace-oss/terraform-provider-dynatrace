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

package generalparameters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomErrorRules []*CustomErrorRule

func (me *CustomErrorRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_error_rule": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(CustomErrorRule).Schema()},
		},
	}
}

func (me CustomErrorRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("custom_error_rule", me)
}

func (me *CustomErrorRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("custom_error_rule", me)
}

type CustomErrorRule struct {
	Condition        *CompareOperation `json:"condition"`        // Request attribute condition
	RequestAttribute string            `json:"requestAttribute"` // Request attribute
}

func (me *CustomErrorRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeList,
			Description: "Request attribute condition",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(CompareOperation).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"request_attribute": {
			Type:        schema.TypeString,
			Description: "Request attribute",
			Required:    true,
		},
	}
}

func (me *CustomErrorRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":         me.Condition,
		"request_attribute": me.RequestAttribute,
	})
}

func (me *CustomErrorRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":         &me.Condition,
		"request_attribute": &me.RequestAttribute,
	})
}
