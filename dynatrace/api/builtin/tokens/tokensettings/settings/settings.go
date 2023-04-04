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

package tokensettings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	NewTokenFormat bool `json:"newDynatraceTokenFormatEnabled"` // Check out this [blog post](http://www.dynatrace.com/blog/further-increased-security-of-your-api-tokens-by-automating-token-protection/) to find out more about the new Dynatrace API token format.
	PersonalTokens bool `json:"patEnabled"`                     // Allow users of this environment to generate personal access tokens based on user permissions. \n Note that existing personal access tokens will become unusable while this setting is disabled.
}

func (me *Settings) Name() string {
	return "token_settings"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"new_token_format": {
			Type:        schema.TypeBool,
			Description: "Check out this [blog post](http://www.dynatrace.com/blog/further-increased-security-of-your-api-tokens-by-automating-token-protection/) to find out more about the new Dynatrace API token format.",
			Required:    true,
		},
		"personal_tokens": {
			Type:        schema.TypeBool,
			Description: "Allow users of this environment to generate personal access tokens based on user permissions. \n Note that existing personal access tokens will become unusable while this setting is disabled.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"new_token_format": me.NewTokenFormat,
		"personal_tokens":  me.PersonalTokens,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"new_token_format": &me.NewTokenFormat,
		"personal_tokens":  &me.PersonalTokens,
	})
}
