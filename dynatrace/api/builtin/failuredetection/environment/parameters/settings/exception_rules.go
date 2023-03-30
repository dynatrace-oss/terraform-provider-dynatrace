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

package parameters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ExceptionRules struct {
	CustomErrorRules           CustomErrorRules `json:"customErrorRules,omitempty"`         // Some custom error situations are only detectable via a return value or other means. To support such cases, [define a request attribute](https://dt-url.net/ys5k0p4y) that captures the required data. Then define a custom error rule that determines if the request has failed based on the value of the request attribute.
	CustomHandledExceptions    Exceptions       `json:"customHandledExceptions,omitempty"`  // There may be situations where your application code handles exceptions gracefully in a manner that these failures aren't detected by Dynatrace. Use this setting to define specific gracefully-handled exceptions that should be treated as service failures.
	IgnoreAllExceptions        bool             `json:"ignoreAllExceptions"`                // Ignore all exceptions
	IgnoreSpanFailureDetection bool             `json:"ignoreSpanFailureDetection"`         // Ignore span failure detection
	IgnoredExceptions          Exceptions       `json:"ignoredExceptions,omitempty"`        // Some exceptions that are thrown by legacy or 3rd-party code indicate a specific response, not an error. Use this setting to instruct Dynatrace to treat such exceptions as non-failed requests.. If an exception matching any of the defined patterns occurs in a request, it will not be considered as a failure. Other exceptions occurring at the same request might still mark the request as failed.
	SuccessForcingExceptions   Exceptions       `json:"successForcingExceptions,omitempty"` // Define exceptions which indicate that a service call should not be considered as failed. E.g. an exception indicating that the client aborted the operation.. If an exception matching any of the defined patterns occurs on the entry node of the service, it will be considered successful. Compared to ignored exceptions, the request will be considered successful even if other exceptions occur in the same request.
}

func (me *ExceptionRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_error_rules": {
			Type:        schema.TypeList,
			Description: "Some custom error situations are only detectable via a return value or other means. To support such cases, [define a request attribute](https://dt-url.net/ys5k0p4y) that captures the required data. Then define a custom error rule that determines if the request has failed based on the value of the request attribute.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(CustomErrorRules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"custom_handled_exceptions": {
			Type:        schema.TypeList,
			Description: "There may be situations where your application code handles exceptions gracefully in a manner that these failures aren't detected by Dynatrace. Use this setting to define specific gracefully-handled exceptions that should be treated as service failures.",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(Exceptions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"ignore_all_exceptions": {
			Type:        schema.TypeBool,
			Description: "Ignore all exceptions",
			Required:    true,
		},
		"ignore_span_failure_detection": {
			Type:        schema.TypeBool,
			Description: "Ignore span failure detection",
			Required:    true,
		},
		"ignored_exceptions": {
			Type:        schema.TypeList,
			Description: "Some exceptions that are thrown by legacy or 3rd-party code indicate a specific response, not an error. Use this setting to instruct Dynatrace to treat such exceptions as non-failed requests.. If an exception matching any of the defined patterns occurs in a request, it will not be considered as a failure. Other exceptions occurring at the same request might still mark the request as failed.",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(Exceptions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"success_forcing_exceptions": {
			Type:        schema.TypeList,
			Description: "Define exceptions which indicate that a service call should not be considered as failed. E.g. an exception indicating that the client aborted the operation.. If an exception matching any of the defined patterns occurs on the entry node of the service, it will be considered successful. Compared to ignored exceptions, the request will be considered successful even if other exceptions occur in the same request.",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(Exceptions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *ExceptionRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"custom_error_rules":            me.CustomErrorRules,
		"custom_handled_exceptions":     me.CustomHandledExceptions,
		"ignore_all_exceptions":         me.IgnoreAllExceptions,
		"ignore_span_failure_detection": me.IgnoreSpanFailureDetection,
		"ignored_exceptions":            me.IgnoredExceptions,
		"success_forcing_exceptions":    me.SuccessForcingExceptions,
	})
}

func (me *ExceptionRules) HandlePreconditions() {
	// ---- CustomHandledExceptions Exceptions -> {"expectedValue":false,"property":"ignoreAllExceptions","type":"EQUALS"}
	// ---- IgnoredExceptions Exceptions -> {"expectedValue":false,"property":"ignoreAllExceptions","type":"EQUALS"}
	// ---- SuccessForcingExceptions Exceptions -> {"expectedValue":false,"property":"ignoreAllExceptions","type":"EQUALS"}
}

func (me *ExceptionRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_error_rules":            &me.CustomErrorRules,
		"custom_handled_exceptions":     &me.CustomHandledExceptions,
		"ignore_all_exceptions":         &me.IgnoreAllExceptions,
		"ignore_span_failure_detection": &me.IgnoreSpanFailureDetection,
		"ignored_exceptions":            &me.IgnoredExceptions,
		"success_forcing_exceptions":    &me.SuccessForcingExceptions,
	})
}
