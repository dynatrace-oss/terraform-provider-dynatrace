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

package connection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name  string `json:"name"`  // A unique and clearly identifiable connection name to your GitLab instance.
	Token string `json:"token"` // The GitLab token to use for authentication. Please note that this token is not refreshed and can expire.
	Url   string `json:"url"`   // The GitLab URL instance you want to connect. For example, https://gitlab.com
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "A unique and clearly identifiable connection name to your GitLab instance.",
			Required:    true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "The GitLab token to use for authentication. Please note that this token is not refreshed and can expire.",
			Required:    true,
			Sensitive:   true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The GitLab URL instance you want to connect. For example, https://gitlab.com",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"token": "${state.secret_value}",
		"url":   me.Url,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"token": &me.Token,
		"url":   &me.Url,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.Token != "" {
		me.Token = "#######"
	}
	return []string{"REST API didn't provide token data"}
}
