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

package contextrules

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecurityContextRule struct {
	Query            string          `json:"query"`                      // Matcher
	RuleName         string          `json:"ruleName"`                   // Rule name
	Value            *string         `json:"value,omitempty"`            // Literal value to be set
	ValueSource      ValueSourceEnum `json:"valueSource"`                // Possible Values: `FIELD`, `LITERAL`
	ValueSourceField *string         `json:"valueSourceField,omitempty"` // Name of field used to copy value
}

func (me *SecurityContextRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"query": {
			Type:        schema.TypeString,
			Description: "Matcher",
			Required:    true,
		},
		"rule_name": {
			Type:        schema.TypeString,
			Description: "Rule name",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Literal value to be set",
			Optional:    true, // precondition
		},
		"value_source": {
			Type:        schema.TypeString,
			Description: "Possible Values: `FIELD`, `LITERAL`",
			Required:    true,
		},
		"value_source_field": {
			Type:        schema.TypeString,
			Description: "Name of field used to copy value",
			Optional:    true, // precondition
		},
	}
}

func (me *SecurityContextRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"query":              me.Query,
		"rule_name":          me.RuleName,
		"value":              me.Value,
		"value_source":       me.ValueSource,
		"value_source_field": me.ValueSourceField,
	})
}

func (me *SecurityContextRule) HandlePreconditions() error {
	if (me.Value == nil) && (string(me.ValueSource) == "LITERAL") {
		return fmt.Errorf("'value' must be specified if 'value_source' is set to '%v'", me.ValueSource)
	}
	if (me.ValueSourceField == nil) && (string(me.ValueSource) == "FIELD") {
		return fmt.Errorf("'value_source_field' must be specified if 'value_source' is set to '%v'", me.ValueSource)
	}
	return nil
}

func (me *SecurityContextRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"query":              &me.Query,
		"rule_name":          &me.RuleName,
		"value":              &me.Value,
		"value_source":       &me.ValueSource,
		"value_source_field": &me.ValueSourceField,
	})
}
