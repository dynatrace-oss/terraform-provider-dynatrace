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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TechnologyAttributes struct {
	CustomMatcher *string `json:"customMatcher,omitempty"` // Custom matching condition which should be used instead of technology matcher.
	TechnologyID  string  `json:"technologyId"`            // Technology ID
}

func (me *TechnologyAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_matcher": {
			Type:        schema.TypeString,
			Description: "Custom matching condition which should be used instead of technology matcher.",
			Optional:    true, // nullable
		},
		"technology_id": {
			Type:        schema.TypeString,
			Description: "Technology ID",
			Required:    true,
		},
	}
}

func (me *TechnologyAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"custom_matcher": me.CustomMatcher,
		"technology_id":  me.TechnologyID,
	})
}

func (me *TechnologyAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"custom_matcher": &me.CustomMatcher,
		"technology_id":  &me.TechnologyID,
	})
}
