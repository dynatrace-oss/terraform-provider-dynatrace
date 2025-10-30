/**
* @license
* Copyright 2025 Dynatrace LLC
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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CostAllocationAttributes struct {
	Value *GenericValueAssignment `json:"value"` // The strategy to set the cost allocation field
}

func (me *CostAllocationAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:        schema.TypeList,
			Description: "The strategy to set the cost allocation field",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(GenericValueAssignment).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *CostAllocationAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"value": me.Value,
	})
}

func (me *CostAllocationAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"value": &me.Value,
	})
}
