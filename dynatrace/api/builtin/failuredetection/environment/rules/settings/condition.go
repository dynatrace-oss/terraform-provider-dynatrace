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

package rules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Conditions []*Condition

func (me *Conditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Condition).Schema()},
		},
	}
}

func (me Conditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("condition", me)
}

func (me *Conditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("condition", me)
}

type Condition struct {
	Attribute Attributes `json:"attribute"` // Possible Values: `PG_NAME`, `PG_TAG`, `SERVICE_MANAGEMENT_ZONE`, `SERVICE_NAME`, `SERVICE_TAG`, `SERVICE_TYPE`
	Predicate *Predicate `json:"predicate"` // Condition to check the attribute against
}

func (me *Condition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "Possible Values: `PG_NAME`, `PG_TAG`, `SERVICE_MANAGEMENT_ZONE`, `SERVICE_NAME`, `SERVICE_TAG`, `SERVICE_TYPE`",
			Required:    true,
		},
		"predicate": {
			Type:        schema.TypeList,
			Description: "Condition to check the attribute against",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Predicate).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Condition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"attribute": me.Attribute,
		"predicate": me.Predicate,
	})
}

func (me *Condition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"attribute": &me.Attribute,
		"predicate": &me.Predicate,
	})
}
