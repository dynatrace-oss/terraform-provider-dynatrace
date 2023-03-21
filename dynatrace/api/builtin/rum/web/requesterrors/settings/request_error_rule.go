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

package requesterrors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RequestErrorRules []*RequestErrorRule

func (me *RequestErrorRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(RequestErrorRule).Schema()},
		},
	}
}

func (me RequestErrorRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("error_rule", me)
}

func (me *RequestErrorRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("error_rule", me)
}

type RequestErrorRule struct {
	CaptureSettings       *CaptureSettings `json:"captureSettings"`       // Capture settings
	ConsiderCspViolations bool             `json:"considerCspViolations"` // Match by errors that have CSP violations
	ConsiderFailedImages  bool             `json:"considerFailedImages"`  // Match by errors that have failed image requests
	ErrorCodes            *string          `json:"errorCodes,omitempty"`  // Match by error code
	FilterSettings        *FilterSettings  `json:"filterSettings"`        // Filter settings
}

func (me *RequestErrorRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"capture_settings": {
			Type:        schema.TypeList,
			Description: "Capture settings",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(CaptureSettings).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"consider_csp_violations": {
			Type:        schema.TypeBool,
			Description: "Match by errors that have CSP violations",
			Required:    true,
		},
		"consider_failed_images": {
			Type:        schema.TypeBool,
			Description: "Match by errors that have failed image requests",
			Required:    true,
		},
		"error_codes": {
			Type:        schema.TypeString,
			Description: "Match by error code",
			Optional:    true,
		},
		"filter_settings": {
			Type:        schema.TypeList,
			Description: "Filter settings",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(FilterSettings).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *RequestErrorRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"capture_settings":        me.CaptureSettings,
		"consider_csp_violations": me.ConsiderCspViolations,
		"consider_failed_images":  me.ConsiderFailedImages,
		"error_codes":             me.ErrorCodes,
		"filter_settings":         me.FilterSettings,
	})
}

func (me *RequestErrorRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"capture_settings":        &me.CaptureSettings,
		"consider_csp_violations": &me.ConsiderCspViolations,
		"consider_failed_images":  &me.ConsiderFailedImages,
		"error_codes":             &me.ErrorCodes,
		"filter_settings":         &me.FilterSettings,
	})
}
