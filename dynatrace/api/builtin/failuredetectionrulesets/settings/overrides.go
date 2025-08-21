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

type Overrides struct {
	ForceSuccessOnExceptions              SingleExceptions                       `json:"forceSuccessOnExceptions,omitempty"`              // Define escaped exceptions that should force success.. Evaluated expression: `iAny(`span.events`[][`span_event.name`] == \"exception\" and `span.events`[][`exception.escaped`] != false)`\n\nFailure detection result: `reason=\"exception\"`, `verdict=\"success\"`, `exception_ids`
	ForceSuccessOnGrpcResponseStatusCodes *ForceSuccessOnGrpcResponseStatusCodes `json:"forceSuccessOnGrpcResponseStatusCodes,omitempty"` // Evaluated attribute: `rpc.grpc.status_code`\n\nFailure detection result: `reason=\"grpc_code\"`, `verdict=\"success\"`
	ForceSuccessOnHttpResponseStatusCodes *ForceSuccessOnHttpResponseStatusCodes `json:"forceSuccessOnHttpResponseStatusCodes,omitempty"` // Evaluated attribute: `http.response.status_code`\n\nFailure detection result: `reason=\"http_code\"`, `verdict=\"success\"`
	ForceSuccessOnSpanStatusOk            *ForceSuccessOnSpanStatusOk            `json:"forceSuccessOnSpanStatusOk"`                      // Evaluated attribute: `span.status_code`\n\nFailure detection result: `reason=\"span_status\"`, `verdict=\"success\"`
	ForceSuccessWithCustomRules           CustomRules                            `json:"forceSuccessWithCustomRules,omitempty"`           // Override failures based on span and request attribute conditions.. Failure detection result: `reason=\"custom_rule\"`, `verdict=\"success\"`, `custom_rule_name`
}

func (me *Overrides) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"force_success_on_exceptions": {
			Type:        schema.TypeList,
			Description: "Define escaped exceptions that should force success.. Evaluated expression: `iAny(`span.events`[][`span_event.name`] == \"exception\" and `span.events`[][`exception.escaped`] != false)`\n\nFailure detection result: `reason=\"exception\"`, `verdict=\"success\"`, `exception_ids`",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(SingleExceptions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"force_success_on_grpc_response_status_codes": {
			Type:             schema.TypeList,
			Description:      "Evaluated attribute: `rpc.grpc.status_code`\n\nFailure detection result: `reason=\"grpc_code\"`, `verdict=\"success\"`",
			Optional:         true,
			Elem:             &schema.Resource{Schema: new(ForceSuccessOnGrpcResponseStatusCodes).Schema()},
			MinItems:         1,
			MaxItems:         1,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return new == "0" },
		},
		"force_success_on_http_response_status_codes": {
			Type:             schema.TypeList,
			Description:      "Evaluated attribute: `http.response.status_code`\n\nFailure detection result: `reason=\"http_code\"`, `verdict=\"success\"`",
			Optional:         true,
			Elem:             &schema.Resource{Schema: new(ForceSuccessOnHttpResponseStatusCodes).Schema()},
			MinItems:         1,
			MaxItems:         1,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return new == "0" },
		},
		"force_success_on_span_status_ok": {
			Type:        schema.TypeList,
			Description: "Evaluated attribute: `span.status_code`\n\nFailure detection result: `reason=\"span_status\"`, `verdict=\"success\"`",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ForceSuccessOnSpanStatusOk).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"force_success_with_custom_rules": {
			Type:        schema.TypeList,
			Description: "Override failures based on span and request attribute conditions.. Failure detection result: `reason=\"custom_rule\"`, `verdict=\"success\"`, `custom_rule_name`",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(CustomRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Overrides) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"force_success_on_exceptions":                 me.ForceSuccessOnExceptions,
		"force_success_on_grpc_response_status_codes": me.ForceSuccessOnGrpcResponseStatusCodes,
		"force_success_on_http_response_status_codes": me.ForceSuccessOnHttpResponseStatusCodes,
		"force_success_on_span_status_ok":             me.ForceSuccessOnSpanStatusOk,
		"force_success_with_custom_rules":             me.ForceSuccessWithCustomRules,
	})
}

func (me *Overrides) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"force_success_on_exceptions":                 &me.ForceSuccessOnExceptions,
		"force_success_on_grpc_response_status_codes": &me.ForceSuccessOnGrpcResponseStatusCodes,
		"force_success_on_http_response_status_codes": &me.ForceSuccessOnHttpResponseStatusCodes,
		"force_success_on_span_status_ok":             &me.ForceSuccessOnSpanStatusOk,
		"force_success_with_custom_rules":             &me.ForceSuccessWithCustomRules,
	})
}
