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

package externalwebrequest

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ContextRoot struct {
	ContributionType ContributionType       `json:"contributionType"`          // Possible Values: `OriginalValue`, `OverrideValue`, `TransformURL`, `TransformValue`
	SegmentCount     *int                   `json:"segmentCount,omitempty"`    // The number of segments of the URL to be kept. The URL is divided by slashes (/), the indexing starts with 1 at context root. For example, if you specify 2 for the `www.dynatrace.com/support/help/dynatrace-api/` URL, the value of `support/help` is used.
	Transformations  ReducedTransformations `json:"transformations,omitempty"` // Choose how the value will be transformed before contributing to the Service Id. All of the Transformations are always applied. Transformations are applied in the order they are specified, and the output of the previous transformation is the input for the next one. The resulting value contributes to the Service Id and can be found on the Service screen under **Properties and tags**.
	ValueOverride    *ValueOverride         `json:"valueOverride,omitempty"`   // The value to be used instead of the detected value.
}

func (me *ContextRoot) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"contribution_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `OriginalValue`, `OverrideValue`, `TransformURL`, `TransformValue`",
			Required:    true,
		},
		"segment_count": {
			Type:        schema.TypeInt,
			Description: "The number of segments of the URL to be kept. The URL is divided by slashes (/), the indexing starts with 1 at context root. For example, if you specify 2 for the `www.dynatrace.com/support/help/dynatrace-api/` URL, the value of `support/help` is used.",
			Optional:    true, // precondition
		},
		"transformations": {
			Type:        schema.TypeList,
			Description: "Choose how the value will be transformed before contributing to the Service Id. All of the Transformations are always applied. Transformations are applied in the order they are specified, and the output of the previous transformation is the input for the next one. The resulting value contributes to the Service Id and can be found on the Service screen under **Properties and tags**.",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(ReducedTransformations).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"value_override": {
			Type:        schema.TypeList,
			Description: "The value to be used instead of the detected value.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ValueOverride).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *ContextRoot) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"contribution_type": me.ContributionType,
		"segment_count":     me.SegmentCount,
		"transformations":   me.Transformations,
		"value_override":    me.ValueOverride,
	})
}

func (me *ContextRoot) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"contribution_type": &me.ContributionType,
		"segment_count":     &me.SegmentCount,
		"transformations":   &me.Transformations,
		"value_override":    &me.ValueOverride,
	})
	if me.SegmentCount == nil && me.ContributionType == ContributionTypes.Transformurl {
		me.SegmentCount = opt.NewInt(0)
	}
	return err
}
