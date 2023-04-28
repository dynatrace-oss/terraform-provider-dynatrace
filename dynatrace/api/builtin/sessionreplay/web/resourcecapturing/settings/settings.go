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

package resourcecapturing

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID                          *string  `json:"-" scope:"applicationId"`                          // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	EnableResourceCapturing                bool     `json:"enableResourceCapturing"`                          // (Field has overlap with `dynatrace_web_application`) When turned on, all CSS resources from all sessions are captured. For details, see [Resource capture](https://dt-url.net/sr-resource-capturing).
	ResourceCaptureUrlExclusionPatternList []string `json:"resourceCaptureUrlExclusionPatternList,omitempty"` // (Field has overlap with `dynatrace_web_application`) Add exclusion rules to avoid the capture of resources from certain pages.
}

func (me *Settings) Name() string {
	return *me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
		},
		"enable_resource_capturing": {
			Type:        schema.TypeBool,
			Description: "(Field has overlap with `dynatrace_web_application`) When turned on, all CSS resources from all sessions are captured. For details, see [Resource capture](https://dt-url.net/sr-resource-capturing).",
			Required:    true,
		},
		"resource_capture_url_exclusion_pattern_list": {
			Type:        schema.TypeSet,
			Description: "(Field has overlap with `dynatrace_web_application`) Add exclusion rules to avoid the capture of resources from certain pages.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":                              me.ApplicationID,
		"enable_resource_capturing":                   me.EnableResourceCapturing,
		"resource_capture_url_exclusion_pattern_list": me.ResourceCaptureUrlExclusionPatternList,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":                              &me.ApplicationID,
		"enable_resource_capturing":                   &me.EnableResourceCapturing,
		"resource_capture_url_exclusion_pattern_list": &me.ResourceCaptureUrlExclusionPatternList,
	})
}
