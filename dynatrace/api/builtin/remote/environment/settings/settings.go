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

package environment

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name         string `json:"name"`         // Name
	NetworkScope Scope  `json:"networkScope"` // Possible Values: `CLUSTER`, `EXTERNAL`, `INTERNAL`
	Token        string `json:"token"`        // Provide a valid token created on the remote environment.
	Uri          string `json:"uri"`          // Specify the full URI to the remote environment. Your local environment will have to be able to connect this URI on a network level.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"network_scope": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CLUSTER`, `EXTERNAL`, `INTERNAL`",
			Required:    true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "Provide a valid token created on the remote environment.",
			Sensitive:   true,
			Required:    true,
		},
		"uri": {
			Type:        schema.TypeString,
			Description: "Specify the full URI to the remote environment. Your local environment will have to be able to connect this URI on a network level.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":          me.Name,
		"network_scope": me.NetworkScope,
		"token":         "${state.secret_value}",
		"uri":           me.Uri,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":          &me.Name,
		"network_scope": &me.NetworkScope,
		"token":         &me.Token,
		"uri":           &me.Uri,
	})
}

const credsNotProvided = "REST API didn't provide token data"

func (me *Settings) FillDemoValues() []string {
	me.Token = "################"
	return []string{credsNotProvided}
}
