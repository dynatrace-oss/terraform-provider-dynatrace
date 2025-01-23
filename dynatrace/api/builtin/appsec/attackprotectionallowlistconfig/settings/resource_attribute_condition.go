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
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceAttributeCondition struct {
	Matcher                ResourceAttributeValueMatcher `json:"matcher"`                          // Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EXIST`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `EXISTS`, `NOT_EQUALS`, `STARTS_WITH`
	ResourceAttributeKey   string                        `json:"resourceAttributeKey"`             // Resource attribute key
	ResourceAttributeValue *string                       `json:"resourceAttributeValue,omitempty"` // Resource attribute value
}

func (me *ResourceAttributeCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"matcher": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EXIST`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `EXISTS`, `NOT_EQUALS`, `STARTS_WITH`",
			Required:    true,
		},
		"resource_attribute_key": {
			Type:        schema.TypeString,
			Description: "Resource attribute key",
			Required:    true,
		},
		"resource_attribute_value": {
			Type:        schema.TypeString,
			Description: "Resource attribute value",
			Optional:    true, // nullable & precondition
		},
	}
}

func (me *ResourceAttributeCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"matcher":                  me.Matcher,
		"resource_attribute_key":   me.ResourceAttributeKey,
		"resource_attribute_value": me.ResourceAttributeValue,
	})
}

func (me *ResourceAttributeCondition) HandlePreconditions() error {
	if (me.ResourceAttributeValue == nil) && (!slices.Contains([]string{"EXISTS", "DOES_NOT_EXIST"}, string(me.Matcher))) {
		return fmt.Errorf("'resource_attribute_value' must be specified if 'matcher' is set to '%v'", me.Matcher)
	}
	return nil
}

func (me *ResourceAttributeCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"matcher":                  &me.Matcher,
		"resource_attribute_key":   &me.ResourceAttributeKey,
		"resource_attribute_value": &me.ResourceAttributeValue,
	})
}
