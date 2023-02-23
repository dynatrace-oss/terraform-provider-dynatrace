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

package workloaddetection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MatchFilter struct {
	MatchOperator MatchEnum `json:"matchOperator"`       // Possible Values: `CONTAINS`, `ENDS`, `EQUALS`, `EXISTS`, `NOT_CONTAINS`, `NOT_ENDS`, `NOT_EQUALS`, `NOT_STARTS`, `STARTS`
	Namespace     *string   `json:"namespace,omitempty"` // Namespace name
}

func (me *MatchFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"match_operator": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `ENDS`, `EQUALS`, `EXISTS`, `NOT_CONTAINS`, `NOT_ENDS`, `NOT_EQUALS`, `NOT_STARTS`, `STARTS`",
			Required:    true,
		},
		"namespace": {
			Type:        schema.TypeString,
			Description: "Namespace name",
			Optional:    true,
		},
	}
}

func (me *MatchFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"match_operator": me.MatchOperator,
		"namespace":      me.Namespace,
	})
}

func (me *MatchFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"match_operator": &me.MatchOperator,
		"namespace":      &me.Namespace,
	})
}
