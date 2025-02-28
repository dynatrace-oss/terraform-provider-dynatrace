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

package attackprotectionallowlistconfig

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AgentSideCriterias []*AgentSideCriteria

func (me *AgentSideCriterias) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(AgentSideCriteria).Schema()},
		},
	}
}

func (me AgentSideCriterias) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *AgentSideCriterias) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

type AgentSideCriteria struct {
	CriteriaKey                AgentSideAttributeKey     `json:"criteriaKey"`                          // Possible Values: `ACTOR_IP`, `DETECTION_TYPE`, `ENTRY_POINT_PAYLOAD`, `ENTRY_POINT_PAYLOAD_DOMAIN`, `ENTRY_POINT_PAYLOAD_PORT`, `ENTRY_POINT_URL_PATH`
	CriteriaMatcher            AgentSideAttributeMatcher `json:"criteriaMatcher"`                      // Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_STARTS_WITH`, `ENDS_WITH`, `EQUALS`, `IP_CIDR`, `NOT_EQUALS`, `NOT_IN_IP_CIDR`, `STARTS_WITH`
	CriteriaValueDetectionType *DetectionType            `json:"criteriaValueDetectionType,omitempty"` // Possible Values: `CMD_INJECTION`, `JNDI_INJECTION`, `SQL_INJECTION`, `SSRF`
	CriteriaValueFreeText      *string                   `json:"criteriaValueFreeText,omitempty"`      // Value
}

func (me *AgentSideCriteria) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"criteria_key": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ACTOR_IP`, `DETECTION_TYPE`, `ENTRY_POINT_PAYLOAD`, `ENTRY_POINT_PAYLOAD_DOMAIN`, `ENTRY_POINT_PAYLOAD_PORT`, `ENTRY_POINT_URL_PATH`",
			Required:    true,
		},
		"criteria_matcher": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_STARTS_WITH`, `ENDS_WITH`, `EQUALS`, `IP_CIDR`, `NOT_EQUALS`, `NOT_IN_IP_CIDR`, `STARTS_WITH`",
			Required:    true,
		},
		"criteria_value_detection_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CMD_INJECTION`, `JNDI_INJECTION`, `SQL_INJECTION`, `SSRF`",
			Optional:    true, // precondition
		},
		"criteria_value_free_text": {
			Type:        schema.TypeString,
			Description: "Value",
			Optional:    true, // precondition
		},
	}
}

func (me *AgentSideCriteria) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"criteria_key":                  me.CriteriaKey,
		"criteria_matcher":              me.CriteriaMatcher,
		"criteria_value_detection_type": me.CriteriaValueDetectionType,
		"criteria_value_free_text":      me.CriteriaValueFreeText,
	})
}

func (me *AgentSideCriteria) HandlePreconditions() error {
	if (me.CriteriaValueDetectionType == nil) && (string(me.CriteriaKey) == "DETECTION_TYPE") {
		return fmt.Errorf("'criteria_value_detection_type' must be specified if 'criteria_key' is set to '%v'", me.CriteriaKey)
	}
	if (me.CriteriaValueDetectionType != nil) && (string(me.CriteriaKey) != "DETECTION_TYPE") {
		return fmt.Errorf("'criteria_value_detection_type' must not be specified if 'criteria_key' is set to '%v'", me.CriteriaKey)
	}
	if (me.CriteriaValueFreeText == nil) && (string(me.CriteriaKey) != "DETECTION_TYPE") {
		return fmt.Errorf("'criteria_value_free_text' must be specified if 'criteria_key' is set to '%v'", me.CriteriaKey)
	}
	return nil
}

func (me *AgentSideCriteria) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"criteria_key":                  &me.CriteriaKey,
		"criteria_matcher":              &me.CriteriaMatcher,
		"criteria_value_detection_type": &me.CriteriaValueDetectionType,
		"criteria_value_free_text":      &me.CriteriaValueFreeText,
	})
}
