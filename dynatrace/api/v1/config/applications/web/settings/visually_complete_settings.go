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

// VisuallyCompleteSettings Settings for VisuallyComplete
type VisuallyCompleteSettings struct {
	ExcludeURLRegex      *string `json:"excludeUrlRegex"`             // A RegularExpression used to exclude images and iframes from being detected by the VC module
	IgnoredMutationsList *string `json:"ignoredMutationsList"`        // Query selector for mutation nodes to ignore in VC and SI calculation
	MutationTimeout      *int32  `json:"mutationTimeout,omitempty"`   // Determines the time in ms VC waits after an action closes to start calculation. Defaults to 50. Valid values range from 0 to 5000.
	InactivityTimeout    *int32  `json:"inactivityTimeout,omitempty"` // The time in ms the VC module waits for no mutations happening on the page after the load action. Defaults to 1000. Valid values range from 0 to 30000.
	Threshold            *int32  `json:"threshold,omitempty"`         // Minimum visible area in pixels of elements to be counted towards VC and SI. Defaults to 50. Valid values range from 0 to 10000.
}

func (me *VisuallyCompleteSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"exclude_url_regex": {
			Type:        schema.TypeString,
			Description: "A RegularExpression used to exclude images and iframes from being detected by the VC module",
			Optional:    true,
		},
		"ignored_mutations_list": {
			Type:        schema.TypeString,
			Description: "Query selector for mutation nodes to ignore in VC and SI calculation",
			Optional:    true,
		},
		"mutation_timeout": {
			Type:        schema.TypeInt,
			Description: "Determines the time in ms VC waits after an action closes to start calculation. Defaults to 50. Valid values range from 0 to 5000.",
			Optional:    true,
		},
		"inactivity_timeout": {
			Type:        schema.TypeInt,
			Description: "The time in ms the VC module waits for no mutations happening on the page after the load action. Defaults to 1000. Valid values range from 0 to 30000.",
			Optional:    true,
		},
		"threshold": {
			Type:        schema.TypeInt,
			Description: "Minimum visible area in pixels of elements to be counted towards VC and SI. Defaults to 50. Valid values range from 0 to 10000.",
			Optional:    true,
		},
	}
}

func (me *VisuallyCompleteSettings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"exclude_url_regex":      me.ExcludeURLRegex,
		"ignored_mutations_list": me.IgnoredMutationsList,
		"mutation_timeout":       me.MutationTimeout,
		"inactivity_timeout":     me.InactivityTimeout,
		"threshold":              me.Threshold,
	})
}

func (me *VisuallyCompleteSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"exclude_url_regex":      &me.ExcludeURLRegex,
		"ignored_mutations_list": &me.IgnoredMutationsList,
		"mutation_timeout":       &me.MutationTimeout,
		"inactivity_timeout":     &me.InactivityTimeout,
		"threshold":              &me.Threshold,
	})
}
