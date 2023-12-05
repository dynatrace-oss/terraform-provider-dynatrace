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

package masking

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	IsEmailMaskingEnabled     bool    `json:"isEmailMaskingEnabled"`     // Exclude email addresses from URLs
	IsFinancialMaskingEnabled bool    `json:"isFinancialMaskingEnabled"` // Exclude IBANs and payment card numbers from URLs
	IsNumbersMaskingEnabled   bool    `json:"isNumbersMaskingEnabled"`   // Exclude hexadecimal IDs and consecutive numbers above 5 digits from URLs
	IsQueryMaskingEnabled     bool    `json:"isQueryMaskingEnabled"`     // Exclude query parameters from URLs and web requests
	ProcessGroupID            *string `json:"-" scope:"processGroupId"`  // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
}

func (me *Settings) Name() string {
	if me.ProcessGroupID == nil {
		return "environment"
	}
	return *me.ProcessGroupID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"is_email_masking_enabled": {
			Type:        schema.TypeBool,
			Description: "Exclude email addresses from URLs",
			Required:    true,
		},
		"is_financial_masking_enabled": {
			Type:        schema.TypeBool,
			Description: "Exclude IBANs and payment card numbers from URLs",
			Required:    true,
		},
		"is_numbers_masking_enabled": {
			Type:        schema.TypeBool,
			Description: "Exclude hexadecimal IDs and consecutive numbers above 5 digits from URLs",
			Required:    true,
		},
		"is_query_masking_enabled": {
			Type:        schema.TypeBool,
			Description: "Exclude query parameters from URLs and web requests",
			Required:    true,
		},
		"process_group_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"is_email_masking_enabled":     me.IsEmailMaskingEnabled,
		"is_financial_masking_enabled": me.IsFinancialMaskingEnabled,
		"is_numbers_masking_enabled":   me.IsNumbersMaskingEnabled,
		"is_query_masking_enabled":     me.IsQueryMaskingEnabled,
		"process_group_id":             me.ProcessGroupID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"is_email_masking_enabled":     &me.IsEmailMaskingEnabled,
		"is_financial_masking_enabled": &me.IsFinancialMaskingEnabled,
		"is_numbers_masking_enabled":   &me.IsNumbersMaskingEnabled,
		"is_query_masking_enabled":     &me.IsQueryMaskingEnabled,
		"process_group_id":             &me.ProcessGroupID,
	})
}
