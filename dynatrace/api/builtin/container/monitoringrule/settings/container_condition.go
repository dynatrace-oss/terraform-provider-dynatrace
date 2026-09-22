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

package monitoringrule

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ContainerCondition struct {
	Operator ConditionOperator `json:"operator"`        // Condition operator. Possible values: `CONTAINS`, `ENDS`, `EQUALS`, `EXISTS`, `NOT_CONTAINS`, `NOT_ENDS`, `NOT_EQUALS`, `NOT_EXISTS`, `NOT_STARTS`, `STARTS`
	Property ContainerItem     `json:"property"`        // Container property. Possible values: `CONTAINER_NAME`, `IMAGE_NAME`, `KUBERNETES_BASEPODNAME`, `KUBERNETES_CONTAINERNAME`, `KUBERNETES_FULLPODNAME`, `KUBERNETES_NAMESPACE`, `KUBERNETES_PODUID`
	Value    *string           `json:"value,omitempty"` // Condition value
}

func (me *ContainerCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"operator": {
			Type:        schema.TypeString,
			Description: "Condition operator. Possible values: `CONTAINS`, `ENDS`, `EQUALS`, `EXISTS`, `NOT_CONTAINS`, `NOT_ENDS`, `NOT_EQUALS`, `NOT_EXISTS`, `NOT_STARTS`, `STARTS`",
			Required:    true,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Container property. Possible values: `CONTAINER_NAME`, `IMAGE_NAME`, `KUBERNETES_BASEPODNAME`, `KUBERNETES_CONTAINERNAME`, `KUBERNETES_FULLPODNAME`, `KUBERNETES_NAMESPACE`, `KUBERNETES_PODUID`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Condition value",
			Optional:    true, // precondition
		},
	}
}

func (me *ContainerCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"operator": me.Operator,
		"property": me.Property,
		"value":    me.Value,
	})
}

func (me *ContainerCondition) HandlePreconditions() error {
	if (me.Value != nil) && (slices.Contains([]string{"EXISTS", "NOT_EXISTS"}, string(me.Operator))) {
		return fmt.Errorf("'value' must not be specified unless 'operator' is not one of ['EXISTS', 'NOT_EXISTS']; got 'operator'='%v'", me.Operator)
	}
	if (me.Value == nil) && (!slices.Contains([]string{"EXISTS", "NOT_EXISTS"}, string(me.Operator))) {
		return fmt.Errorf("'value' must be specified when 'operator' is not one of ['EXISTS', 'NOT_EXISTS']; got 'operator'='%v'", me.Operator)
	}
	return nil
}

func (me *ContainerCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"operator": &me.Operator,
		"property": &me.Property,
		"value":    &me.Value,
	})
}
