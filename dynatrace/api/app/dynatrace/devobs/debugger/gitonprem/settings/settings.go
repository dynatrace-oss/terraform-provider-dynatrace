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

package gitonprem

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	IncludeCredentials *bool     `json:"IncludeCredentials,omitempty"` // If turned on, requests to your Gitlab server will have the `credentials` option set to `include`. Otherwise, it will be set to `omit`.
	Provider           Providers `json:"Provider"`                     // Possible Values: `AzureOnPrem`, `BitbucketOnPrem`, `GithubOnPrem`, `GitlabOnPrem`
	Url                string    `json:"Url"`                          // An HTTP/HTTPS URL of your server
}

func (me *Settings) Name() string {
	return me.Url
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"include_credentials": {
			Type:        schema.TypeBool,
			Description: "If turned on, requests to your Gitlab server will have the `credentials` option set to `include`. Otherwise, it will be set to `omit`.",
			Optional:    true, // precondition
		},
		"git_provider": {
			Type:        schema.TypeString,
			Description: "Possible Values: `AzureOnPrem`, `BitbucketOnPrem`, `GithubOnPrem`, `GitlabOnPrem`",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "An HTTP/HTTPS URL of your server",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"include_credentials": me.IncludeCredentials,
		"git_provider":        me.Provider,
		"url":                 me.Url,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.IncludeCredentials == nil) && (string(me.Provider) == "GitlabOnPrem") {
		me.IncludeCredentials = opt.NewBool(false)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"include_credentials": &me.IncludeCredentials,
		"git_provider":        &me.Provider,
		"url":                 &me.Url,
	})
}
