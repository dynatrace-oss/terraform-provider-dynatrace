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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type PublicDomainTransformationSet struct {
	ContributionType ContributionTypeWithOverride `json:"contributionType"`           // Possible Values: `OriginalValue`, `OverrideValue`, `TransformValue`
	CopyFromHostName *bool                        `json:"copyFromHostName,omitempty"` // Use the detected host name instead of the request's domain name.
	Transformations  Transformations              `json:"transformations,omitempty"`  // Choose how to transform a value before it contributes to the Service Id. Note that all of the Transformations are always applied. Transformations are applied in the order they are specified, and the output of the previous transformation is the input for the next one. The resulting value contributes to the Service Id and can be found on the **Service overview page** under **Properties and tags**.
	ValueOverride    *ValueOverride               `json:"valueOverride,omitempty"`    // The value to be used instead of the detected value.
}

func (me *PublicDomainTransformationSet) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"contribution_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `OriginalValue`, `OverrideValue`, `TransformValue`",
			Required:    true,
		},
		"copy_from_host_name": {
			Type:        schema.TypeBool,
			Description: "Use the detected host name instead of the request's domain name.",
			Optional:    true, // precondition
		},
		"transformations": {
			Type:        schema.TypeList,
			Description: "Choose how to transform a value before it contributes to the Service Id. Note that all of the Transformations are always applied. Transformations are applied in the order they are specified, and the output of the previous transformation is the input for the next one. The resulting value contributes to the Service Id and can be found on the **Service overview page** under **Properties and tags**.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(Transformations).Schema()},
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

func (me *PublicDomainTransformationSet) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"contribution_type":   me.ContributionType,
		"copy_from_host_name": me.CopyFromHostName,
		"transformations":     me.Transformations,
		"value_override":      me.ValueOverride,
	})
}

func (me *PublicDomainTransformationSet) HandlePreconditions() error {
	if me.CopyFromHostName == nil && slices.Contains([]string{"OriginalValue", "TransformValue"}, string(me.ContributionType)) {
		me.CopyFromHostName = opt.NewBool(false)
	}
	if me.ValueOverride == nil && (string(me.ContributionType) == "OverrideValue") {
		return fmt.Errorf("'value_override' must be specified if 'contribution_type' is set to '%v'", me.ContributionType)
	}
	if me.ValueOverride != nil && (string(me.ContributionType) != "OverrideValue") {
		return fmt.Errorf("'value_override' must not be specified if 'contribution_type' is set to '%v'", me.ContributionType)
	}
	// ---- Transformations Transformations -> {"expectedValue":"TransformValue","property":"contributionType","type":"EQUALS"}
	return nil
}

func (me *PublicDomainTransformationSet) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"contribution_type":   &me.ContributionType,
		"copy_from_host_name": &me.CopyFromHostName,
		"transformations":     &me.Transformations,
		"value_override":      &me.ValueOverride,
	})
}
