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

package keyperformancemetrics

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	FrustratingIfReportedOrWebRequestError bool        `json:"frustratingIfReportedOrWebRequestError"` // Treat user actions with reported errors or web request errors as erroneous and rate their performance as Frustrating. Turn off this setting if errors should not affect the Apdex rate.
	Scope                                  string      `json:"-" scope:"scope"`                        // The scope of this setting (DEVICE_APPLICATION_METHOD, MOBILE_APPLICATION, CUSTOM_APPLICATION)
	Thresholds                             *Thresholds `json:"thresholds"`
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"frustrating_if_reported_or_web_request_error": {
			Type:        schema.TypeBool,
			Description: "Treat user actions with reported errors or web request errors as erroneous and rate their performance as Frustrating. Turn off this setting if errors should not affect the Apdex rate.",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (DEVICE_APPLICATION_METHOD, MOBILE_APPLICATION, CUSTOM_APPLICATION)",
			Required:    true,
			ForceNew:    true,
		},
		"thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Thresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"frustrating_if_reported_or_web_request_error": me.FrustratingIfReportedOrWebRequestError,
		"scope":      me.Scope,
		"thresholds": me.Thresholds,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"frustrating_if_reported_or_web_request_error": &me.FrustratingIfReportedOrWebRequestError,
		"scope":      &me.Scope,
		"thresholds": &me.Thresholds,
	})
}
