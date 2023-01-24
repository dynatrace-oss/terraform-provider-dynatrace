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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceTimingSettings configures resource timings capture
type ResourceTimingSettings struct {
	W3CResourceTimings                        bool                       `json:"w3cResourceTimings"`                        // W3C resource timings for third party/CDN enabled/disabled
	NonW3CResourceTimings                     bool                       `json:"nonW3cResourceTimings"`                     // Timing for JavaScript files and images on non-W3C supported browsers enabled/disabled
	NonW3CResourceTimingsInstrumentationDelay int32                      `json:"nonW3cResourceTimingsInstrumentationDelay"` // Instrumentation delay for monitoring resource and image resource impact in browsers that don't offer W3C resource timings. \n\nValid values range from 0 to 9999.\n\nOnly effective if **nonW3cResourceTimings** is enabled
	ResourceTimingCaptureType                 *ResourceTimingCaptureType `json:"resourceTimingCaptureType"`                 // Defines how detailed resource timings are captured.\n\nOnly effective if **w3cResourceTimings** or **nonW3cResourceTimings** is enabled. Possible values are `CAPTURE_ALL_SUMMARIES`, `CAPTURE_FULL_DETAILS` and `CAPTURE_LIMITED_SUMMARIES`
	ResourceTimingsDomainLimit                *int32                     `json:"resourceTimingsDomainLimit"`                // Limits the number of domains for which W3C resource timings are captured.\n\nOnly effective if **resourceTimingCaptureType** is `CAPTURE_LIMITED_SUMMARIES`. Valid values range from 0 to 50.
}

func (me *ResourceTimingSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"w3c_resource_timings": {
			Type:        schema.TypeBool,
			Description: "W3C resource timings for third party/CDN enabled/disabled",
			Optional:    true,
		},
		"non_w3c_resource_timings": {
			Type:        schema.TypeBool,
			Description: "Timing for JavaScript files and images on non-W3C supported browsers enabled/disabled",
			Optional:    true,
		},
		"instrumentation_delay": {
			Type:        schema.TypeInt,
			Description: "Instrumentation delay for monitoring resource and image resource impact in browsers that don't offer W3C resource timings. \n\nValid values range from 0 to 9999.\n\nOnly effective if `nonW3cResourceTimings` is enabled",
			Required:    true,
		},
		"resource_timing_capture_type": {
			Type:        schema.TypeString,
			Description: "Defines how detailed resource timings are captured.\n\nOnly effective if **w3cResourceTimings** or **nonW3cResourceTimings** is enabled. Possible values are `CAPTURE_ALL_SUMMARIES`, `CAPTURE_FULL_DETAILS` and `CAPTURE_LIMITED_SUMMARIES`",
			Optional:    true,
		},
		"resource_timings_domain_limit": {
			Type:        schema.TypeInt,
			Description: "Limits the number of domains for which W3C resource timings are captured.\n\nOnly effective if **resourceTimingCaptureType** is `CAPTURE_LIMITED_SUMMARIES`. Valid values range from 0 to 50.",
			Optional:    true,
		},
	}
}

func (me *ResourceTimingSettings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"w3c_resource_timings":          me.W3CResourceTimings,
		"non_w3c_resource_timings":      me.NonW3CResourceTimings,
		"instrumentation_delay":         me.NonW3CResourceTimingsInstrumentationDelay,
		"resource_timing_capture_type":  me.ResourceTimingCaptureType,
		"resource_timings_domain_limit": me.ResourceTimingsDomainLimit,
	}); err != nil {
		return err
	}
	if me.ResourceTimingCaptureType != nil && len(string(*me.ResourceTimingCaptureType)) == 0 {
		delete(properties, "resource_timing_capture_type")
	}
	return nil
}

func (me *ResourceTimingSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"w3c_resource_timings":          &me.W3CResourceTimings,
		"non_w3c_resource_timings":      &me.NonW3CResourceTimings,
		"instrumentation_delay":         &me.NonW3CResourceTimingsInstrumentationDelay,
		"resource_timing_capture_type":  &me.ResourceTimingCaptureType,
		"resource_timings_domain_limit": &me.ResourceTimingsDomainLimit,
	})
	if err != nil {
		return err
	}
	return nil
}
