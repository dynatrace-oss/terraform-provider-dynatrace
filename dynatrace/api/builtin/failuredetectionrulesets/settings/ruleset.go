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

package failuredetectionrulesets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Ruleset struct {
	Condition                     *string                        `json:"condition,omitempty"` // Limits the scope of the failure detection ruleset using [DQL matcher](https://dt-url.net/l603wby) conditions on span and resource attributes.. A ruleset is applied only if the condition matches, otherwise the evaluation continues.\n\nIf empty, the condition will always match.
	Description                   *string                        `json:"description,omitempty"`
	FailOnCustomRules             CustomRules                    `json:"failOnCustomRules,omitempty"`   // Define failure reasons based on span and request attributes.. Failure detection result: `reason=\"custom_rule\"`, `verdict=\"failure\"`, `custom_rule_name`
	FailOnExceptions              *FailOnExceptions              `json:"failOnExceptions"`              // Evaluated expression: `iAny(`span.events`[][`span_event.name`] == \"exception\" and `span.events`[][`exception.escaped`] != false)`\n\nFailure detection result: `reason=\"exception\"`, `verdict=\"failure\"`, `exception_ids`
	FailOnGrpcStatusCodes         *FailOnGrpcStatusCodes         `json:"failOnGrpcStatusCodes"`         // Evaluated attribute: `rpc.grpc.status_code`\n\nFailure detection result: `reason=\"grpc_code\"`, `verdict=\"failure\"`
	FailOnHttpResponseStatusCodes *FailOnHttpResponseStatusCodes `json:"failOnHttpResponseStatusCodes"` // Evaluated attribute: `http.response.status_code`\n\nFailure detection result: `reason=\"http_code\"`, `verdict=\"failure\"`
	FailOnSpanStatusError         *FailOnSpanStatusError         `json:"failOnSpanStatusError"`         // Evaluated attribute: `span.status_code`\n\nFailure detection result: `reason=\"span_status\"`, `verdict=\"failure\"`
	Overrides                     *Overrides                     `json:"overrides"`
	RulesetName                   string                         `json:"rulesetName"` // Ruleset name
}

func (me *Ruleset) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "Limits the scope of the failure detection ruleset using [DQL matcher](https://dt-url.net/l603wby) conditions on span and resource attributes.. A ruleset is applied only if the condition matches, otherwise the evaluation continues.\n\nIf empty, the condition will always match.",
			Optional:    true, // nullable
		},
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"fail_on_custom_rules": {
			Type:        schema.TypeList,
			Description: "Define failure reasons based on span and request attributes.. Failure detection result: `reason=\"custom_rule\"`, `verdict=\"failure\"`, `custom_rule_name`",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(CustomRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fail_on_exceptions": {
			Type:        schema.TypeList,
			Description: "Evaluated expression: `iAny(`span.events`[][`span_event.name`] == \"exception\" and `span.events`[][`exception.escaped`] != false)`\n\nFailure detection result: `reason=\"exception\"`, `verdict=\"failure\"`, `exception_ids`",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FailOnExceptions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fail_on_grpc_status_codes": {
			Type:        schema.TypeList,
			Description: "Evaluated attribute: `rpc.grpc.status_code`\n\nFailure detection result: `reason=\"grpc_code\"`, `verdict=\"failure\"`",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FailOnGrpcStatusCodes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fail_on_http_response_status_codes": {
			Type:        schema.TypeList,
			Description: "Evaluated attribute: `http.response.status_code`\n\nFailure detection result: `reason=\"http_code\"`, `verdict=\"failure\"`",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FailOnHttpResponseStatusCodes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fail_on_span_status_error": {
			Type:        schema.TypeList,
			Description: "Evaluated attribute: `span.status_code`\n\nFailure detection result: `reason=\"span_status\"`, `verdict=\"failure\"`",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FailOnSpanStatusError).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"overrides": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Overrides).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"ruleset_name": {
			Type:        schema.TypeString,
			Description: "Ruleset name",
			Required:    true,
		},
	}
}

func (me *Ruleset) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":                          me.Condition,
		"description":                        me.Description,
		"fail_on_custom_rules":               me.FailOnCustomRules,
		"fail_on_exceptions":                 me.FailOnExceptions,
		"fail_on_grpc_status_codes":          me.FailOnGrpcStatusCodes,
		"fail_on_http_response_status_codes": me.FailOnHttpResponseStatusCodes,
		"fail_on_span_status_error":          me.FailOnSpanStatusError,
		"overrides":                          me.Overrides,
		"ruleset_name":                       me.RulesetName,
	})
}

func (me *Ruleset) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":                          &me.Condition,
		"description":                        &me.Description,
		"fail_on_custom_rules":               &me.FailOnCustomRules,
		"fail_on_exceptions":                 &me.FailOnExceptions,
		"fail_on_grpc_status_codes":          &me.FailOnGrpcStatusCodes,
		"fail_on_http_response_status_codes": &me.FailOnHttpResponseStatusCodes,
		"fail_on_span_status_error":          &me.FailOnSpanStatusError,
		"overrides":                          &me.Overrides,
		"ruleset_name":                       &me.RulesetName,
	})
}
