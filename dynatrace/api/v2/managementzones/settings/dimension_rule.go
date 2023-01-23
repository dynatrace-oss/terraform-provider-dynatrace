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

package managementzones

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No documentation available
type DimensionRule struct {
	AppliesTo  DimensionType       `json:"appliesTo"`            // Type
	Conditions DimensionConditions `json:"conditions,omitempty"` // Conditions
}

func (me *DimensionRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"applies_to": {
			Type:        schema.TypeString,
			Description: "Type",
			Required:    true,
		},
		"dimension_conditions": {
			Type:        schema.TypeList,
			Description: "Conditions",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(DimensionConditions).Schema()},
			Optional:    true,
		},
	}
}

func (me *DimensionRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"applies_to":           me.AppliesTo,
		"dimension_conditions": me.Conditions,
	})
}

func (me *DimensionRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"applies_to":           &me.AppliesTo,
		"dimension_conditions": &me.Conditions,
	})
}
