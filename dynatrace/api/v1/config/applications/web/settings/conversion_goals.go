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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings/useraction"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings/visit"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConversionGoals []*ConversionGoal

func (me *ConversionGoals) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"goal": {
			Type:        schema.TypeList,
			Description: "A conversion goal of the application",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ConversionGoal).Schema()},
		},
	}
}

func (me ConversionGoals) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("goal", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *ConversionGoals) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("goal", me)
}

// ConversionGoal A conversion goal of the application
type ConversionGoal struct {
	Name                  string                  `json:"name"`                            // The name of the conversion goal. Valid length within 1 and 50 characters.
	ID                    *string                 `json:"id,omitempty"`                    // The ID of conversion goal. \n\n Omit it while creating a new conversion goal
	Type                  *ConversionGoalType     `json:"type,omitempty"`                  // The type of the conversion goal. Possible values are `Destination`, `UserAction`, `VisitDuration` and `VisitNumActions`
	DestinationDetails    *DestinationDetails     `json:"destinationDetails,omitempty"`    // Configuration of a destination-based conversion goal
	UserActionDetails     *useraction.Details     `json:"userActionDetails,omitempty"`     // Configuration of a user action-based conversion goal
	VisitDurationDetails  *visit.DurationDetails  `json:"visitDurationDetails,omitempty"`  // Configuration of a visit duration-based conversion goal
	VisitNumActionDetails *visit.NumActionDetails `json:"visitNumActionDetails,omitempty"` // Configuration of a number of user actions-based conversion goal
}

func (me *ConversionGoal) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the conversion goal. Valid length within 1 and 50 characters.",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of conversion goal. \n\n Omit it while creating a new conversion goal",
			Optional:    true,
			Computed:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the conversion goal. Possible values are `Destination`, `UserAction`, `VisitDuration` and `VisitNumActions`",
			Optional:    true,
		},
		"destination": {
			Type:        schema.TypeList,
			Description: "Configuration of a destination-based conversion goal",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(DestinationDetails).Schema()},
		},
		"user_action": {
			Type:        schema.TypeList,
			Description: "Configuration of a destination-based conversion goal",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(useraction.Details).Schema()},
		},
		"visit_duration": {
			Type:        schema.TypeList,
			Description: "Configuration of a destination-based conversion goal",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(visit.DurationDetails).Schema()},
		},
		"visit_num_action": {
			Type:        schema.TypeList,
			Description: "Configuration of a destination-based conversion goal",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(visit.NumActionDetails).Schema()},
		},
	}
}

func (me *ConversionGoal) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name": me.Name,
		// "id":               me.ID,
		"type":             me.Type,
		"destination":      me.DestinationDetails,
		"user_action":      me.UserActionDetails,
		"visit_duration":   me.VisitDurationDetails,
		"visit_num_action": me.VisitNumActionDetails,
	})
}

func (me *ConversionGoal) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":             &me.Name,
		"id":               &me.ID,
		"type":             &me.Type,
		"destination":      &me.DestinationDetails,
		"user_action":      &me.UserActionDetails,
		"visit_duration":   &me.VisitDurationDetails,
		"visit_num_action": &me.VisitNumActionDetails,
	})
}
