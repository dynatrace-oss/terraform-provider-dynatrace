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

// Defines the Apdex settings of an application
type ApdexSettings struct {
	Threshold                    *int `json:"threshold,omitempty"`                    // no documentation available
	ToleratedThreshold           *int `json:"toleratedThreshold,omitempty"`           // Maximal value of apdex, which is considered as satisfied user experience. Values between 0 and 60000 are allowed.
	FrustratingThreshold         *int `json:"frustratingThreshold,omitempty"`         // Maximal value of apdex, which is considered as tolerable user experience. Values between 0 and 240000 are allowed.
	ToleratedFallbackThreshold   *int `json:"toleratedFallbackThreshold,omitempty"`   // Fallback threshold of an XHR action, defining a satisfied user experience, when the configured KPM is not available. Values between 0 and 60000 are allowed.
	FrustratingFallbackThreshold *int `json:"frustratingFallbackThreshold,omitempty"` // Fallback threshold of an XHR action, defining a tolerable user experience, when the configured KPM is not available. Values between 0 and 240000 are allowed.
}

func (me *ApdexSettings) IsEmpty() bool {
	return me.Threshold == nil && me.ToleratedThreshold == nil && me.FrustratingThreshold == nil && me.ToleratedFallbackThreshold == nil && me.FrustratingFallbackThreshold == nil
}

func (me *ApdexSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"threshold": {
			Type:        schema.TypeInt,
			Description: "no documentation available",
			Optional:    true,
		},
		"tolerated_threshold": {
			Type:        schema.TypeInt,
			Description: "Maximal value of apdex, which is considered as satisfied user experience. Values between 0 and 60000 are allowed.",
			Optional:    true,
		},
		"frustrating_threshold": {
			Type:        schema.TypeInt,
			Description: "Maximal value of apdex, which is considered as tolerable user experience. Values between 0 and 240000 are allowed.",
			Optional:    true,
		},
		"tolerated_fallback_threshold": {
			Type:        schema.TypeInt,
			Description: "Fallback threshold of an XHR action, defining a satisfied user experience, when the configured KPM is not available. Values between 0 and 60000 are allowed.",
			Optional:    true,
		},
		"frustrating_fallback_threshold": {
			Type:        schema.TypeInt,
			Description: "Fallback threshold of an XHR action, defining a tolerable user experience, when the configured KPM is not available. Values between 0 and 240000 are allowed.",
			Optional:    true,
		},
	}
}

func (me *ApdexSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"threshold":                      me.Threshold,
		"tolerated_threshold":            me.ToleratedThreshold,
		"frustrating_threshold":          me.FrustratingThreshold,
		"tolerated_fallback_threshold":   me.ToleratedFallbackThreshold,
		"frustrating_fallback_threshold": me.FrustratingFallbackThreshold,
	})
}

func (me *ApdexSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"threshold":                      &me.Threshold,
		"tolerated_threshold":            &me.ToleratedThreshold,
		"frustrating_threshold":          &me.FrustratingThreshold,
		"tolerated_fallback_threshold":   &me.ToleratedFallbackThreshold,
		"frustrating_fallback_threshold": &me.FrustratingFallbackThreshold,
	})
}
