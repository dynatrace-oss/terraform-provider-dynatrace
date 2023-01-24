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

package errors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HTTPErrorRules []*HTTPErrorRule

func (me *HTTPErrorRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Description: "Configuration of the HTTP error in the web application",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(HTTPErrorRule).Schema()},
		},
	}
}

func (me HTTPErrorRules) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("rule", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *HTTPErrorRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

// HTTPErrorRule represents configuration of the HTTP error in the web application
type HTTPErrorRule struct {
	ConsiderUnknownErrorCode bool                 `json:"considerUnknownErrorCode"` // If `true`, match by errors that have unknown HTTP status code
	ConsiderBlockedRequests  bool                 `json:"considerBlockedRequests"`  // If `true`, match by errors that have CSP Rule violations
	ErrorCodes               *string              `json:"errorCodes,omitempty"`     // The HTTP status code or status code range to match by. \n\nThis field is required if **considerUnknownErrorCode** AND **considerBlockedRequests** are both set to `false`
	FilterByURL              bool                 `json:"filterByUrl"`              // If `true`, filter errors by URL
	Filter                   *HTTPErrorRuleFilter `json:"filter,omitempty"`         // The matching rule for the URL. Popssible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`.
	URL                      *string              `json:"url,omitempty"`            // The URL to look for
	Capture                  bool                 `json:"capture"`                  // Capture (`true`) or ignore (`false`) the error
	ImpactApdex              bool                 `json:"impactApdex"`              // Include (`true`) or exclude (`false`) the error in Apdex calculation
	ConsiderForAI            bool                 `json:"considerForAi"`            // Include (`true`) or exclude (`false`) the error in Davis AI [problem detection and analysis](https://dt-url.net/a963kd2)
}

func (me *HTTPErrorRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"consider_unknown_error_code": {
			Type:        schema.TypeBool,
			Description: "If `true`, match by errors that have unknown HTTP status code",
			Optional:    true,
		},
		"consider_blocked_requests": {
			Type:        schema.TypeBool,
			Description: "If `true`, match by errors that have CSP Rule violations",
			Optional:    true,
		},
		"error_codes": {
			Type:        schema.TypeString,
			Description: "The HTTP status code or status code range to match by. \n\nThis field is required if **considerUnknownErrorCode** AND **considerBlockedRequests** are both set to `false`",
			Optional:    true,
		},
		"filter_by_url": {
			Type:        schema.TypeBool,
			Description: "If `true`, filter errors by URL",
			Optional:    true,
		},
		"filter": {
			Type:        schema.TypeString,
			Description: "The matching rule for the URL. Popssible values are `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH` and `EQUALS`.",
			Optional:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL to look for",
			Optional:    true,
		},
		"capture": {
			Type:        schema.TypeBool,
			Description: "Capture (`true`) or ignore (`false`) the error",
			Optional:    true,
		},
		"impact_apdex": {
			Type:        schema.TypeBool,
			Description: "Include (`true`) or exclude (`false`) the error in Apdex calculation",
			Optional:    true,
		},
		"consider_for_ai": {
			Type:        schema.TypeBool,
			Description: "Include (`true`) or exclude (`false`) the error in Davis AI [problem detection and analysis](https://dt-url.net/a963kd2)",
			Optional:    true,
		},
	}
}

func (me *HTTPErrorRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"consider_unknown_error_code": me.ConsiderUnknownErrorCode,
		"consider_blocked_requests":   me.ConsiderBlockedRequests,
		"error_codes":                 me.ErrorCodes,
		"filter_by_url":               me.FilterByURL,
		"filter":                      me.Filter,
		"url":                         me.URL,
		"capture":                     me.Capture,
		"impact_apdex":                me.ImpactApdex,
		"consider_for_ai":             me.ConsiderForAI,
	})
}

func (me *HTTPErrorRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"consider_unknown_error_code": &me.ConsiderUnknownErrorCode,
		"consider_blocked_requests":   &me.ConsiderBlockedRequests,
		"error_codes":                 &me.ErrorCodes,
		"filter_by_url":               &me.FilterByURL,
		"filter":                      &me.Filter,
		"url":                         &me.URL,
		"capture":                     &me.Capture,
		"impact_apdex":                &me.ImpactApdex,
		"consider_for_ai":             &me.ConsiderForAI,
	})
}
