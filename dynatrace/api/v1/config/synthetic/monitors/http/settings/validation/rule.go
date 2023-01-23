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

package validation

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Rules is a list of validation rules
type Rules []*Rule

type Rule struct {
	Type        Type   `json:"type"`        // The type of the rule. Possible values are `patternConstraint`, `regexConstraint`, `httpStatusesList` and `certificateExpiryDateConstraint`
	PassIfFound bool   `json:"passIfFound"` // The validation condition. `true` means validation succeeds if the specified content/element is found. `false` means validation fails if the specified content/element is found. Always specify `false` for `certificateExpiryDateConstraint` to fail the monitor if SSL cedrtificate expiry is within the specified number of days
	Value       string `json:"value"`       // The content to look for
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the rule. Possible values are `patternConstraint`, `regexConstraint`, `httpStatusesList` and `certificateExpiryDateConstraint`",
			Required:    true,
		},
		"pass_if_found": {
			Type:        schema.TypeBool,
			Description: " The validation condition. `true` means validation succeeds if the specified content/element is found. `false` means validation fails if the specified content/element is found. Always specify `false` for `certificateExpiryDateConstraint` to fail the monitor if SSL certificate expiry is within the specified number of days",
			Optional:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The content to look for",
			Required:    true,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("pass_if_found", me.PassIfFound); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}
	return nil
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	if err := decoder.Decode("pass_if_found", &me.PassIfFound); err != nil {
		return err
	}
	if err := decoder.Decode("value", &me.Value); err != nil {
		return err
	}
	return nil
}

// Type The type of the rule. Possible values are `patternConstraint`, `regexConstraint`, `httpStatusesList` and `certificateExpiryDateConstraint`
type Type string

// ValidationRuleTypes offers the known enum values
var Types = struct {
	PatternConstraint               Type
	RegexConstraint                 Type
	HTTPStatusesList                Type
	CertificateExpiryDateConstraint Type
}{
	`patternConstraint`,
	`regexConstraint`,
	`httpStatusesList`,
	`certificateExpiryDateConstraint`,
}
