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
	Name     string `json:"name"`     // The name of the Jenkins connection
	Password string `json:"password"` // The password of the user or API token obtained from the Jenkins UI (Dashboard > User > Configure > API Token)
	Url      string `json:"url"`      // Base URL of your Jenkins instance (e.g. https://[YOUR_JENKINS_DOMAIN]/)
	Username string `json:"username"` // The name of your Jenkins user (e.g. jenkins)
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the Jenkins connection",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "The password of the user or API token obtained from the Jenkins UI (Dashboard > User > Configure > API Token)",
			Required:    true,
			Sensitive:   true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "Base URL of your Jenkins instance (e.g. https://[YOUR_JENKINS_DOMAIN]/)",
			Required:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The name of your Jenkins user (e.g. jenkins)",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":     me.Name,
		"password": "${state.secret_value}",
		"url":      me.Url,
		"username": me.Username,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":     &me.Name,
		"password": &me.Password,
		"url":      &me.Url,
		"username": &me.Username,
	})
}

func (me *Settings) FillDemoValues() []string {
	if me.Password != "" {
		me.Password = "#######"
	}
	return []string{"REST API didn't provide password data"}
}
